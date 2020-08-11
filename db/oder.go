package db

//Oder Oder
type Oder struct {
	Model
	PaymentTypeID        int    `gorm:"type:int NOT NULL COMMENT '第3方ID'"`
	PlatformID           int    `gorm:"type:int NOT NULL COMMENT '平台 ID'"`
	OrderClientID        string `gorm:"type:varchar(50) NOT NULL COMMENT '訂單客戶ID'"`
	OrderDate            string `gorm:"type:timestamp NOT NULL COMMENT '訂單時間'"`
	OrderOriginalData    string `gorm:"type:varchar(1000) DEFAULT NULL COMMENT '訂單資料'"`
	OrderPrice           string `gorm:"type:decimal(10,2) NOT NULL COMMENT '訂單價格'"`
	RedirectURL          string `gorm:"type: varchar(1000) DEFAULT NULL COMMENT '導向網址'"`
	CallbackOriginalData string `gorm:"type:varchar(1000) DEFAULT NULL"`
	ReceivedCallbackDate string `gorm:"type:timestamp NULL DEFAULT NULL"`
	CallbackURL          string `gorm:"type:varchar(1000) DEFAULT NULL COMMENT '回調網址'"`
	PaymentID            string `gorm:"type:varchar(50) NOT NULL COMMENT '支付端訂單編號'"`
	StageType            int    `gorm:"type:int DEFAULT NULL COMMENT '1正式，0測試'"`
	Status               string `gorm:"type:varchar(1) DEFAULT '1' NOT NULL COMMENT '狀態，0為禁用，1為啟用'"`
}
