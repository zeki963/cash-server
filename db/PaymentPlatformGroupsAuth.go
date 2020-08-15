package db

//PaymentPlatformGroupsAuth PaymentPlatformGroupsAuth
type PaymentPlatformGroupsAuth struct {
	Model
	GroupID string `gorm:"type:int(11) NOT NULL COMMENT '群組ID'"`
	TypeID  string `gorm:"type:int(11) NOT NULL COMMENT '第3方ID' "`
}
