package service

import (
	"cash-server/db"
	"cash-server/model"
	"cash-server/pkg/util"
	"strconv"
)

//LogConnectAdd 紀錄連線
func LogConnectAdd(p db.LogConnect) {
	model.LogConnectAdd(p)
}

//PaymentPlatformAdd 新增平台帳號
func PaymentPlatformAdd(p db.PaymentPlatform) bool {
	if model.PaymentPlatformCheck(p.PlatformAccount) {
		return false
	}
	model.PaymentPlatformAdd(p)
	return true
}

//PaymentPlatformFind 查整組資料
func PaymentPlatformFind(p db.PaymentPlatform) (all db.PaymentPlatform) {
	all = model.PlatformQueryExist(p)
	return all
}

//PlatformQueryOne 查詢存在
func PlatformQueryOne(p db.PaymentPlatform) bool {
	if model.PlatformQueryOne(p) == "" {
		util.Test("PlatformAccount 不存在")
		return false
	}
	util.Test("PlatformAccount 存在")
	return true
}

//PlatformQueryStatus 查詢狀態
func PlatformQueryStatus(p db.PaymentPlatform) bool {
	a := model.PlatformQueryExist(p)
	if a.Status != "1" {
		return false
	}
	return true
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
