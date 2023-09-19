package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Cat struct {
	Id       int    `gorm:"id"`
	ParentId int    `gorm:"parent_id"`
	Code     string `gorm:"code"`
	CatName  string `gorm:"cat_name"`
	Status   int    `gorm:"status"`
}

type Book struct {
	Id             int    `gorm:"id"`
	CatId          int    `gorm:"cat_id"`
	Name           string `gorm:"name"`
	Author         string `gorm:"author"`
	Intro          string `gorm:"intro"`
	Url            string `gorm:"url"`
	Img            string `gorm:"img"`
	CreatedAt      int64  `gorm:"created_at"`
	UpdatedAt      int64  `gorm:"updated_at"`
	SourceUrl      string `gorm:"source_url"`
	Status         int    `gorm:"status"`
	LastChapterNum int    `gorm:"last_chapter_num"`
	LastChapter    string `gorm:"last_chapter"`
	LastUpdate     string `gorm:"last_update"`
	Click          int    `gorm:"click"`
	IsHot          int    `gorm:"is_hot"`
	IsNew          int    `gorm:"is_new"`
	IsOver         int    `gorm:"is_over"`
}

type Chapter struct {
	Id        int    `gorm:"id"`
	BookId    int    `gorm:"book_id"`
	ParentId  int    `gorm:"parent_id"`
	Url       string `gorm:"url"`
	Title     string `gorm:"title"`
	CreatedAt int64  `gorm:"created_at"`
	UpdatedAt int64  `gorm:"updated_at"`
	SourceUrl string `gorm:"source_url"`
	Sort      int    `gorm:"sort"`
}

type ContentTable struct {
	Id        int    `gorm:"id"`
	ChapterId int    `gorm:"chapter_id"`
	BookId    int    `gorm:"book_id"`
	Content   string `gorm:"content"`
}

var (
	db             *gorm.DB
	ExitChan       chan bool
	GetHtmlChan    chan map[string]interface{}
	ResolvHtmlChan chan map[string]interface{}
)
var mysqlDns = "root:123456@(127.0.0.1)/book?charset=utf8&parseTime=True&loc=Local"

func init() {
	var err error
	db, err = gorm.Open("mysql", mysqlDns)
	if err != nil {
		panic(err)
	}

	ExitChan = make(chan bool)
	GetHtmlChan = make(chan map[string]interface{}, 10000)
	ResolvHtmlChan = make(chan map[string]interface{}, 10000)

	start := make(map[string]interface{})
	start["url"] = "http://www.biquge.info/"
	start["types"] = "cat"

	GetHtmlChan <- start
}

func main() {
	fmt.Println("开启协程")
	go product()
	go consume()

	ticker := time.NewTicker(time.Second * 60 * 60 * 24)
	for {
		select {
		case <-ticker.C:
			ExitChan <- true
		}
	}

	<-ExitChan
}

func product() {
	for {
		time.Sleep(time.Second * 2)
		obj := <-GetHtmlChan
		go product_request(obj)
	}
}
func product_request(obj map[string]interface{}) {
	url := obj["url"].(string)
	fmt.Println("开始请求地址:", url)
	htmlReader, err := http_request("GET", url, "")
	if err != nil {
		fmt.Println("获取url出错:", url, err)
		return
	}
	fmt.Println("结束请求地址:", url)
	obj["reader"] = htmlReader
	ResolvHtmlChan <- obj
}

func consume() {
	for {
		fmt.Println("开始解析数据")
		time.Sleep(time.Microsecond * 100)

		obj := <-ResolvHtmlChan
		types := obj["types"].(string)
		switch types {
		case "cat":
			fmt.Println("进入cat")
			go resolvCat(obj)
		case "book":
			fmt.Println("进入book")
			go resolvBook(obj)
		case "bookNext":
			fmt.Println("进入bookNext")
			go resolvBookNext(obj)
		case "bookPage":
			fmt.Println("进入bookPage")
			go resolvBookPage(obj)
		case "chapter":
			fmt.Println("进入chapter")
			go resolvChapter(obj)
		case "content":
			go resolvContent(obj)
			fmt.Println("进入content")
		default:
			fmt.Println("未知类型")
		}
	}
}

