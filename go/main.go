package main

import (
	"app/common"
	"app/config"
	"app/router"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"html/template"
	"net/http"
)




func init()  {

	config.GetDb()
}


func CustomRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			//if r := recover(); r != nil {
			//  log.Printf("崩溃信息:%s", r)
			//}
			if err, ok := recover().(string); ok {
				log.Printf("您可以在这里完成告警任务,邮件,微信等告警")
				c.String(http.StatusInternalServerError, fmt.Sprintf("error: %s", err))
			}else{
				log.Printf("aaa %#v",err)
			}
			c.AbortWithStatus(http.StatusBadRequest)
		}()
		c.Next()
	}
}

func main()  {

	r := route.SetRoute()
	r.SetFuncMap(template.FuncMap{
		"unescaped":common.Unescaped,
	})

	r.Use(gin.Logger())
	r.Use(CustomRecovery())
	r.LoadHTMLGlob("view/*")

	//r.StaticFS("/static",http.Dir("./public/static"))
	//r.StaticFile("mix-manifest.json","./public/mix-manifest.json")
	//r.StaticFile("favicon.ico", "./public/favicon.ico")

	r.Run(":9090")
}

