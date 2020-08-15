package service

import (
	"cash-server/db"
	"cash-server/model"
	"strconv"
)

//PaymentPlatformAdd 新增平台帳號
func PaymentPlatformAdd(p db.PaymentPlatform) bool {
	if model.PaymentPlatformCheck(p.PlatformAccount) {
		return false
	}
	model.PaymentPlatformAdd(p)
	return true
}

//LogConnectAdd 紀錄連線
func LogConnectAdd(p db.LogConnect) {
	model.LogConnectAdd(p)
}

//PaymentPlatformFind 查整組資料
func PaymentPlatformFind(p db.PaymentPlatform) (all db.PaymentPlatform) {
	all = model.PlatformQueryExist(p)
	return all
}

//PlatformQueryStatus 查詢狀態
func PlatformQueryStatus(p db.PaymentPlatform) bool {
	a := model.PlatformQueryExist(p)
	if a.Status == "1" {
		return true
	}
	return false
}

//PlatformGroupAuthQuery 查詢群組狀態
func PlatformGroupAuthQuery(p db.PaymentPlatform, typeid string) bool {
	p2 := model.PlatformQueryExist(p)
	m := model.PlatformGroupAuthQuery(strconv.Itoa(p2.PlatformGroupID), typeid)
	if m.GroupID != "" && m.TypeID != "" {
		return true
	}
	return false
}
