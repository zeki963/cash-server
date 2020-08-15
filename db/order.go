package db

import "time"

//Order Order
type Order struct {
	Model
	PaymentTypeID        int       `gorm:"type:int NOT NULL COMMENT '第3方ID'"`
	PlatformID           int       `gorm:"type:int NOT NULL COMMENT '平台 ID'"`
	OrderClientID        string    `gorm:"type:varchar(50) NOT NULL COMMENT '訂單客戶ID'"`
	OrderDate            time.Time `gorm:"type:timestamp NOT NULL COMMENT '訂單時間'"`
	OrderItemID          string    `gorm:"type:int NOT NULL COMMENT '訂單商品項目'"`
	OrderItemPrice       string    `gorm:"type:decimal(10,2) NOT NULL COMMENT '訂商品價格'"`
	OrderSubID           string    `gorm:"type:varchar(100) NOT NULL COMMENT '商品子流水編號'"`
	OrderOriginalData    string    `gorm:"type:varchar(1000) DEFAULT NULL COMMENT '訂單資料'"`
	CallbackOriginalData string    `gorm:"type:varchar(1000) DEFAULT NULL COMMENT '收到原始資料'"`
	ReceivedCallbackDate time.Time `gorm:"type:timestamp COMMENT '收到時間'"`
	CallbackURL          string    `gorm:"type:varchar(1000) DEFAULT NULL COMMENT '回調網址'"`
	PaymentID            string    `gorm:"type:varchar(1000) NOT NULL COMMENT '支付端訂單編號'"`
	StageType            int       `gorm:"type:int DEFAULT NULL COMMENT '1正式，0測試'"`
	Status               string    `gorm:"type:varchar(1) DEFAULT '1' NOT NULL COMMENT '狀態，0為禁用，1為啟用'"`
}
