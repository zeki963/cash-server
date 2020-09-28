package router

import (
	"cash-server/configs"

	"cash-server/controller"
	// cash-server/docs swag DOC
	_ "cash-server/docs"
	"cash-server/pkg/util"

	"github.com/gin-contrib/cors"

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
	//CORS config
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowCredentials = true
	config.AddAllowHeaders("authorization")
	config.AllowHeaders = []string{"access-control-allow-credentials,access-control-allow-methods,access-control-allow-origin,content-type"}
	r.Use(cors.New(config))

	if configs.GetGlobalConfig().Logconf.LoggerToFile {
		r.Use(util.LoggerToFile())
	}
	if configs.GetGlobalConfig().Logconf.LoggerToDB {
		r.Use(controller.LoggerToDB())
	}
	//r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.StaticFile("/favicon.ico", "./favicon.ico")

	//DEMO 用
	r.Static("/static", "./templates/static")
	r.Any("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "200"})
	})
	rMycardSand := r.Group("mycard")
	{
		rMycardSand.POST("/order", controller.MycardSandOrderAdd)    //新增 mycard 建單 Add
		rMycardSand.Any("/ordercallback", controller.CallbackMycard) //給Mycard廠商用ReturnURL 3.2
	}
	rMycardCallSand := r.Group("mycardcall", MycardCallAuth()) //rmycardC  Mycard Call back!! 榜定白名單
	{
		rMycardCallSand.Any("/transactioncallback", controller.TransactionCallback) //給Mycard廠商用ReturnURL 3.6
		rMycardCallSand.Any("/transactioncheck", controller.TransactionCheck)       //給Mycard廠商用ReturnURL 3.7
	}
	rCasino := r.Group("casino")
	{
		rCasino.Any("/shop", controller.CasinoGetItem)      //查詢 casino 查詢商品清單
		rCasino.GET("/user/:Acc", controller.CasinoGetUser) //查詢 casino 查詢商品清單
	}
	rAdmin := r.Group("admin")
	{
		rAdmin.POST("/platform", controller.PlatformAdd)                         //新增 platform 帳號
		rAdmin.GET("/platform/:Acc", controller.PlatformGet)                     //查詢 platform 帳號
		rAdmin.GET("/platforms", controller.PlatformGetAll)                      //查詢 全部 platform 帳號
		rAdmin.PUT("/platform/status/:Acc", controller.PlatformStatusEnable)     //修改 platform 帳號狀態
		rAdmin.DELETE("/platform/status/:Acc", controller.PlatformStatusDisable) //修改 platform 帳號狀態
		rAdmin.GET("/order/:Subid", controller.OrderGet)                         //查詢單筆交易
		rAdmin.GET("/orders", controller.OrderGetAll)                            //查詢交易清單
		rAdmin.GET("/inventories")                                               //查詢交易報告
	}

	if configs.GetGlobalConfig().Swagger == true { //Swagger
		r.Any("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
	r.NoRoute(controller.NoResponse)
	return r
}
