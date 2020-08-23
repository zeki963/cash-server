package db

import (
	"cash-server/pkg/util"
	"fmt"
	"reflect"
	"time"
)

//Order Order
type Order struct {
	Model
	PaymentTypeID        int       `gorm:"type:int NOT NULL COMMENT '第3方ID'"`
	PlatformID           int       `gorm:"type:int NOT NULL COMMENT '平台 ID'"`
	OrderClientID        string    `gorm:"type:varchar(50) NOT NULL COMMENT '訂單客戶ID'"`
	OrderDate            time.Time `gorm:"type:timestamp NOT NULL COMMENT '訂單時間'"`
	OrderItemID          string    `gorm:"type:varchar(100) NOT NULL COMMENT '訂單商品項目'"`
	OrderItemPrice       string    `gorm:"type:decimal(10,2) NOT NULL COMMENT '訂商品價格'"`
	OrderSubID           string    `gorm:"type:varchar(100) NOT NULL COMMENT '商品子流水編號'"`
	OrderOriginalData    string    `gorm:"type:varchar(1000) DEFAULT NULL COMMENT '訂單資料'"`
	OrderGameSubID       int32     `gorm:"type:varchar(100) DEFAULT NULL COMMENT '遊戲資料庫編號'"`
	CallbackOriginalData string    `gorm:"type:varchar(1000) DEFAULT NULL COMMENT '收到原始資料'"`
	ReceivedCallbackDate time.Time `gorm:"type:timestamp COMMENT '收到時間'"`
	PaymentAuth          string    `gorm:"type:varchar(1000) DEFAULT NULL COMMENT '交易認証編號'"`
	PaymentID            string    `gorm:"type:varchar(1000) NOT NULL COMMENT '支付端訂單編號'"`
	StageType            int       `gorm:"type:int DEFAULT NULL COMMENT '1正式，0測試'"`
	Status               string    `gorm:"type:varchar(1) DEFAULT '0' NOT NULL COMMENT '狀態，0為禁用，1為啟用'"`
	MycardTradeNo        string    `gorm:"type:varchar(1000) DEFAULT  NULL COMMENT '通路編號'"`
	PaymentConfirm       string    `gorm:"type:varchar(1) DEFAULT '0' NOT NULL COMMENT '狀態，0為禁用，1為啟用'"`
}

//IOrder IOrder
type IOrder interface {
	DBFind() (string, string)
}

//DBFind 新增
func (i Order) DBFind() (string, string) {
	b := Struct2Map(i)
	var Str string
	var columnName string
	for k, v := range b {
		if v != "" && reflect.TypeOf(v).String() == "string" {
			Str = fmt.Sprintf("%s", v)
			columnName = k
		}
	}
	return util.UnMarshal(columnName), Str
}
