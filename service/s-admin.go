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

//PlatformAdd 新增平台帳號
func PlatformAdd(p db.Platform) bool {
	if model.PlatformCheck(p.PlatformAccount) {
		return false
	}
	model.PlatformAdd(p)
	return true
}

//PlatformFind 查整組資料
func PlatformFind(p db.Platform) (all db.Platform) {
	all = model.PlatformQueryExist(p)
	return all
}

//PlatformQueryOne 查詢存在
func PlatformQueryOne(p db.Platform) bool {
	if model.PlatformQueryOne(p) == "" {
		util.Test("PlatformAccount 不存在")
		return false
	}
	util.Test("PlatformAccount 存在")
	return true
}

//PlatformStatusEnable 帳號狀態轉啟用
func PlatformStatusEnable(p db.Platform) bool {
	a := model.PlatformQueryExist(p)
	if a.Status == "1" {
		util.Test("PlatformStatusEnable 帳號狀態 啟用了")
		return false
	}
	util.Test("PlatformStatusEnable 帳號狀態轉啟用")
	model.PlatformStatusEnableUpdata(p)
	return true
}

//PlatformStatusDisable 帳號狀態轉啟用
func PlatformStatusDisable(p db.Platform) bool {
	a := model.PlatformQueryExist(p)
	if a.Status == "0" {
		util.Test("PlatformStatusEnable 帳號狀態 關閉了")
		return false
	}
	util.Test("PlatformStatusEnable 帳號狀態轉未啟用")
	model.PlatformStatusDisableUpdata(p)
	return true
}

//PlatformQueryStatus 查詢狀態
func PlatformQueryStatus(p db.Platform) bool {
	a := model.PlatformQueryExist(p)
	if a.Status != "1" {
		return false
	}
	return true
}

//PlatformGroupAuthQuery 查詢群組狀態
func PlatformGroupAuthQuery(p db.Platform, typeid string) bool {
	p2 := model.PlatformQueryExist(p)
	m := model.PlatformGroupAuthQuery(strconv.Itoa(p2.PlatformGroupID), typeid)
	if m.GroupID != "" && m.TypeID != "" {
		return true
	}
	return false
}
