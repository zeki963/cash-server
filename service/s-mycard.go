package service

import (
	"cash-server/db"
	"cash-server/model"
)

//PlatformQueryStatus 查詢狀態
func PlatformQueryStatus(p db.PaymentPlatform) bool {
	a := model.PlatformQueryExist(p)
	if a.Status == "1" {
		return true
	}
	return false
}

//PlatformGroupAuthQuery 查詢群組狀態
func PlatformGroupAuthQuery(groupid string, typeid string) bool {
	m := model.PlatformGroupAuthQuery(groupid, typeid)
	if m.GroupID != "" && m.TypeID != "" {
		return true
	}
	return false
}
