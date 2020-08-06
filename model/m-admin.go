package model

import (
	"cash-server/db"
	"cash-server/pkg/util"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

//Registerserver 舊版本
type Registerserver struct {
	id    string
	name  string
	token string
	time  string
}

//Platform struct
type Platform struct {
	Platformid          int    `json:"platform_id"`
	Platformaccount     string `json:"platform_account"`
	Platformpassword    string `json:"platform_password"`
	Platformname        string `json:"platform_name"`
	Platformgroupid     string `json:"platform_group_id"`
	PlatformEmail       string `json:"platform_email"`
	Platformtoken       string `json:"platform_token"`
	Platformtokensecret string `json:"platform_token_secret"`
	Status              string `json:"status"`
	Createdate          string `json:"create_date"`
	Updatedate          string `json:"update_date"`
}

//---------------------------payment_platform  表單相關-------------------------------------

//PaymentPlatformAdd Platform 註冊
func PaymentPlatformAdd(p db.PaymentPlatform) error {
	db.SQLDBX.Create(&p)
	return nil
}

//PaymentPlatformCheck Platform 查詢
func PaymentPlatformCheck(account string) bool {
	model := &db.PaymentPlatform{}
	a := db.SQLDBX.Where("platform_account = ?", account).First(&model)
	if a.Error != nil {
		return false
	}
	return true
}

//PlatformQueryExist 查詢帳號存在
func PlatformQueryExist(account string) string {
	model := &db.PaymentPlatform{}
	a := db.SQLDBX.Where("platform_account = ?", account).First(&model)
	if a.Error != nil {
		return "err"
	}
	return model.PlatformPassword
}

//PlatformQueryStatus 查詢帳號開通狀態
func PlatformQueryStatus(account string, password string) string {
	model := &db.PaymentPlatform{}
	a := db.SQLDBX.Where("platform_account = ?", account).First(&model)
	if a.Error != nil {
		return "err"
	}
	log.Println("PlatformQueryStatus 查詢帳號開通狀態 -> 帳號：", account, "狀態：", model.Status)
	return model.Status
}

//PlatformTokenQueryStatus 查詢帳號開通狀態
func PlatformTokenQueryStatus(token string) string {
	model := &db.PaymentPlatform{}
	a := db.SQLDBX.Where("platform_token = ?", token).First(&model)
	if a.Error != nil {
		return "err"
	}
	log.Println("PlatformTokenQueryStatus 查詢帳號狀態 -> token：", token)
	return model.Status
}

// //PlatformQueryInfo  查詢帳號詳細資料 乖 用id查
// func PlatformQueryInfo(taskID string) *db.PaymentPlatform {
// 	model := &db.PaymentPlatform{}
// 	a := db.SQLDBX.Where("platform_id  = ?", taskID).First(&model)
// 	if a.Error != nil {
// 		return &db.PaymentPlatform{}
// 	}
// 	fmt.Println("PlatformQueryInfo 查詢帳號資料 ", taskID)
// 	return model
// }

//PlatformQueryInfo  查詢帳號詳細資料 乖 用id查
func PlatformQueryInfo(taskID string) string {
	sql := "SELECT * FROM payment_platform where platform_id =?"
	rows, err := db.GetJSON(sql, taskID)
	fmt.Println("PlatformQueryInfo 查詢帳號資料 ", rows)
	if err != nil {
		//log.Fatal(err)
		util.Error(err.Error())
		return "err"
	}
	return rows
}

//PlatformQueryInfoAll  查詢ALL帳號資料
func PlatformQueryInfoAll() {
	var platforms []Platform
	rows, err := db.SqlDB.Query("select * from payment_platform")
	if err != nil {
		//log.Fatal(err)
		util.Error(err.Error())
	}
	for rows.Next() {
		var platformid int
		var platformaccount, platformpassword, platformname, platformgroupid, platformemail, platformtoken, platformtokensecret, status, createdate, updatedate string
		rows.Scan(&platformid, &platformaccount, &platformpassword, &platformname, &platformgroupid, &platformemail, &platformtoken, &platformtokensecret, &status, &createdate, &updatedate)
		platforms = append(platforms, Platform{platformid, platformaccount, platformpassword, platformname, platformgroupid, platformemail, platformtoken, platformtokensecret, status, createdate, updatedate})
	}
	platformsBytes, _ := json.Marshal(&platforms)
	fmt.Println(string(platformsBytes))
}

//PlatformQueryInfodata  查詢帳號資料
func PlatformQueryInfodata(token string) Platform {
	var platforms Platform
	rows, err := db.SqlDB.Query("select * from payment_platform WHERE platform_token=? ", token)
	if err != nil {
		//log.Fatal(err)
		util.Error(err.Error())
	}
	for rows.Next() {
		rows.Scan(&platforms.Platformid, &platforms.Platformaccount, &platforms.Platformpassword, &platforms.Platformname, &platforms.Platformgroupid, &platforms.PlatformEmail, &platforms.Platformtoken, &platforms.Platformtokensecret, &platforms.Status, &platforms.Createdate, &platforms.Updatedate)
		//fmt.Printf("%+v", platforms)
	}
	return platforms
}

//   ---------------------------payment_platform_group_auth  表單相關-------------------------------------

//PlatformQueryGroupAuth 查詢帳號資料
//[typeid 1=mycard]
func PlatformQueryGroupAuth(typeid string, groupid string) bool {
	rows, err := db.SqlDB.Query("select * from payment_platform_group_auth WHERE group_id=? and type_id=? ", groupid, typeid)
	if err != nil {
		//log.Fatal(err)
		util.Error(err.Error())
		return false
	}
	for rows.Next() {
		//fmt.Printf("%+v", platforms)
	}
	return true
}

// //PlatformTokenQueryStatus 查詢帳號開通狀態
// func PlatformQueryGroupAuth(typeid string, groupid string) bool {
// 	model := &db.PaymentPlatform{}
// 	a := db.SQLDBX.Where(" group_id=? and type_id=? = ?", groupid, typeid).First(&model)
// 	if a.Error != nil {
// 		return false
// 	}
// 	return true
// }

//   ---------------------------log_connect  表單相關-------------------------------------

//InsertLogConnect 寫入紀錄
func InsertLogConnect(statusCode int, latencyTime time.Duration, clientIP, reqMethod string, reqURL string, time string) error {
	if _, err := db.SqlDB.Exec(
		"INSERT INTO log_connect (statusCode, latencyTime, clientIP, reqMethod, reqURL, create_date,update_date) VALUES (?,?,?,?,?,?,?)",
		statusCode, latencyTime, clientIP, reqMethod, reqURL, time, time); err != nil {
		util.Error("DB INSERT fail")
		util.Error(err.Error())
		return err
	}
	return nil
}
