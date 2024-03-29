package service

import (
	"cash-server/configs"
	"cash-server/db"
	casinogrpc "cash-server/grpc/casino"
	"cash-server/model"
	"cash-server/pkg/util"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

//OrderAdd 新增Order(order單)
func OrderAdd(o db.Order) {
	o.Status = "0"
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
	util.Test("OrderCallbackSave ")
	//交易成功為 3; 交易失敗為 0
	if callbackform.PayResult == "3" {
		model.MycardOrderCallbackSave(db.Order{OrderSubID: callbackform.FacTradeSeq}, callbackform)
		ToMycardTradeQuery(model.MycardOrderAuthGet(db.Order{OrderSubID: callbackform.FacTradeSeq}))
	}
}

// ToMycardTradeQuery 驗證 MyCard 交易結果 (Server to Server) 3.3
func ToMycardTradeQuery(AuthCode string) bool {
	authURL := "https://testb2b.mycard520.com.tw/MyBillingPay/v1.1/TradeQuery"
	if configs.GetGlobalConfig().RunEnv == "prod" {
		authURL = "https://b2b.mycard520.com.tw/MyBillingPay/v1.1/TradeQuery"
	}
	toServerVal := "AuthCode=" + AuthCode
	util.Test(fmt.Sprint("ToMycardTradeQuery Auth : ", toServerVal))
	resp, err := http.Post(authURL,
		"application/x-www-form-urlencoded",
		strings.NewReader(toServerVal))
	if err != nil {
		util.Error(err.Error())
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		util.Error(err.Error())
	}
	util.Trace("[Mycard Resp] >> " + string(body))
	var Form db.ToMycardTradeQueryForm
	json.Unmarshal(body, &Form)
	util.Test(fmt.Sprint(Form))
	//交易成功為 3; 交易失敗為 0
	if Form.PayResult == "3" {
		model.MycardTradeQuery(db.Order{OrderSubID: Form.FacTradeSeq}, &Form)
		// -> 3.4
		ToMycardPaymentConfirm(AuthCode)
		//TODO grpc 給order 傳送確認資料 要確認遊戲別
		switch model.OrderSubidQuery(db.Order{OrderSubID: Form.FacTradeSeq}).PaymentTypeID {
		case 2:
			if casinogrpc.SendItemResult(model.OrderSubidQuery(db.Order{OrderSubID: Form.FacTradeSeq})) {
				return false
			}
		}
		return true
	}
	return false
}

// ToMycardPaymentConfirm 確認 MyCard 交易，並進行請款 (Server to Server) 3.4
func ToMycardPaymentConfirm(AuthCode string) bool {
	authURL := "https://testb2b.mycard520.com.tw/MyBillingPay/v1.1/PaymentConfirm"
	if configs.GetGlobalConfig().RunEnv == "prod" {
		authURL = "https://b2b.mycard520.com.tw/MyBillingPay/v1.1/PaymentConfirm"
	}
	toServerVal := "AuthCode=" + AuthCode
	util.Test(fmt.Sprint("toServerVal : ", toServerVal))
	resp, err := http.Post(authURL,
		"application/x-www-form-urlencoded",
		strings.NewReader(toServerVal))
	if err != nil {
		util.Error(err.Error())
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		util.Error(err.Error())
	}
	util.Trace("[Mycard Resp] >> " + string(body))
	var Form db.ToMycardPaymentConfirmForm
	json.Unmarshal(body, &Form)
	util.Test(fmt.Sprintf("%+v", Form))
	if Form.ReturnCode == "1" {
		util.Test("交易成功")
		model.MycardPaymentConfirm(db.Order{OrderSubID: Form.FacTradeSeq}, &Form)
		return true
	}
	util.Test("交易失敗")
	return false
}

//Transactioncallback 回應更新Order(子單號,交易結果,回應訊息) 3-6
func Transactioncallback(callbackform *db.TransactioncallbackForm) {
	for _, keyword := range callbackform.FacTradeSeq {
		util.Test(fmt.Sprintf("SubID =>  %+v", keyword))
		if callbackform.FacServiceID == "CQIG" {
			model.TransactionCallBackTSave(db.Order{OrderSubID: keyword}, callbackform)
			// -> 3.3 -> 3.4
			ToMycardTradeQuery(model.MycardOrderAuthGet(db.Order{OrderSubID: keyword}))
		}
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

//OrderQueryOneMyCardTradeNo 查詢存在 MyCardTradeNo
func OrderQueryOneMyCardTradeNo(o db.Order) db.Order {
	newo := model.OrderSubidQueryMyCardTradeNo(o)
	if newo.OrderClientID == "" {
		util.Test("MyCardTradeNo 不存在")
		return newo
	}
	util.Test("MyCardTradeNo 存在")
	return newo
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
