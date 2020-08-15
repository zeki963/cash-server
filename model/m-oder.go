package model

import (
	"cash-server/db"
)

//---------------------------Oder  表單相關-------------------------------------

//OrderAdd Oder 新增
func OrderAdd(o db.Order) {
	db.SQLDBX.Create(&o)
}

//OrderSave Oder 更新
func OrderSave(o db.Order) bool {
	dbrut := db.SQLDBX.Save(o)
	return dbErrBool(dbrut)
}
