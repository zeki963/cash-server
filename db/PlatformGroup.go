package db

//PlatformGroup PlatformGroup
type PlatformGroup struct {
	Model
	GroupName     string `gorm:"type:varchar(20) NOT NULL  COMMENT '群組名稱'"`
	GroupDescribe string `gorm:"type:varchar(45)  COMMENT '群組詳細'"`
	GroupOrder    int    `gorm:"type:int(11) NOT NULL DEFAULT 0 COMMENT '群組單數'"`
	Status        string `gorm:"type:varchar(1) DEFAULT '0' COMMENT '狀態/0禁用/1啟用'"`
}
