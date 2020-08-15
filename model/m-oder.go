package model

import (
	"cash-server/db"
	"cash-server/pkg/util"
	"fmt"
)

//---------------------------Oder  表單相關-------------------------------------

//OrderAdd Oder 新增
func OrderAdd(o db.Order) {
	db.SQLDBX.Create(&o)
}

//OrderSubidQuery Oder subid查詢
func OrderSubidQuery(o db.Order) db.Order {
	var model db.Order
	db.SQLDBX.Where("order_sub_id = ?", o.OrderSubID).First(&model)
	return model
}

//OrderCallbackTSave Oder 更新
func OrderCallbackTSave(o db.Order, Mrep db.Mycardresp) bool {
	db.SQLDBX.First(&o, OrderSubidQuery(o))
	fmt.Printf("%+v", o)
	o.PaymentID = Mrep.AuthCode
	o.CallbackURL = "http://test.mycard520.com.tw/MyCardPay/?AuthCode=" + Mrep.AuthCode
	o.ReceivedCallbackDate = util.GetUTCTime()
	o.CallbackOriginalData = db.Struct2JSON(Mrep)
	dbrut := db.SQLDBX.Save(o)
	return dbErrBool(dbrut)
}

//OrderCallbackFSave Oder 更新
func OrderCallbackFSave(o db.Order, Mrep db.Mycardresp) bool {
	db.SQLDBX.First(&o, OrderSubidQuery(o))
	fmt.Printf("%+v", o)
	o.PaymentID = "FAIL"
	o.ReceivedCallbackDate = util.GetUTCTime()
	o.CallbackOriginalData = db.Struct2JSON(Mrep)
	dbrut := db.SQLDBX.Save(o)
	return dbErrBool(dbrut)
}
