package db

//PaymentType PaymentType
type PaymentType struct {
	Model
	TypeName string `gorm:"type:varchar(100) NOT NULL  COMMENT '第3方金流名稱'"`
	Status   string `gorm:"type:varchar(1) DEFAULT '1' NOT NULL COMMENT '狀態，0為禁用，1為啟用'"`
}
