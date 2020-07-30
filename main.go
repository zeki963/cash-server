package main

import (
	"cash-server/configs"
	"cash-server/db"
	"cash-server/model"
	"cash-server/pkg/flag"
	"cash-server/pkg/util"
	"cash-server/routers"
	"cash-server/routers/api/mycard"
)

func init() {
	util.Success(" < - SERVER INIT - > ")
	flag.StartFlag()
	// <== 測試CODE == >
	model.PlatformQueryInfodata("0OEaFvXvGXVjlanvAXZugA")
	mycard.Savedb()
}

// @title  金流SERVER API
// @version 2020.07
// @description This is a sample API server .
// @contact.name  CQI-SERVICES
// @contact.url https://www.cqiserv.com/
// @contact.email zor@cqiserv.com
func main() {
	// server start
	util.Success(" < - SERVER START - > ")
	//DB config
	defer db.SqlDB.Close()
	//ROUTER
	router := routers.InitRouter()
	//PORT
	if configs.GetGlobalConfig().HTTPS == true {
		router.RunTLS(":8443", "templates/server.crt", "templates/server.key")
	} else {
		if err := router.Run(":8080"); err != nil {
			util.Error(err.Error())
		}
	}
}
