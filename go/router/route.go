package route

import (
	"app/controller"
	"net/http"

	//"app/middleware"
	"github.com/gin-gonic/gin"
)


func SetRoute() *gin.Engine  {
	r := gin.Default()

	r.GET("/",controller.Index)
	r.GET("/cat/:code",controller.Cat)
	r.GET("/detail/:id",controller.Detail)

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"code" : "404",
			"msg" : "页面不存在",
		})
	})

	//全局注册中间件
	//r.Use(middleware.Http{}.Auth(true))

	/*
	videoGroup := r.Group("/home")
	{
		//路由组注册中间件
		//videoGroup.Use(middleware.Http{}.Auth())
		videoGroup.GET("/index",controller.Home{}.Index)
	}
	uploadGroup := r.Group("/upload")
	{
		uploadGroup.POST("/one",controller.Upload)
		uploadGroup.POST("/mult",controller.Uploads)
	}

	redirectGroup := r.Group("/redirect")
	{
		redirectGroup.GET("/index",controller.Redirect)
		redirectGroup.GET("/route",func(context *gin.Context){
			controller.Redirect_route(r,context)
		})
	}

	 */
	return r
}