func resolvCat(obj map[string]interface{}) {
	htmlRead := obj["reader"].(io.ReadCloser)

	doc, err := goquery.NewDocumentFromReader(htmlRead)
	if err != nil {
		fmt.Println("html 生成文档失败:", err)
		return
	}

	body, err := ioutil.ReadAll(htmlRead)
	if len(string(body)) < 1000 {
		href, exists := doc.Find("center>h3>a").Attr("href")
		if exists {
			start := make(map[string]interface{})
			start["url"] = "http://www.biquge.info" + href
			start["types"] = "cat"
			GetHtmlChan <- start
			return
		}
	}

	doc.Find(".nav>ul>li").Each(func(i int, selection *goquery.Selection) {
		url, _ := selection.Find("a").Attr("href")
		catName := selection.Find("a").Text()

		cat := Cat{}
		db.Table("b_cat").Where("cat_name=?", catName).First(&cat)

		bookStartObj := map[string]interface{}{}
		if cat.Id > 0 {
			//数据库已经存在
			fmt.Println("查找分类数据ID: ", cat.Id, cat.CatName)
			bookStartObj["url"] = url
			bookStartObj["cat"] = cat
			bookStartObj["types"] = "book"
		} else {
			//数据库不存在
			catInsert := Cat{
				ParentId: 0,
				Code:     catName,
				CatName:  catName,
				Status:   0,
			}
			db.Table("b_cat").Create(&catInsert)
			fmt.Println("插入分类数据ID: ", catInsert.Id, catName)
			bookStartObj["url"] = url
			bookStartObj["cat"] = catInsert
			bookStartObj["types"] = "book"
		}

		//发送到详情生产者处理
		GetHtmlChan <- bookStartObj
	})
}

func resolvBook(obj map[string]interface{}) {
	htmlRead := obj["reader"].(io.ReadCloser)
	sourceUrl := obj["url"].(string)
	catObj, ok := obj["cat"]
	if !ok {
		return
	}
	cat := catObj.(Cat)

	if cat.CatName == "我的书架" {
		return
	}

	doc, err := goquery.NewDocumentFromReader(htmlRead)
	if err != nil {
		fmt.Println("html 生成文档失败:", err)
		return
	}
	fmt.Println("类型:", cat.CatName)
	if cat.CatName == "全部小说" || cat.CatName == "全本小说" || cat.CatName == "排行榜单" {
		lastNumStr := doc.Find("#pagelink>a.last").Text()
		fmt.Println("获取总页数字符串: #pagelink>a.last ", lastNumStr)
		lastNum, err := strconv.Atoi(lastNumStr)
		fmt.Println("总页数:", lastNum)
		if err == nil {
			switch cat.CatName {
			case "排行榜单":
				for i := 1; i <= lastNum; i++ {
					bookPageObj := map[string]interface{}{}
					tempUrl := fmt.Sprintf("http://www.biquge.info/paihangbang_allvisit/%d.html", i)
					fmt.Println("排行榜单", tempUrl)
					bookPageObj["url"] = tempUrl
					bookPageObj["types"] = "bookPage"
					GetHtmlChan <- bookPageObj
				}
			default:
				for i := 1; i <= lastNum; i++ {
					bookPageObj := map[string]interface{}{}
					tempUrl := fmt.Sprintf("%s%d", sourceUrl, i)
					fmt.Println(cat.CatName, tempUrl)
					bookPageObj["url"] = tempUrl
					bookPageObj["types"] = "bookPage"
					GetHtmlChan <- bookPageObj
				}
			}
		} else {
			fmt.Printf("页数转换错误 %v \n", err)
		}
	} else {

		lastNumStr := doc.Find("#pagelink>a.last").Text()
		fmt.Println(" else 获取总页数字符串: #pagelink>a.last ", lastNumStr)
		lastNum, err := strconv.Atoi(lastNumStr)
		fmt.Println("总页数:", lastNum)
		if err == nil {
			for i := 1; i <= lastNum; i++ {
				bookPageObj := map[string]interface{}{}
				tempUrl := strings.Replace(sourceUrl, "1.html", fmt.Sprintf("%d.html", i), 1)
				fmt.Println(cat.CatName, tempUrl)
				bookPageObj["url"] = tempUrl
				bookPageObj["types"] = "bookNext"
				bookPageObj["cat"] = cat
				GetHtmlChan <- bookPageObj
			}
		} else {
			fmt.Printf("页数转换错误 %v \n", err)
		}

	}
}

