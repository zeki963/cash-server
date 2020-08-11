package db

import (
	"cash-server/configs"
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
	gdb.AutoMigrate(&LogConnect{}, &PaymentPlatform{}, &PaymentPlatformGroupsAuth{}, &PaymentPlatformGroup{}, &PaymentType{}, &Oder{})
	model := &PaymentPlatformGroup{}
	//檢查初始資料
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
	model2 := &PaymentType{}
	if b := gdb.Where("id = ?", 1).First(&model2); b.Error != nil {
		var dfmycad PaymentType
		dfmycad.TypeName = "mycard"
		dfmycad.PaymentMemberID = "1234"
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
