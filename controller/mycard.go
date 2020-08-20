package controller

import (
	"cash-server/configs"
	"cash-server/db"
	casinogrpc "cash-server/grpc/casino"
	"cash-server/pkg/encryption"
	"cash-server/pkg/util"
	"cash-server/service"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// MycardSandOrderAdd  application/json  application/x-www-form-urlencoded
// @Summary 新增 mycard 交易單
// @Tags MyCard-SandBox
// @Produce  json
// @Accept  application/x-www-form-urlencoded
// @Param OrderClientID formData string true "玩家帳號ID"
// @Param OrderItemID formData string true "商品代號"
// @Param OrderItemPrice formData string true "商品價格"
// @Param PlatformToken formData string true "平台token代號"
// @success 200 {object}  Resp "OK"
// @success 411 {object}  Resp "Error"
// @Router /mycardsandbox/order [post]
func MycardSandOrderAdd(c *gin.Context) {
	var p db.Platform
	c.Bind(&p)
	var o db.Order
	c.Bind(&o)
	o.StageType = 0
	o.OrderOriginalData = ""
	o.OrderDate = util.GetUTCTime()
	o.ReceivedCallbackDate = util.GetUTCTime()
	var gameOrderStatus bool = false
	if o.OrderClientID != "" && o.OrderItemID != "" && o.OrderItemPrice != "" {
		util.Test(fmt.Sprint("要求Mycard認証 -> token:", p.PlatformToken))
		util.Test(fmt.Sprint("PlatformQueryStatus: ", service.PlatformQueryStatus(p), " / PlatformGroupAuthQuery: ", service.PlatformGroupAuthQuery(p, "1")))
		//平台 Status確認  && GroupAuth確認
		if service.PlatformQueryStatus(p) && service.PlatformGroupAuthQuery(p, "1") {
			//解析 Token 依各group確認該遊戲帳戶是否存在
			p = service.PlatformFind(p)
			switch p.PlatformGroupID {
			case 2:
				s2iPrice, _ := strconv.ParseInt(o.OrderItemPrice, 10, 32)
				if casinogrpc.VetifyItem(o.OrderItemID, int32(s2iPrice)) {
					//casino資料庫建單
					var gameOrderSubID int32
					casinoUser := casinogrpc.VetifyUserGUID(o.OrderClientID)
					gameOrderStatus, gameOrderSubID = casinogrpc.SendItemBuy(casinoUser, o.OrderItemID)
					o.OrderGameSubID = strconv.FormatInt(int64(gameOrderSubID), 10)
					o.OrderSubID = service.GroupOrderGet(p.PlatformGroupID, o.StageType)
					o.PaymentTypeID = p.PlatformGroupID
					o.PlatformID = int(p.ID)
				}
			}
			if gameOrderStatus != false {
				service.OrderAdd(o)
				//資料庫建單
				util.Test(fmt.Sprintf("Order 資料： %+v", o))
				nmycarderp := toMycardSandAuthGlobal(o.OrderClientID, o.OrderItemID, o.OrderItemPrice, strconv.Itoa(p.PlatformGroupID), o.OrderSubID)
				service.OrderSave(o, nmycarderp)
				if nmycarderp.ReturnCode == "1" {
					//給前端3-2
					c.Redirect(http.StatusMovedPermanently, "http://test.mycard520.com.tw/MyCardPay/?AuthCode="+nmycarderp.AuthCode)
				} else {
					c.JSON(411, resp(3001, nmycarderp))
				}
			} else {
				c.JSON(411, resp(2001, nil))
			}
		} else {
			c.JSON(411, resp(1004, nil))
		}
	} else {
		c.JSON(411, resp(1001, nil))
	}
}

// authGlobal 向 MyCard 要求交易授權碼 (Server to Server) 3.1
func toMycardSandAuthGlobal(userid string, itemid string, itemprice string, serverid string, subid string) (Mresp db.Mycardresp) {
	util.Info("<< ToMycardAuthGlobal 向 MyCard 要求交易授權碼 >>")
	var (
		authURL              = "https://testb2b.mycard520.com.tw/MyBillingPay/v1.1/AuthGlobal"
		facServiceID  string = configs.GetGlobalConfig().Mycard.FacServiceID //廠商服務代碼
		facTradeSeq   string = subid                                         //廠商交易序號 - 廠商自訂，每筆訂單編號不得重覆，為訂單資料 key 值(只能用英數、底線(_)及連字號(-))
		tradeType     string = "2"                                           //交易模式 - 1:Android SDK (手遊適用) 2:WEB
		serverID      string = serverid                                      //伺服器代號 - 用戶在廠商端的伺服器編號(僅允許 0-9 a-z A-Z . _ - )
		customerID    string = userid                                        //會員代號 - 用戶在廠商端的會員唯一識別編號(僅允許 0-9 a-z A-Z . _ -# $ % & * ~ : / ^ ! + @)
		paymentType   string                                                 //付費方式 / 付費方式群組代碼 - 參閱 4.1 付費方式和品項代碼查詢 可不填**
		itemCode      string                                                 //品項代碼  - 參閱 4.1 付費方式和品項代碼查詢 可不填**
		productName   string = itemid                                        //產品名稱 - 用戶購買的產品名稱(不可以輸入 ' < > 其他皆可)
		amount        string = itemprice                                     //交易金額 - 可以為整數，若有小數點最多 2 位
		currency      string = "TWD"                                         //幣別
		sandBoxMode   string = "true"                                        //是否為測試環境
		facReturnURL  string = "https://35.201.250.106/"                     //廠商回傳網址
		facReturnURL2 string = "payCallback"
	)
	const (
		Key string = "CQIGamesQ1FJR2FtZXM" //我們的KEY
	)
	preHashValue := facServiceID + facTradeSeq + tradeType + serverID + customerID + paymentType + itemCode + productName + amount + currency + sandBoxMode + encryption.Urlencode(facReturnURL) + facReturnURL2 + Key //準備加密字串
	util.Test(fmt.Sprint("preHashValue : ", preHashValue))
	hash := (encryption.Sha256encode(preHashValue))
	util.Test(fmt.Sprint("hash : ", hash))
	toServerVal := "FacServiceId=" + facServiceID + "&FacTradeSeq=" + facTradeSeq + "&TradeType=" + tradeType +
		"&ServerId=" + serverID + "&CustomerId=" + customerID + "&PaymentType=" + paymentType + "&ItemCode=" + itemCode +
		"&ProductName=" + productName + "&Amount=" + amount + "&Currency=" + currency + "&SandBoxMode=" + sandBoxMode +
		"&FacReturnURL=" + facReturnURL + facReturnURL2 + "&Hash=" + hash
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
	//Resp
	util.Trace("[Mycard Resp] >> " + string(body))
	//JSON
	var nmycarderp db.Mycardresp
	data := []byte(body)
	json.Unmarshal(data, &nmycarderp)
	return nmycarderp
}

//CreateOrder  建立Order
func CreateOrder(o db.Order) {
	util.Test("CreateMycardOrder")
}

// CallbackMycard  這是給Mycard 摳背專用的 3.2
func CallbackMycard(c *gin.Context) {
	util.Info("<< Mycard 摳背專用的 3.2 >>")
	form := &db.OrderMycard{}
	if err := c.BindJSON(form); err != nil {
		util.Error(err.Error())
	}
	service.OrderCallbackSave(form.FacTradeSeq, form.PayResult, db.Struct2JSON(form))
}

// toMycardTradeQuery 驗證 MyCard 交易結果 (Server to Server) 3.3
func toMycardTradeQuery(AuthCode string) {
	type toMycardTradeQueryForm struct {
		ReturnCode    string //查詢結果代碼
		ReturnMsg     string //ReturnCode 訊息描述
		PayResult     string //交易結果代碼
		FacTradeSeq   string //※廠商交易序號
		PaymentType   string //付費方式
		Amount        string //※金額
		Currency      string //※幣別
		MyCardTradeNo string //※交易成功時間
		MyCardType    string //※通路代碼
		PromoCode     string //※活動代碼
		SerialID      string //※連續扣款序號
	}
	authURL := "https://testb2b.mycard520.com.tw/MyBillingPay/v1.1/TradeQuery"
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
	//Resp
	util.Trace("[Mycard Resp] >> " + string(body))
	//假資料
	body = []byte(`{"ReturnCode":"1","ReturnMsg":"查詢成功","PayResult":"3","FacTradeSeq":"FacTradeSeq0001","PaymentType":"INGAME","Amount":"150","Currency":"TWD", "MyCardTradeNo":"MAAAAA0000000001","MyCardType":"1","PromoCode":"A0000","SerialId":"1"}`)
	var Form toMycardTradeQueryForm
	json.Unmarshal(body, &Form)
	util.Test(fmt.Sprint(Form))
	if Form.PayResult == "3" {
		util.Test("交易成功")
	} else {
		util.Test("交易失敗")
	}
}

// toMycardPaymentConfirm 確認 MyCard 交易，並進行請款 (Server to Server) 3.4
func toMycardPaymentConfirm(AuthCode string) {
	type toMycardPaymentConfirmForm struct {
		ReturnCode  string `json:"ReturnCode"`  //請款結果代碼 Payment Result
		ReturnMsg   string `json:"ReturnMsg"`   //※ReturnCode 訊息描述
		FacTradeSeq string `json:"FacTradeSeq"` //※廠商交易序號
		TradeSeq    string `json:"TradeSeq"`    //※MyCard 交易序號
		SerialID    string `json:"SerialId"`    //※連續扣款序號
	}
	authURL := "https://testb2b.mycard520.com.tw/MyBillingPay/v1.1/TradeQuery"
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
	//Resp
	util.Trace("[Mycard Resp] >> " + string(body))
	//假資料
	body = []byte(`{"ReturnCode":"1","ReturnMsg":"請款成功","FacTradeSeq":"FacTradeSeq0001","TradeSeq":"KDS1512080000050","SerialId":""}`)
	var Form toMycardPaymentConfirmForm
	json.Unmarshal(body, &Form)
	util.Test(fmt.Sprint(Form))
	if Form.ReturnCode == "1" {
		util.Test("交易成功")
	} else {
		util.Test("交易失敗")
	}
}

// Transactioncallback  這是給Mycard 摳背專用的 3.6
//{"ReturnCode":"1","ReturnMsg":"QueryOK","FacServiceId":"MyCardSDK","TotalNum":2,"FacTradeSeq":["FacTradeSeq0001","FacTradeSeq0002"]}
func Transactioncallback(c *gin.Context) {
	util.Info("<< Mycard 摳背專用的 3.6 >>")
	type TransactioncallbackForm struct {
		ReturnCode   string   `form:"ReturnCode" binding:"required"`   //回傳結果代碼
		ReturnMsg    string   `form:"ReturnMsg" binding:"required"`    //ReturnCode 訊息描述
		FacTradeSeq  []string `form:"FacTradeSeq" binding:"required"`  //廠商交易序號
		FacServiceID string   `form:"FacServiceID" binding:"required"` //廠商服務代碼
		TotalNum     string   `form:"TotalNum" binding:"required"`     //交易筆數
	}
	form := &TransactioncallbackForm{}
	c.BindJSON(form)
	util.Test(fmt.Sprint(form))
	c.JSON(200, resp(200, form))
}

// TransactionCheck  這是給Mycard 摳背專用的 3.7
func TransactionCheck(c *gin.Context) {
	util.Info("<< Mycard 摳背專用的 3.7 >>")
	type TransactionCheckForm struct {
		StartDateTime string `form:"StartDateTime" binding:"required"` //※開始日期(UTC+8) yyyy-mm-ddThr:mi:se(24 )
		EndDateTime   string `form:"EndDateTime" binding:"required"`   // ※結束日期(UTC+8) 2014-12-01T00:00:00
		MyCardTradeNo string `form:"MyCardTradeNo" binding:"required"` //MyCardTradeNo
	}
	type TransactionCheckBackForm struct {
		PaymentType   string //※付費方式
		TradeSeq      string //交易序號
		MyCardTradeNo string //MyCardTradeNo
		FacTradeSeq   string //※廠商交易序號
		CustomerID    string //※會員代號
		Amount        string //※金額
		Currency      string //※幣別
		TradeDateTime string //※交易成功時間
	}
	var backform TransactionCheckBackForm
	form := &TransactionCheckForm{}
	c.BindJSON(form)
	fmt.Println(form)
	if form.MyCardTradeNo != "" {
		util.Test("多筆查詢")
	} else {
		util.Test("單筆查詢")
		var o db.Order
		o.OrderSubID = form.MyCardTradeNo
		service.OrderQueryOne(o)
		backform.Amount = o.OrderItemPrice
		backform.PaymentType = "PaymentType"
		backform.TradeSeq = "KDS1512080000050"
		backform.MyCardTradeNo = "MyCardTradeNo"
		backform.FacTradeSeq = o.OrderSubID
		backform.CustomerID = o.OrderClientID
		backform.Currency = "NTD"
		backform.TradeDateTime = o.OrderOriginalData
	}

	jsonbackform, err := json.Marshal(backform)
	if err != nil {
		log.Fatal(err)
	}
	c.JSON(200, resp(200, jsonbackform))
}
