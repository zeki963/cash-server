package service

import (
	"cash-server/db"
	"cash-server/model"
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
