package service

import (
	"cash-server/db"
	"cash-server/model"
)

//OrderAdd 新增Oder帳號
func OrderAdd(o db.Order) {
	model.OrderAdd(o)
}
