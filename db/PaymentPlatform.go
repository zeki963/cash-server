package db

import (
	"cash-server/pkg/util"
	"fmt"
	"reflect"
)

//PaymentPlatform PaymentPlatform
type PaymentPlatform struct {
	Model
	PlatformAccount     string `gorm:"type:varchar(100) NOT NULL COMMENT '帳號'"`
	PlatformPassword    string `gorm:"type:varchar(100) NOT NULL COMMENT '密碼'"`
	PlatformName        string `gorm:"type:varchar(100) NOT NULL  COMMENT '別稱'"`
	PlatformGroupID     int    `gorm:"type:int NOT NULL DEFAULT '1' COMMENT '群組'"`
	PlatformEmail       string `gorm:"type:varchar(100) NOT NULL COMMENT '信箱'"`
	PlatformToken       string `gorm:"type:varchar(50) NOT NULL COMMENT '金鑰'"`
	PlatformTokenSecret string `gorm:"type:varchar(100) NOT NULL COMMENT '加密金鑰'"`
	Status              string `gorm:"type:varchar(1) DEFAULT '0' COMMENT '狀態/0禁用/1啟用'"`
}

//IPaymentPlatform IPaymentPlatform
type IPaymentPlatform interface {
	Echo() bool
	Add() error
}

//Echo 呼叫
func (i PaymentPlatform) Echo() bool {
	fmt.Println(i)
	return true
}

//Add 新增
func (i PaymentPlatform) Add() bool {
	SQLDBX.Create(&i)
	return true
}

//DBFind 新增
func (i PaymentPlatform) DBFind() (string, string) {
	b := Struct2Map(i)
	var Str string
	var columnName string
	for k, v := range b {
		if v != "" && reflect.TypeOf(v).String() == "string" {
			Str = fmt.Sprintf("%s", v)
			columnName = k
		}
	}
	//q := util.UnMarshal(columnName) + "=" + Str
	return util.UnMarshal(columnName), Str
}
