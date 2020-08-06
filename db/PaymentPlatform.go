package db

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
