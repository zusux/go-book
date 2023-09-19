package middleware

import (
	"app/common"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Http struct {
}

//登录中间件
func (h Http) Auth(doCheck bool)gin.HandlerFunc  {

	return func(c *gin.Context){
		if doCheck {
			tokenString := c.GetHeader("X-Token")
			if tokenString != ""{
				uid,err := common.ParseToken(tokenString)
				if err != nil{
					c.Abort()
					c.JSON(http.StatusOK,gin.H{
						"code":"400",
						"message":err.Error(),
						"data":nil,
					})
				}else{
					c.Set("uid",uid)
					c.Next()
				}
			}else{
				c.Abort()
				c.JSON(http.StatusOK,gin.H{
					"code":"400",
					"message":"请先登录",
					"data":nil,
				})
			}
		}else{
			c.Next()
		}
	}
}
