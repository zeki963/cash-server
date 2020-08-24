package main

import (
	"cash-server/configs"
	"cash-server/db"
	casinogrpc "cash-server/grpc/casino"
	"cash-server/pkg/util"
	"cash-server/router"
	"cash-server/service"
	"strconv"
)

// @title  金流SERVER API
// @version 2020.08
// @description This is a sample API server .
// @contact.name  CQI-SERVICES
// @contact.url https://www.cqiserv.com/
// @contact.email zor@cqiserv.com
func main() {
	// <== 測試CODE == >
	util.Test(" < - TEST  - > ")
	if configs.GetGlobalConfig().RunMode == "test" {
		//casinoUser := casinogrpc.VetifyUserGUID("20312543")
		//casinogrpc.SendItemBuy(casinoUser, "com.cqicasino.happycity.tier1")
		casinogrpc.SendItemResult(service.OrderFind(db.Order{OrderSubID: "casinox000000111"}))
		casinogrpc.VetifyItem("com.cqicasino.happycity.tier1", 300)
		//casinogrpc.SendItemResult("70ccafc9-0c66-47a1-af93-a075a9f08ce9", "com.cqicasino.happycity.tier1")
	}
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
