package main

import (
	"cash-server/configs"
	"cash-server/db"
	"cash-server/model"
	"cash-server/pkg/util"
	"cash-server/router"
	"fmt"
)

// @title  金流SERVER API
// @version 2020.07
// @description This is a sample API server .
// @contact.name  CQI-SERVICES
// @contact.url https://www.cqiserv.com/
// @contact.email zor@cqiserv.com
func main() {
	// <== 測試CODE == >
	fmt.Println(model.PlatformQueryExist("test"))
	// <== 測試CODE == >
	// server start
	util.Success(" < - SERVER START - > ")
	//DB config
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
	util.Success(" < - SERVER INIT - > ")
	configs.LoadGlobalConfig("")
	util.Info("[MODE] " + configs.GetGlobalConfig().RunMode)
	util.Info(" < - MyDB INIT - >")
	if err := db.Initgorm(); err != nil {
		util.Error(err.Error())
	} else {
		util.Info("DB Host :" + configs.GetGlobalConfig().MySQL.Host)
	}
}
