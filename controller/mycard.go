package controller

import (
	"cash-server/configs"
	"cash-server/db"
	casinogrpc "cash-server/grpc/casino"
	"cash-server/model"
	"cash-server/pkg/encryption"
	"cash-server/pkg/util"
	"cash-server/service"
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	var checkgameOrderStatus bool = false
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
					casinoUser := casinogrpc.VetifyUserID(o.OrderClientID)
					checkgameOrderStatus, gameOrderSubID = casinogrpc.SendItemBuy(casinoUser, o.OrderItemID)
					o.OrderGameSubID = gameOrderSubID
					o.OrderSubID = service.GroupOrderGet(p.PlatformGroupID, o.StageType)
					o.PaymentTypeID = p.PlatformGroupID
					o.PlatformID = int(p.ID)
				}
			}
			if checkgameOrderStatus != false {
				service.OrderAdd(o)
				//資料庫建單
				util.Test(fmt.Sprintf("Order 資料： %+v", o))
				nmycarderp, toServerVal := toMycardSandAuthGlobal(o.OrderClientID, o.OrderItemID, o.OrderItemPrice, strconv.Itoa(p.PlatformGroupID), o.OrderSubID)
				service.OrderSave(o, nmycarderp, toServerVal)
				if nmycarderp.ReturnCode == "1" {
					//給前端3-2
					util.Test("Redirect mycard page")
					//c.Redirect(301, "http://test.mycard520.com.tw/MyCardPay/?AuthCode="+nmycarderp.AuthCode)
					c.JSON(200, resp(200, fmt.Sprint("http://test.mycard520.com.tw/MyCardPay/?AuthCode="+nmycarderp.AuthCode)))
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
func toMycardSandAuthGlobal(userid string, itemid string, itemprice string, serverid string, subid string) (Mresp db.Mycardresp, toServerVal string) {
	util.Info("<< ToMycardAuthGlobal 向 MyCard 要求交易授權碼 >>")
	var (
		authURL             = "https://testb2b.mycard520.com.tw/MyBillingPay/v1.1/AuthGlobal"
		facServiceID string = configs.GetGlobalConfig().Mycard.FacServiceID               //廠商服務代碼
		facTradeSeq  string = subid                                                       //廠商交易序號 - 廠商自訂，每筆訂單編號不得重覆，為訂單資料 key 值(只能用英數、底線(_)及連字號(-))
		tradeType    string = "2"                                                         //交易模式 - 1:Android SDK (手遊適用) 2:WEB
		serverID     string = serverid                                                    //伺服器代號 - 用戶在廠商端的伺服器編號(僅允許 0-9 a-z A-Z . _ - )
		customerID   string = userid                                                      //會員代號 - 用戶在廠商端的會員唯一識別編號(僅允許 0-9 a-z A-Z . _ -# $ % & * ~ : / ^ ! + @)
		paymentType  string                                                               //付費方式 / 付費方式群組代碼 - 參閱 4.1 付費方式和品項代碼查詢 可不填**
		itemCode     string                                                               //品項代碼  - 參閱 4.1 付費方式和品項代碼查詢 可不填**
		productName  string = itemid                                                      //產品名稱 - 用戶購買的產品名稱(不可以輸入 ' < > 其他皆可)
		amount       string = itemprice                                                   //交易金額 - 可以為整數，若有小數點最多 2 位
		currency     string = "TWD"                                                       //幣別
		sandBoxMode  string = "true"                                                      //是否為測試環境
		facReturnURL string = "https://test-cash.cqiserv.com/mycardsandbox/ordercallback" //廠商回傳網址
	)
	const (
		Key string = "CQIGamesQ1FJR2FtZXM" //我們的KEY
	)
	preHashValue := facServiceID + facTradeSeq + tradeType + serverID + customerID + paymentType + itemCode + productName + amount + currency + sandBoxMode + encryption.Urlencode(facReturnURL) + Key //準備加密字串
	util.Test(fmt.Sprint("preHashValue : ", preHashValue))
	hash := (encryption.Sha256encode(preHashValue))
	util.Test(fmt.Sprint("hash : ", hash))
	toServerVal = "FacServiceId=" + facServiceID + "&FacTradeSeq=" + facTradeSeq + "&TradeType=" + tradeType +
		"&ServerId=" + serverID + "&CustomerId=" + customerID + "&PaymentType=" + paymentType + "&ItemCode=" + itemCode +
		"&ProductName=" + productName + "&Amount=" + amount + "&Currency=" + currency + "&SandBoxMode=" + sandBoxMode +
		"&FacReturnURL=" + facReturnURL + "&Hash=" + hash
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
	return nmycarderp, toServerVal
}

//CreateOrder  建立Order
func CreateOrder(o db.Order) {
	util.Test("CreateMycardOrder")
}

// CallbackMycard  這是給Mycard 摳背專用的 3.2
func CallbackMycard(c *gin.Context) {
	util.Info("<< Mycard 摳背專用的 3.2 >>")
	OrderMycard := &db.OrderMycard{}
	if err := c.Bind(&OrderMycard); err != nil {
		util.Error(err.Error())
	}
	util.Test(fmt.Sprintf("%+v", OrderMycard))
	//Mycard 驗證成功
	if OrderMycard.PayResult == "3" {
		service.OrderCallbackSave(OrderMycard)
		//TODO  Redirect  網址要依遊戲別更換
		switch (service.OrderFind(db.Order{OrderSubID: OrderMycard.FacTradeSeq}).PaymentTypeID) {
		case 2:
			c.Redirect(301, "https://www.cqicasino.com/mycard?succeed=true")
		}
	} else {
		// 失敗
		switch (service.OrderFind(db.Order{OrderSubID: OrderMycard.FacTradeSeq}).PaymentTypeID) {
		case 2:
			//c.JSON(411, resp(3002, fmt.Sprint(encryption.Urldecrypt(OrderMycard.ReturnMsg))))
			msg, _ := encryption.Urldecrypt(OrderMycard.ReturnMsg)
			c.Redirect(301, "https://www.cqicasino.com/mycard?succeed=false&msg="+msg)
		}

	}
}

// Transactioncallback  這是給Mycard 摳背專用的 3.6
func Transactioncallback(c *gin.Context) {
	util.Info("<< Mycard 摳背專用的 3.6 >>")
	transactioncallbackForm := &db.TransactioncallbackForm{}
	form := &db.MycardData{}
	if err := c.Bind(form); err != nil {
		util.Error(err.Error())
	}
	json.Unmarshal([]byte(form.DATA), &transactioncallbackForm)
	service.Transactioncallback(transactioncallbackForm)
	c.JSON(200, resp(200, form))
}

// TransactionCheck  這是給Mycard 摳背專用的 3.7
func TransactionCheck(c *gin.Context) {
	util.Info("<< Mycard 摳背專用的 3.7 >>")
	var retu string
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
	//var backform TransactionCheckBackForm
	form := &TransactionCheckForm{}
	c.Bind(&form)
	fmt.Printf("%+v", form)
	if form.MyCardTradeNo != "" || form.StartDateTime != "" {
		if form.StartDateTime != "" {
			util.Test("多筆查詢")
			retu = model.OrderQueryInfoMoreMycard(form.StartDateTime, form.EndDateTime)
		} else {
			util.Test("單筆查詢")
			var o db.Order
			o.MycardTradeNo = form.MyCardTradeNo
			newo := service.OrderQueryOneMyCardTradeNo(o)
			// backform.Amount = newo.OrderItemPrice
			// backform.TradeSeq = newo.PaymentID
			// backform.PaymentType = newo.PaymentType
			// backform.MyCardTradeNo = newo.MycardTradeNo
			// backform.FacTradeSeq = newo.OrderSubID
			// backform.CustomerID = newo.OrderClientID
			// backform.Currency = "NTD"
			// backform.TradeDateTime = newo.OrderOriginalData
			if newo.PaymentType != "" {
				retu = newo.PaymentType + "," + newo.PaymentID + "," + newo.MycardTradeNo + "," + newo.OrderSubID + "," + newo.OrderClientID + "," + newo.OrderItemPrice + "," + "NTD," + newo.OrderDate.Format("2006-01-02T15:04:05") + "<BR>"
			}
		}
		c.String(200, retu)
	} else {
		c.JSON(200, resp(1001, nil))
	}
}
