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
	util.Test(fmt.Sprintf("交易單subid查詢 : %+v", o.OrderSubID))
	var model db.Order
	db.SQLDBX.Where("order_sub_id = ?", o.OrderSubID).First(&model)
	return model
}

//OrderAuthTSave Order 更新認証成功
func OrderAuthTSave(o db.Order, Mrep db.Mycardresp, toServerVal string) bool {
	util.Test(fmt.Sprintf("交易單明細 : %+v", o))
	db.SQLDBX.First(&o, OrderSubidQuery(o))
	o.PaymentID = Mrep.TradeSeq
	o.OrderOriginalData = toServerVal
	o.PaymentAuth = Mrep.AuthCode
	o.ReceivedCallbackDate = util.GetUTCTime()
	o.CallbackOriginalData = db.Struct2JSON(Mrep)
	dbrut := db.SQLDBX.Save(o)
	return dbErrBool(dbrut)
}

//OrderAuthFSave Order 更新認証失敗
func OrderAuthFSave(o db.Order, Mrep db.Mycardresp, toServerVal string) bool {
	util.Test(fmt.Sprintf("交易單明細 : %+v", o))
	db.SQLDBX.First(&o, OrderSubidQuery(o))
	o.PaymentID = "FAIL"
	o.OrderOriginalData = toServerVal
	o.ReceivedCallbackDate = util.GetUTCTime()
	o.CallbackOriginalData = db.Struct2JSON(Mrep)
	dbrut := db.SQLDBX.Save(o)
	return dbErrBool(dbrut)
}

//MycardOrderCallbackSave  Mycard Order 更新 回應成功 3.2
func MycardOrderCallbackSave(o db.Order, callbackform *db.OrderMycard) bool {
	util.Test(fmt.Sprintf("MycardOrderCallbackSave Order 回應成功 3.2 交易單明細 : %+v", o))
	db.SQLDBX.First(&o, OrderSubidQuery(o))
	o.ReceivedCallbackDate = util.GetUTCTime()
	o.MycardTradeNo = callbackform.MyCardTradeNo
	dbrut := db.SQLDBX.Save(o)
	return dbErrBool(dbrut)
}

//MycardOrderAuthGet  Mycard Order  查詢Auth
func MycardOrderAuthGet(o db.Order) string {
	util.Test(fmt.Sprintf("MycardOrderAuthGet  查詢Auth : %+v", o))
	db.SQLDBX.First(&o, OrderSubidQuery(o))
	return o.PaymentAuth
}

//MycardTradeQuery Order 驗證 3.3
func MycardTradeQuery(o db.Order, callbackform *db.ToMycardTradeQueryForm) bool {
	util.Test(fmt.Sprintf("MycardTradeQuery Order 驗證 3.3 交易單明細 : %+v", o))
	db.SQLDBX.First(&o, OrderSubidQuery(o))
	o.ReceivedCallbackDate = util.GetUTCTime()
	o.MycardTradeNo = callbackform.MyCardTradeNo
	o.Status = "1"
	dbrut := db.SQLDBX.Save(o)
	return dbErrBool(dbrut)
}

//MycardPaymentConfirm Order 確認並進行請款 3.4
func MycardPaymentConfirm(o db.Order, callbackform *db.ToMycardPaymentConfirmForm) bool {
	util.Test(fmt.Sprintf("MycardPaymentConfirm Order 確認並進行請款 3.4 交易單明細 : %+v", o))
	db.SQLDBX.First(&o, OrderSubidQuery(o))
	o.ReceivedCallbackDate = util.GetUTCTime()
	o.PaymentConfirm = "1"
	dbrut := db.SQLDBX.Save(o)
	return dbErrBool(dbrut)
}

//TransactionCallBackTSave Order 更新 回應成功
func TransactionCallBackTSave(o db.Order, callbackform *db.TransactioncallbackForm) bool {
	util.Test(fmt.Sprintf("TransactionCallBackTSave Order 更新 回應成功 交易單明細 : %+v", o))
	db.SQLDBX.First(&o, OrderSubidQuery(o))
	o.ReceivedCallbackDate = util.GetUTCTime()
	//TODO
	//o.MycardTradeNo = callbackform.d
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
	util.Test(fmt.Sprintln("All 交易單明細 : ", Orders))
	return Orders
}
