package database

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var SqlDB *sql.DB

//Dbcannot DB連線
func Dbcannot(DBInfo string) {
	var err error
	SqlDB, err = sql.Open("mysql", DBInfo)
	if err != nil {
		log.Fatal(err.Error())
	}
	err = SqlDB.Ping()
	if err != nil {
		log.Fatal(err.Error())
	}
}
