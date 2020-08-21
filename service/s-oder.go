package service

import (
	"cash-server/db"
	"cash-server/model"
	"cash-server/pkg/util"
)

//OrderAdd 新增Order(order單)
func OrderAdd(o db.Order) {
	o.Status = "1"
	model.OrderAdd(o)
}

//OrderSave 更新Order(order單,MyResp回應訊息)
func OrderSave(o db.Order, Mresp db.Mycardresp, toServerVal string) {
	if Mresp.ReturnCode == "1" {
		model.OrderAuthTSave(o, Mresp, toServerVal)
	} else {
		model.OrderAuthFSave(o, Mresp, toServerVal)
	}
}

//OrderCallbackSave 回應更新Order(子單號,交易結果,回應訊息) 3-2
func OrderCallbackSave(callbackform *db.OrderMycard) {
	//交易成功為 3; 交易失敗為 0
	if callbackform.PayResult == "3" {
		//TODO
		model.OrderCallbackTSave(db.Order{OrderSubID: callbackform.FacTradeSeq}, callbackform)
	}
}

//OrderQueryOne 查詢存在
func OrderQueryOne(o db.Order) bool {
	if model.OrderSubidQuery(o).OrderSubID == "" {
		util.Test("OrderSubID 不存在")
		return false
	}
	util.Test("OrderSubID 存在")
	return true
}

//OrderFind 查整組資料
func OrderFind(o db.Order) (all db.Order) {
	all = model.OrderQueryExist(o)
	return all
}

//OrderQueryAll 查詢存在
func OrderQueryAll(o db.Order) []db.Order {
	return model.OrderQueryInfoAllJSON()
}
