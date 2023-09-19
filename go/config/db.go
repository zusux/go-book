package config

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"time"
)

var db *gorm.DB

func init(){
	LoadEnv()
	db = NewConnection()
}

func NewConnection() *gorm.DB  {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", Conf.Mysql.Username, Conf.Mysql.Password, Conf.Mysql.Host, Conf.Mysql.Port, Conf.Mysql.Database)
	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		panic("连接数据库失败, error=" + err.Error())
	}
	db.DB().SetMaxIdleConns(10)
	// SetMaxOpenCons 设置数据库的最大连接数量。
	db.DB().SetMaxOpenConns(100)
	// SetConnMaxLifetiment 设置连接的最大可复用时间。
	db.DB().SetConnMaxLifetime(time.Hour)
	return db
}

func GetDb() *gorm.DB {
	if err := db.DB().Ping(); err != nil{
		db.Close()
		db = NewConnection()
	}
	return db
}