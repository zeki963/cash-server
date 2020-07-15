package router

import (
	// cash-server/docs swag DOC
	config "cash-server/configs"
	_ "cash-server/docs"
	"cash-server/pkg/util"
	"cash-server/routers/api/admin"
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
	//Mode switch
	switch config.GetGlobalConfig().RunMode {
	case "release":
		gin.SetMode(gin.ReleaseMode)
		ginlog()
		util.Info(" < - ROUTER START - > ")
	case "test":
		gin.SetMode(gin.TestMode)
	case "debug":
		gin.SetMode(gin.DebugMode)
	}
	r := gin.Default()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.StaticFile("/favicon.ico", "./favicon.ico")
	r.Static("/statics", "./statics")
	// 根目錄
	r.Any("/", func(context *gin.Context) {
		context.String(http.StatusOK, "hello i'm gin server")
	})

	//Group pay1
	pay1 := r.Group("pay1")
	{
		//POST test
		pay1.POST("/form_post", pay.UrlencodedPost)
		//POST test
		pay1.POST("/json_post", pay.JSONtestPost)
		pay1.POST("/forget", pay.PayIndex)
	}

	//Group testrouter
	testrouter := r.Group("testrouter")
	{
		//ping
		testrouter.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})
		//input
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
	}

	//Group admin
	v1 := r.Group("admin")
	{
		//register  提供註冊
		v1.POST("/register", admin.RegisterServer)
		//list   查詢
		v1.POST("/list", admin.ListServer)
	}
	//system 相關
	if config.GetGlobalConfig().Swagger == true {
		//swag interface
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
	r.NoRoute(NoResponse)
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

//NoResponse 请求的url不存在，返回404
func NoResponse(c *gin.Context) {
	//返回404状态码
	c.JSON(http.StatusNotFound, gin.H{
		"status": 404,
		"error":  "404, page not exists!",
	})
}
