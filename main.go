package main

import (
	db "cash-server/database"
	encode "cash-server/pkg/encode"
	flag "cash-server/pkg/flag"
	util "cash-server/pkg/util"
	router "cash-server/routers"
)

func init() {
	encode.Md5("123")
	println(dblink)
	flag.SettestMod(true)
	flag.StartFlag()
}

func main() {
	// server start
	util.Success(" < - SERVER START - > ")
	//DB config
	defer db.SqlDB.Close()
	//ROUTER
	router := router.InitRouter()
	//PORT
	//router.RunTLS(":8443", "templates/server.crt", "templates/server.key")
	if err := router.Run(":8080"); err != nil {
		util.Error(err.Error())
	}
}
