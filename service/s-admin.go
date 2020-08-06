package service

import (
	"cash-server/db"
	"cash-server/model"
)

//PaymentPlatformAdd 新增
func PaymentPlatformAdd(p db.PaymentPlatform) bool {
	if model.PaymentPlatformCheck(p.PlatformAccount) {
		return false
	}
	model.PaymentPlatformAdd(p)
	return true
}
