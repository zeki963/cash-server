package model

import "cash-server/db"

//---------------------------Oder  表單相關-------------------------------------

//OrderAdd Oder 新增
func OrderAdd(o db.Order) error {
	db.SQLDBX.Create(&o)
	return nil
}
