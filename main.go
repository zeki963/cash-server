package main

import (
	"cash-server/configs"
	"cash-server/db"
	"cash-server/pkg/util"
	"cash-server/router"
)

// @title  金流SERVER API
// @version 2020.08
// @description This is a sample API server .
// @contact.name  CQI-SERVICES
// @contact.url https://www.cqiserv.com/
// @contact.email zor@cqiserv.com
func main() {
	// <== 測試CODE == >
	var a db.PaymentPlatform
	a.PlatformToken = "mSwgzQ6SV5hasRvQ0uJwVg"

	//fmt.Println(model.PlatformQueryExist(a))
	// <== 測試CODE == >
	// server start
	util.Info(" < - SERVER START - > ")
	//DB 連線關閉
	defer db.SQLDBX.Close()
	//ROUTER
	r := router.InitRouter()
	//PORT
	if configs.GetGlobalConfig().HTTPS == true {
		r.RunTLS(":8443", "templates/server.crt", "templates/server.key")
	} else {
		if err := r.Run(":8080"); err != nil {
			util.Error(err.Error())
		}
	}
}

func init() {
	util.Info(" < - SERVER INIT - > ")
	configs.LoadGlobalConfig("")
	util.Success("[MODE] > " + configs.GetGlobalConfig().RunMode)
	util.Info(" < - MyDB INIT - >")
	if err := db.Initgorm(); err != nil {
		util.Error("! DB Connect ERROR !")
		util.Error(err.Error())
	} else {
		util.Success("[DB Host] > " + configs.GetGlobalConfig().MySQL.Host)
	}
}
