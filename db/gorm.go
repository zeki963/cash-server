package db

import (
	"cash-server/configs"
	"time"

	"github.com/jinzhu/gorm"
	//sql
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

//Model 預設模組
type Model struct {
	ID        uint      `gorm:"primary_key; AUTO_INCREMENT"`
	CreatedAt time.Time `gorm:"type:timestamp;NOT NULL"`
	UpdatedAt time.Time `gorm:"type:timestamp;NOT NULL"`
}

//LogConnect  連線紀錄
type LogConnect struct {
	Model
	StatusCode  int
	LatencyTime string
	ClientIP    string `gorm:"type:varchar(15)"`
	ReqMethod   string `gorm:"type:varchar(100)"`
	ReqURL      string `gorm:"type:varchar(100)"`
	Reqbody     string `gorm:"type:varchar(255)"`
	Reqheader   string `gorm:"type:text"`
}

//PaymentPlatform PaymentPlatform
type PaymentPlatform struct {
	Model
	PlatformAccount     string `gorm:"type:varchar(100) NOT NULL COMMENT '帳號'"`
	PlatformPassword    string `gorm:"type:varchar(100) NOT NULL COMMENT '密碼'"`
	PlatformName        string `gorm:"type:varchar(100) NOT NULL"`
	PlatformGroupID     int    `gorm:"type:int NOT NULL DEFAULT '1' COMMENT '1 為空群組'"`
	PlatformEmail       string `gorm:"type:varchar(100) NOT NULL"`
	PlatformToken       string `gorm:"type:varchar(50) NOT NULL COMMENT '金鑰'"`
	PlatformTokenSecret string `gorm:"type:varchar(100) NOT NULL COMMENT '加密金鑰'"`
	Status              string `gorm:"type:varchar(1) DEFAULT '0' COMMENT '狀態，0為禁用，1為啟用'"`
}

//PaymentPlatformGroup PaymentPlatformGroup
type PaymentPlatformGroup struct {
	Model
	GroupName     string `gorm:" varchar(20) NOT NULL COMMENT '群組名稱'"`
	GroupDescribe string `gorm:"varchar(45)  COMMENT '詳細'"`
}

//PaymentPlatformGroupAuth PaymentPlatformGroupAuth
type PaymentPlatformGroupAuth struct {
	Model
	GroupID string `gorm:"int(10) NOT NULL"`
	TypeID  string `gorm:"int(10) NOT NULLL"`
}

//SQLDBX  供外部呼叫
var SQLDBX *gorm.DB

//GetDB  供外部呼叫
func GetDB() *gorm.DB {
	return SQLDBX
}

//Initgorm 初始化gorm
func Initgorm() error {
	var err error
	SQLDBX, err = gorm.Open("mysql", configs.MySQL.DSN(configs.GetGlobalConfig().MySQL))
	if err != nil {
		println(err.Error())
		return err
	}
	SQLDBX.DB().SetConnMaxLifetime(1 * time.Second)
	SQLDBX.DB().SetMaxIdleConns(20)
	SQLDBX.DB().SetMaxOpenConns(2000)
	if configs.GetGlobalConfig().RunMode != "release" {
		SQLDBX.LogMode(true)
	}
	migratetable(SQLDBX)
	return nil
}

//migratetable 初始化表 自動建表
func migratetable(gdb *gorm.DB) {
	// Migrate the schema
	gdb.AutoMigrate(&LogConnect{}, &PaymentPlatform{}, &PaymentPlatformGroupAuth{}, &PaymentPlatformGroup{})
	model := &PaymentPlatformGroup{}
	if a := gdb.Where("id = ?", 1).First(&model); a.Error != nil {
		var PaymentPlatformGroup1 PaymentPlatformGroup
		PaymentPlatformGroup1.GroupName = "nil"
		PaymentPlatformGroup1.GroupDescribe = "預設群組"
		gdb.Create(&PaymentPlatformGroup1)
		var PaymentPlatformGroup2 PaymentPlatformGroup
		PaymentPlatformGroup2.GroupName = "casino"
		PaymentPlatformGroup2.GroupDescribe = "歡樂賭爛城"
		gdb.Create(&PaymentPlatformGroup2)
	}
}
