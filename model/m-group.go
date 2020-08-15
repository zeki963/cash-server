package model

import (
	"cash-server/db"
)

//---------------------------PaymentPlatformGroup  表單相關-------------------------------------

//GroupAdd  群組新增
func GroupAdd(g db.PaymentPlatformGroup) bool {
	dbrtu := db.SQLDBX.Create(&g)
	return dbErrBool(dbrtu)
}

//GroupExist 查詢群組存在
func GroupExist(groupid int) db.PaymentPlatformGroup {
	var model db.PaymentPlatformGroup
	db.SQLDBX.Where(" id = ?", groupid).First(&model)
	return model
}

//GroupOrderUpdate 更新群組訂單單號
func GroupOrderUpdate(groupid int) bool {
	var model db.PaymentPlatformGroup
	dbrtu := db.SQLDBX.Where(" id = ?", groupid).First(&model).Update("group_order", model.GroupOrder+1)
	return dbErrBool(dbrtu)
}

//GroupNameCheck 檢查群組Name存在
func GroupNameCheck(name string) bool {
	var model db.PaymentPlatformGroup
	dbrtu := db.SQLDBX.Where(" group_name = ?", name).First(&model)
	return dbErrBool(dbrtu)
}

//---------------------------PaymentPlatformGroupsAuth  表單相關-------------------------------------

//GroupAuthAdd  群組權限新增
func GroupAuthAdd(ga db.PaymentPlatformGroupsAuth) error {
	db.SQLDBX.Create(&ga)
	return nil
}

//PlatformGroupAuthQuery 查詢帳號權限開通狀態
func PlatformGroupAuthQuery(groupid string, typeid string) db.PaymentPlatformGroupsAuth {
	var model db.PaymentPlatformGroupsAuth
	db.SQLDBX.Where(" group_id= ? AND type_id=?", groupid, typeid).First(&model)
	return model
}
