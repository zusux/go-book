package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/v2/client/grpc"
	"github.com/micro/go-micro/v2/web"
	"log"
	"net/http"
	"strconv"
	booklist "webapi/proto/booklist"
	chapter "webapi/proto/chapter"
	content "webapi/proto/content"
)

type BookRequest struct {
	Id int32 `json:"id"`
	CatId int32 `json:"cat_id"`
	Name string `json:"name"`
	Author string `json:"author"`
	IsHot int32 `json:"is_hot"`
	IsNew int32 `json:"is_new"`
	IsOver int32 `json:"is_over"`
	Page int32 `json:"page"`
	Limit int32 `json:"limit"`
	Order string `json:"order"`
}

type ChapterRequest struct {
	BookId int32 `json:"book_id"`
	Id int32 `json:"id"`
}

type ContentRequest struct {
	BookId int32 `json:"book_id"`
	ChapterId int32 `json:"chapter_id"`
}


func main()  {
	client := grpc.NewClient()
	r := gin.New()

	r.LoadHTMLGlob("view/*")
	r.StaticFS("/static",http.Dir("./public/static"))
	r.StaticFile("favicon.ico", "./public/favicon.ico")
	
	r.GET("/", func(c *gin.Context) {
		catIdStr := c.Query("cat")
		catId,err := strconv.Atoi(catIdStr)
		if err!= nil{
			catId = 0
		}
		c.HTML(http.StatusOK, "index.html", gin.H{
			"cat_id":catId,
		})
	})
	
	r.POST("/content", func(c *gin.Context) {

		var ContentReqInfo ContentRequest
		err := c.BindJSON(&ContentReqInfo)
		if err != nil {
			c.JSON(200, gin.H{
				"code": 400,
				"message": err.Error(),
				"data": nil,
			})
			return
		}

		/*
			调用grpc服务
		*/
		contentService := content.NewContentService("zusux.book.service.content",client)
		response, err := contentService.Call(context.Background(), &content.Request{BookId: ContentReqInfo.BookId,ChapterId: ContentReqInfo.ChapterId})
		if err != nil{
			log.Println(err)
			c.JSON(400,gin.H{
				"code": 400,
				"message":err.Error(),
				"data":"",
			})
		}else{
			c.JSON(200,gin.H{
				"code":200,
				"message":"请求成功",
				"data":response.Content,
			})
		}
	})

	r.POST("/chapter", func(c *gin.Context) {

		var ChapterReqInfo ChapterRequest
		err := c.BindJSON(&ChapterReqInfo)
		if err != nil {
			c.JSON(200, gin.H{
				"code": 400,
				"message": err.Error(),
				"data": nil,
			})
			return
		}


		/*
			调用grpc服务
		*/
		chapterService := chapter.NewChapterService("zusux.book.service.chapter",client)
		chapterResponse, err := chapterService.Call(context.Background(), &chapter.ChapterRequest{BookId: ChapterReqInfo.BookId,Id:ChapterReqInfo.Id})
		if err != nil{
			log.Println(err)
			c.JSON(400,gin.H{
				"code":400,
				"message":err.Error(),
				"data":"",
			})
			return
		}else{
			c.JSON(200,gin.H{
				"code":200,
				"message":"请求成功",
				"data":chapterResponse.ChapterResponses,
			})
		}
	})


	r.GET("/book/:id", func(c *gin.Context) {
		idStr := c.Param("id")
		c.HTML(http.StatusOK, "chapter.html", gin.H{
			"id":idStr,
		})

	})

	r.GET("/content/:bookid/:chapterid", func(c *gin.Context) {
		bookidStr := c.Param("bookid")
		chapteridStr := c.Param("chapterid")
		log.Println(bookidStr,chapteridStr,"bc")
		c.HTML(http.StatusOK, "content.html", gin.H{
			"book_id":bookidStr,
			"chapter_id":chapteridStr,
		})

	})

	r.POST("/book", func(c *gin.Context) {

		Brequest := &booklist.BookRequest{}
		var BookreqInfo BookRequest
		err := c.BindJSON(&BookreqInfo)
		if err != nil {
			c.JSON(200, gin.H{
				"code": 400,
				"message": err.Error(),
				"data": nil,
			})
			return
		}

		if BookreqInfo.Page >0 {
			Brequest.Page = BookreqInfo.Page
		}
		if BookreqInfo.Limit >0 {
			Brequest.Limit = BookreqInfo.Limit
		}
		if BookreqInfo.Id >0 {
			Brequest.Id = BookreqInfo.Id
		}
		if BookreqInfo.CatId >0 {
			Brequest.CatId = BookreqInfo.CatId
		}
		if BookreqInfo.Name != ""{
			Brequest.Name = BookreqInfo.Name
		}
		if BookreqInfo.Author != ""{
			Brequest.Author = BookreqInfo.Author
		}
		if BookreqInfo.IsHot >0 {
			Brequest.IsHot = BookreqInfo.IsHot
		}
		if BookreqInfo.IsNew >0 {
			Brequest.IsNew = BookreqInfo.IsNew
		}
		if BookreqInfo.IsOver >0 {
			Brequest.IsOver = BookreqInfo.IsOver
		}
		if BookreqInfo.Order != ""{
			Brequest.Order = BookreqInfo.Order
		}

		/*
			调用grpc服务
		*/
		booklistService := booklist.NewBooklistService("zusux.book.service.booklist",client)
		booklistResponse, err := booklistService.Call(context.Background(), Brequest)
		if err != nil{
			log.Println(err)
			c.JSON(400,gin.H{
				"code":400,
				"message":err.Error(),
				"data":"",
			})
			return
		}else{
			c.JSON(200,gin.H{
				"code":200,
				"message":"请求成功",
				"data":booklistResponse.BookResponses,
				"limit":BookreqInfo.Limit,
				"page":BookreqInfo.Page,
			})
		}
	})

	webService := web.NewService(
		web.Handler(r),
	)
	webService.Init()

	if  err := webService.Run();err != nil{
		log.Fatal(err)
	}
}