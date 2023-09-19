package common


import (
	"app/config"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"html/template"
	"time"
)


func Unescaped (str interface{}) template.HTML {
	ss,ok := str.(string)
	if ok {
		return template.HTML(ss)
	}
	return template.HTML("")
}

//登录密码
func Password(password string) string  {
	h := md5.New()
	h.Write([]byte(password))
	return hex.EncodeToString(h.Sum(nil))
}

//解密token
func ParseToken(tokenString string) ( string, error) {
	//tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1ODg3ODEzMzUsImlhdCI6MTU4ODc3NzczNSwiaXNzIjoienVzdXguY29tIiwidWlkIjoxfQ.RIfqBwl9WCWHpEG3AJQtBqvOHLyQ-urpZAPcNG-Reuo"
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return config.Conf.Login.Key, nil
	})
	if token.Valid {
		claims := token.Claims
		uid := claims.(jwt.MapClaims)["uid"]
		return fmt.Sprintf("%v",uid),nil
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			return "",errors.New("不是token")
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			// Token is either expired or not active yet
			return "",errors.New("token过期")
		} else {
			return "",err
		}
	} else {
		return "",err
	}
}


func Token(uid string) (string , error)  {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(1)).Unix() //过期时间
	claims["iat"] = time.Now().Unix() //签发时间
	claims["iss"] = "zusux.com" //签发者
	claims["uid"] = uid //唯一身份标识
	token.Claims = claims
	tokenString, err := token.SignedString(config.Conf.Login.Key)
	return tokenString ,err
}

// 在slice中判断是否存在某个元素
func InSlice(need string, array []string) bool{
	for _,v := range array{
		if need == v {
			return true
		}
	}
	return false
}

func DelSlice(slice []string, value string) []string {
	for i :=0;i< len(slice);i++{
		if slice[i] == value{
			if len(slice) == i+1{
				slice = slice[:i]
			}else{
				slice = append(slice[:i],slice[i+1])
			}
		}
	}
	return slice
}