func resolvBookNext(obj map[string]interface{}) {

	htmlRead := obj["reader"].(io.ReadCloser)
	sourceUrl := obj["url"].(string)
	catObj, ok := obj["cat"]
	if !ok {
		return
	}
	cat := catObj.(Cat)

	doc, err := goquery.NewDocumentFromReader(htmlRead)
	if err != nil {
		fmt.Println("html 生成文档失败:", err)
		return
	}
	fmt.Println("类型:", cat.CatName)

	//最近更新
	doc.Find("#newscontent>div>ul>li").Each(func(i int, selection *goquery.Selection) {
		//types := selection.Find(".s1").Text()
		url, _ := selection.Find(".s2>a").Attr("href")
		name := selection.Find(".s2>a").Text()
		lastChapter := selection.Find(".s3>a").Text()
		//newChapterUrl,_ := selection.Find(".s3>a").Attr("href")
		author := selection.Find(".s4").Text()
		//date := selection.Find(".s5").Text()

		//查找
		book := Book{}
		db.Table("b_book").Where("name=?", name).Where("url = ?", url).First(&book)
		chapterStartObj := map[string]interface{}{}
		if book.Id > 0 {
			//更新

			chapterStartObj["url"] = url
			chapterStartObj["book_id"] = book.Id
			chapterStartObj["types"] = "chapter"
			GetHtmlChan <- chapterStartObj
		} else {

			updateAt := time.Now().Unix()
			fmt.Println("更新时间:", updateAt)
			bookInsert := Book{
				Name:        name,
				Author:      author,
				CatId:       cat.Id,
				Url:         url,
				SourceUrl:   sourceUrl,
				CreatedAt:   time.Now().Unix(),
				UpdatedAt:   updateAt,
				LastChapter: lastChapter,
			}
			db.Table("b_book").Create(&bookInsert)
			fmt.Println("插入book数据ID: ", bookInsert.Id)

			chapterStartObj["url"] = url
			chapterStartObj["book_id"] = bookInsert.Id
			chapterStartObj["types"] = "chapter"
			GetHtmlChan <- chapterStartObj
		}
	})

	//最新小说
	doc.Find("#newscontent>.r>ul>li").Each(func(i int, selection *goquery.Selection) {

		url, _ := selection.Find(".s2>a").Attr("href")
		name := selection.Find(".s2>a").Text()
		author := selection.Find(".s4").Text()

		//查找
		book := Book{}
		db.Table("b_book").Where("name=?", name).Where("url = ?", url).First(&book)
		chapterStartObj := map[string]interface{}{}
		if book.Id > 0 {
			//更新

			chapterStartObj["url"] = url
			chapterStartObj["book_id"] = book.Id
			chapterStartObj["types"] = "chapter"
			GetHtmlChan <- chapterStartObj
		} else {

			bookInsert := Book{
				Name:      name,
				Author:    author,
				CatId:     cat.Id,
				Url:       url,
				SourceUrl: sourceUrl,
				IsNew:     1,
				CreatedAt: time.Now().Unix(),
			}
			db.Table("b_book").Create(&bookInsert)
			fmt.Println("插入book数据ID: ", bookInsert.Id)

			chapterStartObj["url"] = url
			chapterStartObj["book_id"] = bookInsert.Id
			chapterStartObj["types"] = "chapter"
			GetHtmlChan <- chapterStartObj
		}
	})

	//热门小说
	doc.Find("#hotcontent>.l>.item").Each(func(i int, selection *goquery.Selection) {

		intro := selection.Find("dl>dd").Text()
		url, _ := selection.Find("dl>dt>a").Attr("href")
		name := selection.Find("dl>dt>a").Text()
		author := selection.Find("dl>dt>span").Text()
		img, _ := selection.Find(".image>a>img").Attr("src")

		//查找
		book := Book{}
		db.Table("b_book").Where("name=?", name).Where("url = ?", url).First(&book)
		chapterStartObj := map[string]interface{}{}
		if book.Id > 0 {
			//更新
			chapterStartObj["url"] = url
			chapterStartObj["book_id"] = book.Id
			chapterStartObj["types"] = "chapter"
			GetHtmlChan <- chapterStartObj
		} else {

			bookInsert := Book{
				Name:      name,
				Author:    author,
				CatId:     cat.Id,
				Url:       url,
				SourceUrl: sourceUrl,
				Img:       img,
				Intro:     intro,
				IsHot:     1,
				CreatedAt: time.Now().Unix(),
			}
			db.Table("b_book").Create(&bookInsert)
			fmt.Println("插入book数据ID: ", bookInsert.Id)

			chapterStartObj["url"] = url
			chapterStartObj["book_id"] = bookInsert.Id
			chapterStartObj["types"] = "chapter"
			GetHtmlChan <- chapterStartObj
		}
	})
}

