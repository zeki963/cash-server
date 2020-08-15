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

// MycardSandOderAdd  application/json  application/x-www-form-urlencoded
// @Summary AuthMycard
// @Tags MyCard
// @Produce  json
// @Accept  application/x-www-form-urlencoded
// @Param userid formData string true "玩家帳號ID"
// @Param groupid formData string true "遊戲GroupID"
// @Param itemid formData string true "商品代號"
// @Param itemprice formData string true "商品價格"
// @Param token formData string true "平台token代號"
// @success 200 {string} string "{"還沒寫好"}"
// @success 400 {string} string "{"status":"FAIL",    "msg": "錯誤訊息"}"
// @Router /mycard/CreateMycardOder [post]
func MycardSandOderAdd(c *gin.Context) {
	var p db.PaymentPlatform
	c.Bind(&p)
	var o db.Order
	c.Bind(&o)
	o.StageType = 0
	var guid string
	if o.OrderClientID != "" && o.OrderItemID != "" && o.OrderItemPrice != "" {
		util.Info("要求Mycard認証 -> token:", p.PlatformToken)
		util.Test(fmt.Sprintln(service.PlatformQueryStatus(p), service.PlatformGroupAuthQuery(p, "1")))
		//平台 Status確認  && GroupAuth確認
		if service.PlatformQueryStatus(p) && service.PlatformGroupAuthQuery(p, "1") {
			//解析 Token 依各group確認該遊戲帳戶是否存在
			p = service.PaymentPlatformFind(p)
			switch p.PlatformGroupID {
			case 2:
				guid = casinogrpc.VetifyUserID(o.OrderClientID)
			}
			o.OrderDate = util.GetUTCTime()
			o.ReceivedCallbackDate = util.GetUTCTime()
			o.OrderSubID = service.GroupOrderGet(p.PlatformGroupID, o.StageType)
			o.PaymentTypeID = p.PlatformGroupID
			o.PlatformID = int(p.ID)
			if guid != "NoAccount" {
				//資料庫建單
				fmt.Printf("%+v", o)
				service.OrderAdd(o)
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
	util.Info("toMycardAuthGlobal")
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
	util.Test(preHashValue)
	hash := (encryption.Sha256encode(preHashValue))
	util.Test(hash)
	toServerVal := "FacServiceId=" + facServiceID + "&FacTradeSeq=" + facTradeSeq + "&TradeType=" + tradeType +
		"&ServerId=" + serverID + "&CustomerId=" + customerID + "&PaymentType=" + paymentType + "&ItemCode=" + itemCode +
		"&ProductName=" + productName + "&Amount=" + amount + "&Currency=" + currency + "&SandBoxMode=" + sandBoxMode +
		"&FacReturnURL=" + facReturnURL + facReturnURL2 + "&Hash=" + hash
	util.Test(toServerVal)
	resp, err := http.Post(authURL,
		"application/x-www-form-urlencoded",
		strings.NewReader(toServerVal))
	if err != nil {
		fmt.Println(err)
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

//CreateOder  建立Oder
func CreateOder(o db.Order) {
	util.Info("CreateMycardOder")
}

// CallbackMycard  這是給Mycard 摳背專用的 3.2
func CallbackMycard(c *gin.Context) {
	type CallbackMycardForm struct {
		ReturnCode  string `form:"ReturnCode" binding:"required"`  //回傳結果代碼
		ReturnMsg   string `form:"ReturnMsg" binding:"required"`   //ReturnCode 訊息描述
		PayResult   string `form:"PayResult" binding:"required"`   //交易結果代碼
		FacTradeSeq string `form:"FacTradeSeq" binding:"required"` //廠商交易序號
		PaymentType string `form:"PaymentType" binding:"required"` //付費方式
		Amount      string `form:"Amount" binding:"required"`      //金額
		Currency    string `form:"Currency" binding:"required"`    //幣別
		//MyCardTradeNo
		//1.PaymentType = INGAME 時，傳 MyCard 卡片號碼
		//2.PaymentType = COSTPOINT 時，傳會員扣點交易序號，格式為 CGM 開頭+數字
		//3.其餘 PaymentType 為 Billing 小額付款交易，傳 Billing 交易序號
		MyCardTradeNo string `form:"MyCardTradeNo" binding:"required"`
		MyCardType    string `form:"MyCardType" binding:"required"` //通路代碼
		PromoCode     string `form:"PromoCode" binding:"required"`  //活動代碼
		Hash          string `form:"Hash" binding:"required"`       //驗證碼
	}
	c.JSON(200, resp(200, nil))
	// 要接DB
}

// Transactioncallback  這是給Mycard 摳背專用的 3.6
func Transactioncallback(c *gin.Context) {
	type TransactioncallbackForm struct {
		ReturnCode   string   `form:"ReturnCode" binding:"required"`   //回傳結果代碼
		ReturnMsg    string   `form:"ReturnMsg" binding:"required"`    //ReturnCode 訊息描述
		FacTradeSeq  []string `form:"FacTradeSeq" binding:"required"`  //廠商交易序號
		FacServiceID string   `form:"FacServiceID" binding:"required"` //廠商服務代碼
		TotalNum     string   `form:"TotalNum" binding:"required"`     //交易筆數
	}
	form := &TransactioncallbackForm{}
	c.BindJSON(form)
	time := util.GETNowsqltime()
	fmt.Println(form)
	c.JSON(200, gin.H{
		"status":       "success",
		"ReturnCode":   form.ReturnCode,
		"ReturnMsg":    form.ReturnMsg,
		"FacServiceId": form.FacServiceID,
		"FacTradeSeq":  form.FacTradeSeq,
		"TotalNum":     form.TotalNum,
		"time":         time,
	})
}

// TransactionCheck  這是給Mycard 摳背專用的 3.7
func TransactionCheck(c *gin.Context) {
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
	form := &TransactionCheckForm{}
	c.BindJSON(form)
	fmt.Println(form)
	var backform TransactionCheckBackForm
	jsonbackform, err := json.Marshal(backform)
	if err != nil {
		log.Fatal(err)
	}
	c.JSON(200, gin.H{
		"message": string(jsonbackform),
	})
}
