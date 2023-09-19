package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"regexp"
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

	//ExitChan = make(chan bool)
	//GetHtmlChan = make(chan map[string]interface{},10000)
	//ResolvHtmlChan = make(chan map[string]interface{},10000)

	bookRecord := make([]Book, 0)
	CatRecord := make([]Cat, 0)
	db.Table("b_cat").Find(&CatRecord)

	for _, v := range CatRecord {
		page := 0
		limit := 10
		for {
			db.Table("b_book").
				Select("id,url").
				Where("cat_id =?", v.Id).
				Limit(limit).
				Offset(page * limit).
				Order("updated_at").
				Find(&bookRecord)
			for _, bv := range bookRecord {
				start := make(map[string]interface{})
				start["url"] = bv.Url
				start["types"] = "chapter"
				start["book_id"] = bv.Id

				htmlReader, err := http_request("GET", bv.Url, "")
				if err != nil {
					fmt.Println("获取url出错:", bv.Url, err)
					continue
				}
				start["reader"] = htmlReader
				resolvChapter(start)
				db.Table("b_book").Where("id =?", bv.Id).Update("updated_at", time.Now().Unix())
			}
			page++
			if len(bookRecord) <= 0 {
				break
			}
		}

	}
}

func main() {
	fmt.Println("开启携程")

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

	img, _ := doc.Find("#fmimg>img").Attr("src")
	intro := doc.Find("#intro>p:nth-child(1)").Text()
	lastUpdateTime := doc.Find("#info>p:nth-child(4)").Text()
	lastUpdateChapter := doc.Find("#info>p:nth-child(5)>a").Text()
	lastUpdateTime = strings.Replace(lastUpdateTime, "最后更新&nbsp;&nbsp;:", "", 1)
	lastUpdateTime = strings.Replace(lastUpdateTime, "最后更新  :", "", 1)

	updateBook := Book{Img: img, Intro: intro, LastChapter: lastUpdateChapter, UpdatedAt: time.Now().Unix(), LastUpdate: lastUpdateTime}
	db.Table("b_book").Where("id=?", book_id).Updates(updateBook)

	fmt.Println("book_id", book_id, updateBook)

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
			fmt.Println("插入chapter数据ID: ", chapterInsert.Id, time.Now().Format("2006-01-02 15:04:05"))

			contentStartObj["url"] = url
			contentStartObj["book_id"] = book_id
			contentStartObj["chapter_id"] = chapterInsert.Id
			contentStartObj["types"] = "content"

			htmlReader, err := http_request("GET", url, "")
			if err != nil {
				fmt.Println("获取url出错:", url, err)
			} else {
				contentStartObj["reader"] = htmlReader
				resolvContent(contentStartObj)
			}

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
		fmt.Println("插入content数据ID: ", contentInsert.Id, time.Now().Format("2006-01-02 15:04:05"))
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
