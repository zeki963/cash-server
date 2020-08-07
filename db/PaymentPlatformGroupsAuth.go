package db

//PaymentPlatformGroupsAuth PaymentPlatformGroupsAuth
type PaymentPlatformGroupsAuth struct {
	Model
	GroupID string `gorm:"int(10) NOT NULL"`
	TypeID  string `gorm:"int(10) NOT NULLL"`
}
