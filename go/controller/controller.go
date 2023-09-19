package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Index(c *gin.Context)  {
	c.HTML(http.StatusOK, "home.html", gin.H{
		"title": "home/index",
	})
}

func Cat(c *gin.Context) {
	c.HTML(http.StatusOK, "cat.html", gin.H{
		"title": "cat/index",
	})
}

func Chapter(c *gin.Context)  {
	c.HTML(http.StatusOK, "chapter.html", gin.H{
		"title": "chapter/index",
	})
}

func Detail(c *gin.Context)  {
	c.HTML(http.StatusOK, "detail.html", gin.H{
		"title": "detail/index",
	})
}