package db

//PlatformGroupsAuth PlatformGroupsAuth
type PlatformGroupsAuth struct {
	Model
	GroupID string `gorm:"type:int(11) NOT NULL COMMENT '群組ID'"`
	TypeID  string `gorm:"type:int(11) NOT NULL COMMENT '第3方ID' "`
	Status  string `gorm:"type:varchar(1) DEFAULT '0' COMMENT '狀態/0禁用/1啟用'"`
}
