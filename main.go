package main

import (
	config "cash-server/configs"
	db "cash-server/database"
	encryption "cash-server/pkg/encryption"
	flag "cash-server/pkg/flag"
	util "cash-server/pkg/util"
	router "cash-server/routers"
)

func init() {
	encryption.Md5("123")
	//println(Dblink)
	flag.StartFlag()
}

// @title  金流SERVER API
// @version 2020.1.08.1.1.3.1.a
// @description This is a sample API server .

// @contact.name  CQI-SERVICES
// @contact.url http://www.cqiserv.com/
// @contact.email zor@cqiserv.com
func main() {
	// server start
	util.Success(" < - SERVER START - > ")
	//DB config
	defer db.SqlDB.Close()
	//ROUTER
	router := router.InitRouter()
	//PORT
	if config.GetGlobalConfig().HTTPS == true {
		router.RunTLS(":8443", "templates/server.crt", "templates/server.key")
	} else {
		if err := router.Run(":8080"); err != nil {
			util.Error(err.Error())
		}
	}
}
