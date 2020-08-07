package db

//PaymentPlatformGroup PaymentPlatformGroup
type PaymentPlatformGroup struct {
	Model
	GroupName     string `gorm:" varchar(20) NOT NULL COMMENT '群組名稱'"`
	GroupDescribe string `gorm:"varchar(45)  COMMENT '詳細'"`
}
