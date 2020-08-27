package db

import (
	"cash-server/configs"
	"encoding/json"
	"fmt"
	"reflect"
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
	Editor    string    `gorm:"type:varchar(200) COMMENT 'Last Editor'"`
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
	if configs.GetGlobalConfig().RunMode == "debug" {
		SQLDBX.LogMode(true)
	}
	migratetable(SQLDBX)
	return nil
}

//migratetable 初始化表 自動建表
func migratetable(gdb *gorm.DB) {
	// Migrate the schema
	gdb.AutoMigrate(&LogConnect{}, &Platform{}, &PlatformGroupsAuth{}, &PlatformGroup{}, &PaymentType{}, &Order{})
	model := &PlatformGroup{}
	//檢查初始資料
	if a := gdb.Where("id = ?", 1).First(&model); a.Error != nil {
		var PlatformGroup1 PlatformGroup
		PlatformGroup1.GroupName = "nil"
		PlatformGroup1.GroupDescribe = "預設群組"
		gdb.Create(&PlatformGroup1)
		var PlatformGroup2 PlatformGroup
		PlatformGroup2.GroupName = "casino"
		PlatformGroup2.GroupDescribe = "歡樂賭爛城"
		gdb.Create(&PlatformGroup2)
	}
	model2 := &PaymentType{}
	if b := gdb.Where("id = ?", 1).First(&model2); b.Error != nil {
		var dfmycad PaymentType
		dfmycad.TypeName = "mycard"
		dfmycad.Status = "1"
		gdb.Create(&dfmycad)
	}
}

//Struct2Map Struct2Map
func Struct2Map(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Name] = v.Field(i).Interface()
	}
	return data
}

//Struct2JSON Struct2JSON
func Struct2JSON(obj interface{}) string {
	j, err := json.Marshal(obj)
	if err != nil {
		fmt.Println(err)
		return err.Error()
	}
	return string(j)
}
