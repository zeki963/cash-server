package routers

import (
	"cash-server/configs"
	"path/filepath"

	// cash-server/docs swag DOC
	_ "cash-server/docs"

	"cash-server/pkg/util"
	"cash-server/routers/api/admin"
	"cash-server/routers/api/mycard"
	"cash-server/routers/api/pay"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

var (
	key string
)

//InitRouter 初始化路由表
func InitRouter() *gin.Engine {
	//Mode switch
	switch configs.GetGlobalConfig().RunMode {
	case "release":
		gin.SetMode(gin.ReleaseMode)
		ginlogmode()
		util.Info(" < - ROUTER START - > ")
	case "test":
		gin.SetMode(gin.TestMode)
	case "debug":
		gin.SetMode(gin.DebugMode)
	}
	r := gin.Default()
	if configs.GetGlobalConfig().Logconf.LoggerToFile {
		r.Use(util.LoggerToFile())
	}
	if configs.GetGlobalConfig().Logconf.LoggerToDB {
		r.Use(util.LoggerToDB())
	}
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.StaticFile("/favicon.ico", "./favicon.ico")
	r.Static("/static", "./templates/static")
	// 根目錄
	r.Any("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{
			"title": "Home page",
		})
		util.Logger().WithFields(logrus.Fields{
			"name": "Info",
		}).Info("有人連進根目錄了", "Info")
	})

	//前端DEMO用
	r.HTMLRender = loadTemplates("./templates")
	r.GET("/demo", func(c *gin.Context) {
		c.HTML(200, "demo.html", gin.H{
			"title": "DEMO page",
		})
	})

	//Group pay1
	pay1 := r.Group("pay1")
	{
		pay1.POST("/A", pay.UrlencodedPost)
		pay1.POST("/B", pay.JSONtestPost)
	}

	//Group mycard
	rmycard := r.Group("mycard")
	{
		//使用 mycard 建單 操作s
		rmycard.POST("/CreateMycardOder", mycard.AuthMycard)
		//查詢 mycard 查詢單筆交易
		rmycard.POST("/QuiryMycardOderOne")
		//查詢 mycard 查詢交易清單
		//TODO 給web查詢 列出清單
		rmycard.POST("/QuiryMycardOderList")

		//給Mycard廠商用ReturnURL
		rmycard.POST("/MycardOderCallback", mycard.CallbackMycard)
	}

	//Group mycard
	rmycard := r.Group("casino")
	{
		//查詢 casino 查詢商品清單
		rmycard.POST("/QuiryMycardOderOne")
	}

	//Group testrouter
	testrouter := r.Group("test")
	{
		testrouter.POST("/A", pay.TestRegisterServer)
		testrouter.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"msg": "pong",
			})
		})
		//input output
		r.GET("/putkey/:key", func(context *gin.Context) {
			key = context.Param("key")
			fmt.Printf("Hello, %s", key)
		})
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

	//Group radmin
	radmin := r.Group("admin")
	{
		//register  提供註冊
		radmin.POST("/register", admin.PlatformRegisterServer)
	}
	//system 相關
	if configs.GetGlobalConfig().Swagger == true {
		//swag interface
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
	r.NoRoute(NoResponse)
	return r
}

//ginlogmode  server logs
func ginlogmode() {
	// Disable log's color
	gin.DisableConsoleColor()
	// Force log's color
	//gin.ForceConsoleColor()
	f, err := os.Create("gin.log")
	if err != nil {
		fmt.Println("Open Log File Failed", err)
	}
	gin.DefaultWriter = io.MultiWriter(f)
}

//NoResponse 不存在，返回404
func NoResponse(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{
		"status": 404,
		"msg":    "page not exists!,你想幹嘛ヽ(`Д´)ノ  ",
	})
}

//loadTemplates 前端DEMO用
func loadTemplates(templatesDir string) multitemplate.Renderer {
	r := multitemplate.NewRenderer()

	layouts, err := filepath.Glob(templatesDir + "/layouts/*.html")
	if err != nil {
		panic(err.Error())
	}
	includes, err := filepath.Glob(templatesDir + "/includes/*.html")
	if err != nil {
		panic(err.Error())
	}
	// Generate our templates map from our layouts/ and includes/ directories
	for _, include := range includes {
		layoutCopy := make([]string, len(layouts))
		copy(layoutCopy, layouts)
		files := append(layoutCopy, include)
		r.AddFromFiles(filepath.Base(include), files...)
	}
	return r
}
