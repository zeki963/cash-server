package model

import (
	//mysql
	db "cash-server/database"
	"cash-server/pkg/util"
)

//Registerserver 1
type Registerserver struct {
	id    string
	name  string
	token string
	time  string
}

//InsertServer 寫入DB
func InsertServer(name string, token string, time string) error {
	if _, err := db.SqlDB.Exec(
		"INSERT INTO registerserver (name, token , time) VALUES (?, ? , ?)", name, token, time,
	); err != nil {
		util.Error("DB INSERT fail")
		util.Error(err.Error())
		return err
	}
	return nil
}
