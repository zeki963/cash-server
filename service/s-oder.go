package service

import (
	"cash-server/db"
	"cash-server/model"
)

//OrderAdd 新增Oder
func OrderAdd(o db.Order) {
	model.OrderAdd(o)
}

//OrderUpdate 更新Oder
func OrderUpdate(o db.Order) {
	model.OrderAdd(o)
}
