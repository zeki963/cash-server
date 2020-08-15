package router

import (
	"cash-server/configs"
	"path/filepath"

	// cash-server/docs swag DOC
	_ "cash-server/docs"

	"cash-server/controller"
	"cash-server/pkg/util"

	"fmt"
	"io"
	"os"

	"github.com/gin-contrib/multitemplate"
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
		r.Use(controller.LoggerToDB())
	}
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.StaticFile("/favicon.ico", "./favicon.ico")
	r.Static("/static", "./templates/static")
	r.HTMLRender = loadTemplates("./templates")
	r.Any("/", controller.Homepage)
	r.GET("/demo", controller.Demopage)

	rmycard := r.Group("mycardsandbox")
	{
		rmycard.POST("/order", controller.MycardSandOderAdd) //新增 mycard 建單 Add
		rmycard.GET("/order/:key")                           //查詢 mycard 查詢單筆交易
		rmycard.GET("/orders")                               //查詢 mycard 查詢交易清單
	}
	rmycardCall := r.Group("mycardcall", MycardCallAuth()) //rmycardC  Mycard Call back!! 榜定白名單
	{
		rmycardCall.POST("/odercallback", controller.CallbackMycard)             //給Mycard廠商用ReturnURL 3.2
		rmycardCall.POST("/transactioncallback", controller.Transactioncallback) //給Mycard廠商用ReturnURL 3.6
		rmycardCall.POST("/transactioncheck", controller.TransactionCheck)       //給Mycard廠商用ReturnURL 3.7
	}
	rcasino := r.Group("casino")
	{
		rcasino.GET("/Shop") //查詢 casino 查詢商品清單
	}
	radmin := r.Group("admin")
	{
		radmin.POST("/platform", controller.PlatformRegisterServerAdd)
		radmin.GET("/platform/:Token", controller.PlatformGet)
		radmin.GET("/platforms")
		radmin.PUT("/platform/status/:Token", controller.PlatformStatusEnable)
		radmin.DELETE("/platform/status/:Token", controller.PlatformStatusDisable)
	}

	if configs.GetGlobalConfig().Swagger == true { //Swagger
		r.Any("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
	r.NoRoute(controller.NoResponse)
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

//MycardCallAuth Mycard call 白名單
func MycardCallAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		whiteList := []string{
			"127.0.0.1",
			"218.32.37.148",
			"220.130.127.125",
		}
		flag := false
		for _, host := range whiteList {
			if c.ClientIP() == host {
				flag = true
				break
			}
		}
		if !flag {
			c.JSON(511, "your ip is not trusted")
			c.Abort()
		}
	}
}
