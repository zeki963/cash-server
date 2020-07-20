package model

import (
	//mysql
	db "cash-server/database"
	"cash-server/pkg/util"
	"fmt"
	"log"
)

//Registerserver 1
type Registerserver struct {
	id    string
	name  string
	token string
	time  string
}

//InsertServer 測試寫入DB
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

//InsertLog 寫入紀錄
func InsertLog(name string, token string, time string) error {

}

//UserRegister User 註冊
func UserRegister(account string, pw string, name string, token string, tokenSecret string, time string) error {
	if _, err := db.SqlDB.Exec(
		"INSERT INTO payment_user (user_account, user_password, user_name, user_token, user_token_secret, status, create_date, update_date) VALUES (?,?,?,?,?,0,?,?)",
		account, pw, name, token, tokenSecret, time, time,
	); err != nil {
		util.Error("DB INSERT fail")
		util.Error(err.Error())
		return err
	}
	return nil
}

//UserQuery User Query 查詢
func UserQuery(account string) {

	rows, err := db.SqlDB.Query("SELECT user_password FROM payment_user WHERE user_account = ?", account)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		var password string
		if err := rows.Scan(&password); err != nil {
			log.Fatal(err)
		}
		fmt.Println("密碼：", password)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}