func resolvBookPage(obj map[string]interface{}) {

	sourceUrl := obj["url"].(string)
	htmlRead := obj["reader"].(io.ReadCloser)

	doc, err := goquery.NewDocumentFromReader(htmlRead)
	if err != nil {
		fmt.Println("html 生成文档失败:", err)
		return
	}
	doc.Find(".novelslistss>ul>li").Each(func(i int, selection *goquery.Selection) {

		types := selection.Find(".s1").Text()
		name := selection.Find(".s2>a").Text()
		url, _ := selection.Find(".s2>a").Attr("href")
		lastChapter := selection.Find(".s3>a").Text()
		author := selection.Find(".s4").Text()
		date := selection.Find(".s5").Text()
		isOverStr := selection.Find(".s7").Text()
		isOver := 0
		if strings.Contains(isOverStr, "完成") {
			isOver = 1
		}

		//查找
		book := Book{}
		db.Table("b_book").Where("name=?", name).Where("url = ?", url).First(&book)
		chapterStartObj := map[string]interface{}{}
		if book.Id > 0 {
			//更新

			chapterStartObj["url"] = url
			chapterStartObj["book_id"] = book.Id
			chapterStartObj["types"] = "chapter"
			GetHtmlChan <- chapterStartObj
		} else {
			cat := Cat{}
			//查找分类
			catName := strings.Trim(types, "[]")
			db.Table("b_cat").Where("cat_name=?", catName).First(&cat)

			catId := 0
			if cat.Id > 0 {
				catId = cat.Id
			}

			ymd := fmt.Sprintf("%d%s", 20, date)
			loc, _ := time.LoadLocation("Asia/Shanghai")                               //设置时区
			tt, _ := time.ParseInLocation("2006-01-02 15:04:05", ymd+" 00:00:00", loc) //2006-01-02 15:04:05是转换的格式如php的"Y-m-d H:i:s"
			updateAt := tt.Unix()

			bookInsert := Book{
				Name:        name,
				Author:      author,
				CatId:       catId,
				Url:         url,
				Status:      1,
				SourceUrl:   sourceUrl,
				CreatedAt:   time.Now().Unix(),
				UpdatedAt:   updateAt,
				LastChapter: lastChapter,
				IsOver:      isOver,
			}
			db.Table("b_book").Create(&bookInsert)
			fmt.Println("插入book数据ID: ", bookInsert.Id)

			chapterStartObj["url"] = url
			chapterStartObj["book_id"] = bookInsert.Id
			chapterStartObj["types"] = "chapter"
			GetHtmlChan <- chapterStartObj
		}
	})
}

