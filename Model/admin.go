package admin

import (
	"database/sql"

	//mysql
	_ "github.com/go-sql-driver/mysql"
)

//Registerserver 1
type Registerserver struct {
	id    string
	name  string
	token string
	time  string
}

func insert() {
	db, err := sql.Open("mysql", "user:password@/dbname")
	if err != nil {
		println(err.Error())
	}

	result, err := db.Exec(
		"INSERT INTO user_info (name, token , time) VALUES (?, ? , ?)",
		"syhlion",
		"18",
		"1",
	)
	println(result)
} 