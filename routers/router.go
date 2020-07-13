package router

import (
	// cash-server/docs swag DOC
	_ "cash-server/docs"
	"cash-server/pkg/util"
	pay "cash-server/routers/api/pay"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

var (
	key string
)

//InitRouter 初始化路由表
func InitRouter() *gin.Engine {
	//log init
	ginlog()
	util.Info(" < - ROUTER START - > ")
	//關閉DEBUG
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.StaticFile("/favicon.ico", "./favicon.ico")
	r.Static("/statics", "./statics")

	//POST test
	r.POST("/form_post", pay.UrlencodedPost)
	//POST test
	r.POST("/json_post", pay.JSONtestPost)
	r.POST("/forget", pay.PayIndex)

	//輸入
	r.GET("/putkey/:key", func(context *gin.Context) {
		key = context.Param("key")
		fmt.Printf("Hello, %s", key)
	})
	//output
	r.GET("/getkey", func(context *gin.Context) {
		if key != "" {
			context.String(http.StatusOK, key)
			fmt.Printf("key is %s", key)
		} else {
			context.String(http.StatusOK, "No key.")
			fmt.Println("No key.")
		}
	})

	// 添加 Get 请求路由
	r.GET("/", func(context *gin.Context) {
		context.String(http.StatusOK, "hello gin")
	})
	//ping
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	//swag interface
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}

//ginlog  server logs
func ginlog() {
	// Disable log's color
	gin.DisableConsoleColor()
	// Force log's color
	// gin.ForceConsoleColor()
	f, err := os.Create("gin.log")
	if err != nil {
		fmt.Println("Open Log File Failed", err)
	}

	gin.DefaultWriter = io.MultiWriter(f)
}
