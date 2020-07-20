package database

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

//SqlDB DB控制項
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
	SqlDB.SetConnMaxLifetime(100 * time.Second) //最大連接周期，超時close
	SqlDB.SetMaxOpenConns(100)                  //設置最大連接數
}