func resolvChapter(obj map[string]interface{}) {
	sourceUrl := obj["url"].(string)
	book_id := obj["book_id"].(int)
	htmlRead := obj["reader"].(io.ReadCloser)

	doc, err := goquery.NewDocumentFromReader(htmlRead)
	if err != nil {
		fmt.Println("html 生成文档失败:", err)
		return
	}
	doc.Find("#list>dl>dd").Each(func(i int, selection *goquery.Selection) {

		title := selection.Find("a").Text()
		url, _ := selection.Find("a").Attr("href")

		url = sourceUrl + url

		//查找
		chapter := Chapter{}
		db.Table("b_chapter").Where("book_id=?", book_id).Where("title = ?", title).First(&chapter)
		contentStartObj := map[string]interface{}{}
		if chapter.Id > 0 {
			//更新

		} else {
			chapterInsert := Chapter{
				Url:       url,
				SourceUrl: sourceUrl,
				BookId:    book_id,
				Title:     title,
				CreatedAt: time.Now().Unix(),
				Sort:      i + 1,
			}
			db.Table("b_chapter").Create(&chapterInsert)
			fmt.Println("插入chapter数据ID: ", chapterInsert.Id)

			contentStartObj["url"] = url
			contentStartObj["book_id"] = book_id
			contentStartObj["chapter_id"] = chapterInsert.Id
			contentStartObj["types"] = "content"
			GetHtmlChan <- contentStartObj
		}
	})
}

func resolvContent(obj map[string]interface{}) {
	book_id := obj["book_id"].(int)
	chapter_id := obj["chapter_id"].(int)
	htmlRead := obj["reader"].(io.ReadCloser)
	doc, err := goquery.NewDocumentFromReader(htmlRead)
	if err != nil {
		fmt.Println("html 生成文档失败:", err)
		return
	}
	content := doc.Find("#content").Text()

	reg3, _ := regexp.Compile("笔.趣.阁")
	content = reg3.ReplaceAllString(content, "")
	content = strings.Replace(content, "ｗWｗ。ｂｉｑUｇE。ｉｎｆｏ", "", -1)

	//查找
	contentTable := ContentTable{}
	db.Table("b_content").Where("book_id=?", book_id).Where("chapter_id = ?", chapter_id).First(&contentTable)
	if contentTable.Id > 0 {
		//更新
	} else {
		contentInsert := ContentTable{
			BookId:    book_id,
			ChapterId: chapter_id,
			Content:   content,
		}
		db.Table("b_content").Create(&contentInsert)
		fmt.Println("插入content数据ID: ", contentInsert.Id)
	}
}

func http_request(method string, url string, msg string) (htmlReader io.ReadCloser, err error) {
	//fmt.Println("请求地址:",url)
	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				conn, err := net.DialTimeout(netw, addr, time.Second*30)
				if err != nil {
					return nil, err
				}
				conn.SetDeadline(time.Now().Add(time.Second * 30))
				return conn, nil
			},
			ResponseHeaderTimeout: time.Second * 30,
		},
	}
	body := bytes.NewBuffer([]byte(msg))
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		fmt.Println("构建错误")
		return
	}
	req.Header.Set("Content-Type", "text/html; charset=UTF-8")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.125 Safari/537.36 Edg/84.0.522.59")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("请求错误:", err.Error())
		return
	}
	if resp.StatusCode != 200 {
		return htmlReader, errors.New(fmt.Sprintf("状态码错误: %d", resp.StatusCode))
	}
	htmlReader = resp.Body
	return
}
