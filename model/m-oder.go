package model

import (
	"cash-server/db"
	"cash-server/pkg/util"
	"fmt"
)

//---------------------------Order  表單相關-------------------------------------

//OrderAdd Order 新增
func OrderAdd(o db.Order) {
	db.SQLDBX.Create(&o)
}

//OrderSubidQuery Order subid查詢
func OrderSubidQuery(o db.Order) db.Order {
	var model db.Order
	db.SQLDBX.Where("order_sub_id = ?", o.OrderSubID).First(&model)
	return model
}

//OrderAuthTSave Order 更新認証成功
func OrderAuthTSave(o db.Order, Mrep db.Mycardresp) bool {
	db.SQLDBX.First(&o, OrderSubidQuery(o))
	util.Test(fmt.Sprintf("交易單明細 : %+v", o))
	o.PaymentID = Mrep.TradeSeq
	o.CallbackURL = "http://test.mycard520.com.tw/MyCardPay/?AuthCode=" + Mrep.AuthCode
	o.ReceivedCallbackDate = util.GetUTCTime()
	o.CallbackOriginalData = db.Struct2JSON(Mrep)
	dbrut := db.SQLDBX.Save(o)
	return dbErrBool(dbrut)
}

//OrderAuthFSave Order 更新認証失敗
func OrderAuthFSave(o db.Order, Mrep db.Mycardresp) bool {
	db.SQLDBX.First(&o, OrderSubidQuery(o))
	util.Test(fmt.Sprintf("交易單明細 : %+v", o))
	o.PaymentID = "FAIL"
	o.ReceivedCallbackDate = util.GetUTCTime()
	o.CallbackOriginalData = db.Struct2JSON(Mrep)
	dbrut := db.SQLDBX.Save(o)
	return dbErrBool(dbrut)
}

//OrderCallbackTSave Order 更新 回應成功
func OrderCallbackTSave(o db.Order, callbackstring string) bool {
	db.SQLDBX.First(&o, OrderSubidQuery(o))
	util.Test(fmt.Sprintf("交易單明細 : %+v", o))
	o.ReceivedCallbackDate = util.GetUTCTime()
	dbrut := db.SQLDBX.Save(o)
	return dbErrBool(dbrut)
}

//OrderQueryExist 查詢帳號存在
func OrderQueryExist(data db.Order) db.Order {
	//model := &db.Order{}
	var model db.Order
	k, v := data.DBFind()
	db.SQLDBX.Where(k+"= ?", v).First(&model)
	return model
}

//OrderQueryInfoAllJSON  查詢ALL帳號資料
func OrderQueryInfoAllJSON() []db.Order {
	var Orders []db.Order
	db.SQLDBX.Find(&Orders)
	util.Test(fmt.Sprint((&Orders)))
	return Orders
}
