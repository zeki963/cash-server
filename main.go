package main

import (
	"cash-server/configs"
	"cash-server/db"
	"cash-server/pkg/util"
	"cash-server/router"
	"strconv"
)

// @title  金流SERVER API
// @version 2020.08
// @description This is a sample API server .
// @contact.name  CQI-SERVICES
// @contact.url https://www.cqiserv.com/
// @contact.email zor@cqiserv.com
func main() {
	// server start
	util.Info(" < - SERVER START - > ")
	//DB 連線關閉
	defer db.SQLDBX.Close()
	//ROUTER
	r := router.InitRouter()
	//PORT
	if configs.GetGlobalConfig().HTTP.HTTPS == true {
		r.RunTLS(":8443", "templates/server.crt", "templates/server.key")
	} else {
		if err := r.Run(":" + strconv.Itoa(configs.GetGlobalConfig().HTTP.HTTPPort)); err != nil {
			util.Error(err.Error())
		}
	}
}

func init() {
	util.Info(" < - SERVER INIT - > ")
	configs.LoadGlobalConfig("")
	util.Success("[MODE] > " + configs.GetGlobalConfig().RunMode)
	util.Info("[GRPC-Casino] > " + configs.GetGlobalConfig().Casino.Alphaip)
	util.Info(" < - MyDB INIT - >")
	if err := db.Initgorm(); err != nil {
		util.Error("! DB Connect ERROR !")
		util.Error(err.Error())
		panic("error!")
	} else {
		util.Success("[DB Host] > " + configs.GetGlobalConfig().MySQL.Host)
	}
}
