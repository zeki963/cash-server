package service

import (
	"cash-server/db"
	"cash-server/model"
)

//OrderAdd 新增Oder
func OrderAdd(o db.Order) {
	o.Status = "1"
	model.OrderAdd(o)
}

//OrderSave 更新Oder
func OrderSave(o db.Order, Mrep db.Mycardresp) {
	if Mrep.ReturnCode == "1" {
		model.OrderCallbackTSave(o, Mrep)
	}
	model.OrderCallbackFSave(o, Mrep)
}
