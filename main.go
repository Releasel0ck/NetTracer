package main

import (
	_ "database/sql"

	"io/ioutil"
	"net/http"
	"strings"

	_ "github.com/mattn/go-sqlite3"

	"NetTracer/controllers"
	"NetTracer/models"
	"NetTracer/util"

	"github.com/gin-gonic/gin"
)

var err error

var DBname string

func main() {

	//如果存在数据库，则选择一个连接并打开。
	dbs := strings.Split(util.GetArbitraryDB(), "$")
	if len(dbs) > 0 {
		DBname = dbs[0]
		if DBname != "" {
			models.ConnectDB(DBname)
		}
	}

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Static("/statics", "./statics")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "NetTracer",
		})
	})
	//获取IP和Port按连接数量的排名数据
	r.POST("/getPortAndIPRank", controllers.GetPortAndIPRank)
	//获取已有的所有数据库
	r.POST("/getDatabase", controllers.GetDatabases)
	//获取当前连接的数据库
	r.POST("/getCurrentDatabase", func(c *gin.Context) {
		c.String(http.StatusOK, DBname)
	})
	//选择一个新的数据库进行连接
	r.POST("/selectDatabase", func(c *gin.Context) {
		data, _ := ioutil.ReadAll(c.Request.Body)
		DBname = string(data)
		models.ConnectDB(DBname)
		c.String(http.StatusOK, DBname)
	})
	//创建一个新的数据库
	r.POST("/createDatabase", controllers.CreateDatabase)
	//测试使用...
	r.GET("/status", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})
	//查询数据
	r.POST("/query", controllers.QueryData)
	//上次并解析网络连接文件
	r.POST("/upload", controllers.ParseText)
	//监听7474端口
	r.Run(":7474")
}
