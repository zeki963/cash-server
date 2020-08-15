package db

//PlatformGroupsAuth PlatformGroupsAuth
type PlatformGroupsAuth struct {
	Model
	GroupID string `gorm:"type:int(11) NOT NULL COMMENT '群組ID'"`
	TypeID  string `gorm:"type:int(11) NOT NULL COMMENT '第3方ID' "`
}
