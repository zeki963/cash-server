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

	//Group mycard
	rmycard := r.Group("mycard")
	{
		//使用 mycard 建單 Add
		rmycard.POST("/oder", controller.MycardOderAdd)
		//查詢 mycard 查詢單筆交易
		rmycard.GET("/oder/:key")
		//查詢 mycard 查詢交易清單
		rmycard.GET("/oders")
		//mycard link Call back!!
		//給Mycard廠商用ReturnURL 3.2
		rmycard.POST("/odercallback", controller.CallbackMycard)
		//給Mycard廠商用ReturnURL 3.6
		rmycard.POST("/transactioncallback", controller.Transactioncallback)
		//給Mycard廠商用ReturnURL 3.7
		rmycard.POST("/transactioncheck", controller.TransactionCheck)
	}

	//Group casino
	casino := r.Group("casino")
	{
		//查詢 casino 查詢商品清單
		casino.GET("/Shop")
	}

	//Group testrouter
	testrouter := r.Group("test")
	{
		testrouter.POST("/A", controller.UrlencodedPost)
		testrouter.POST("/B", controller.JSONtestPost)
		testrouter.GET("/putkey/:key", func(context *gin.Context) {
			key = context.Param("key")
			fmt.Printf("Hello, %s", key)
		})
	}

	//Group radmin
	radmin := r.Group("admin")
	{
		//register  提供註冊
		radmin.POST("/platform", controller.PlatformRegisterServerAdd)
	}

	//Swagger 相關
	if configs.GetGlobalConfig().Swagger == true {
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
