package db

//PaymentType PaymentType
type PaymentType struct {
	Model
	PaymentMemberID string `gorm:"int(255) NOT NULLL"`
	TypeName        string `gorm:"type:varchar(100) NOT NULL"`
	//TypeToken       string `gorm:"type:varchar(50) NOT NULL COMMENT '金鑰'"`
	//TypeTokenSecret string `gorm:"type:varchar(100) NOT NULL COMMENT '加密金鑰'"`
	Status string `gorm:"type:varchar(1) DEFAULT '1' NOT NULL COMMENT '狀態，0為禁用，1為啟用'"`
}
