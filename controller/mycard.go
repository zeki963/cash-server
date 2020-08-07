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
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

//Mycardresp Mycardresp 格式
type Mycardresp struct {
	InGameSaveType int
	ReturnCode     string
	ReturnMsg      string
	AuthCode       string
	TradeSeq       string
}

// MycardOderAdd  application/json  application/x-www-form-urlencoded
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
func MycardOderAdd(c *gin.Context) {
	var p db.PaymentPlatform
	c.Bind(&p)
	groupid := c.PostForm("groupid")
	userid := c.PostForm("userid")
	itemid := c.PostForm("itemid")
	itemprice := c.PostForm("itemprice")
	var guid string
	if userid != "" && itemid != "" && itemprice != "" {
		util.Info("要求Mycard認証 -> token:", p.PlatformToken)
		//平台 Status確認  && GroupAuth確認
		if service.PlatformQueryStatus(p) && service.PlatformGroupAuthQuery(groupid, "1") {
			//依各group確認該遊戲帳戶是否存在
			switch groupid {
			case "2":
				guid = casinogrpc.VetifyUserID(userid)
			}
			if guid != "NoAccount" {
				//test
				c.JSON(http.StatusCreated, gin.H{
					"status":    "success",
					"userid":    userid,
					"itemid":    itemid,
					"itemprice": itemprice,
					"token":     p.PlatformToken,
				})
				//給前端3-2
				//c.JSON(200, resp(200, Nresp))
			} else {
				c.JSON(411, resp(1003, nil))
			}
		} else {
			c.JSON(411, resp(2003, nil))
		}
	} else {
		c.JSON(411, resp(1001, nil))
	}
}

//3.1要求交易授權碼
// authGlobal 向 MyCard 要求交易授權碼 (Server to Server)
func toMycardAuthGlobal(userid string, itemid string, itemprice string, serverid string) {
	util.Info("toMycardAuthGlobal")
	var (
		//測試 http(s)://testb2b.mycard520.com.tw/MyBillingPay/v1.1/AuthGlobal
		//正式 https://b2b.mycard520.com.tw/MyBillingPay/v1.1/AuthGlobal
		authURL = "https://testb2b.mycard520.com.tw/MyBillingPay/v1.1/AuthGlobal"
		//廠商服務代碼
		facServiceID string = configs.GetGlobalConfig().Mycard.FacServiceID
		//廠商交易序號 - 廠商自訂，每筆訂單編號不得重覆，為訂單資料 key 值(只能用英數、底線(_)及連字號(-))
		facTradeSeq string = "1"
		//交易模式 - 1:Android SDK (手遊適用) 2:WEB
		tradeType string = "2"

		//伺服器代號 - 用戶在廠商端的伺服器編號(僅允許 0-9 a-z A-Z . _ - )
		serverID string = serverid
		//會員代號 - 用戶在廠商端的會員唯一識別編號(僅允許 0-9 a-z A-Z . _ -# $ % & * ~ : / ^ ! + @)
		customerID string = userid

		//付費方式 / 付費方式群組代碼 - 參閱 4.1 付費方式和品項代碼查詢 可不填**
		paymentType string
		//品項代碼  - 參閱 4.1 付費方式和品項代碼查詢 可不填**
		itemCode string

		//產品名稱 - 用戶購買的產品名稱(不可以輸入 ' < > 其他皆可)
		productName string = itemid
		//交易金額 - 可以為整數，若有小數點最多 2 位
		amount string = itemprice
		//幣別
		currency string = "TWD"

		//是否為測試環境
		sandBoxMode string = "true"
		//廠商回傳網址
		facReturnURL string = ""
	)
	preHashValue := facServiceID + facTradeSeq + tradeType + serverID + customerID + paymentType + itemCode + productName + amount + currency + sandBoxMode + facReturnURL
	//轉換加密 測試連結 http://test.mycard520.com.tw/FactoryTestTool/MyCardPayCpTest/HASH.aspx
	hash := (encryption.Sha256encode(encryption.Urlencode(preHashValue)))
	toServerVal := "FacServiceId=" + facServiceID + "FacTradeSeq=" + facTradeSeq + "TradeType=" + tradeType +
		"ServerId=" + serverID + "CustomerId=" + customerID + "PaymentType=" + paymentType + "ItemCode=" + itemCode +
		"ProductName=" + productName + "Amount=" + amount + "Currency=" + currency + "SandBoxMode=" + sandBoxMode +
		"FacReturnURL=" + facReturnURL + "Hash=" + hash
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
	//JSON
	util.Trace("Mycardresp:")
	util.Trace(string(body))
	var nmycarderp Mycardresp
	data := []byte(body)
	json.Unmarshal(data, &nmycarderp)
	fmt.Printf("%+v", nmycarderp)
}

//Savedb  建立Oder
func Savedb() {

}

//createMycardOder  建立Oder
func createMycardOder() {
	util.Info("createMycardOder")
	var nmycarderp Mycardresp
	msg := `{"InGameSaveType":"2","ReturnCode":"1","ReturnMsg":"授權成功","AuthCode":"123456","TradeSeq":"KDS1512080000050"}`
	data := []byte(msg)
	json.Unmarshal(data, &nmycarderp)
	fmt.Println(nmycarderp)
	Savedb()
}

// CallbackMycard  這是給Mycard 摳背專用的
func CallbackMycard(c *gin.Context) {
	//回傳結果代碼
	ReturnCode := c.PostForm("ReturnCode")
	//ReturnCode 訊息描述
	ReturnMsg := c.PostForm("ReturnMsg")
	//交易結果代碼
	PayResult := c.PostForm("PayResult")
	//廠商交易序號
	FacTradeSeq := c.PostForm("FacTradeSeq")
	//付費方式
	PaymentType := c.PostForm("PaymentType")
	//金額
	Amount := c.PostForm("Amount")
	//幣別
	Currency := c.PostForm("Currency")
	//MyCardTradeNo
	//1.PaymentType = INGAME 時，傳 MyCard 卡片號碼
	//2.PaymentType = COSTPOINT 時，傳會員扣點交易序號，格式為 CGM 開頭+數字
	//3.其餘 PaymentType 為 Billing 小額付款交易，傳 Billing 交易序號
	MyCardTradeNo := c.PostForm("MyCardTradeNo")
	//通路代碼
	MyCardType := c.PostForm("MyCardType")
	//活動代碼
	PromoCode := c.PostForm("PromoCode")
	//驗證碼
	Hash := c.PostForm("Hash")
	time := util.GETNowsqltime()
	c.JSON(http.StatusCreated, gin.H{
		"status":        "success",
		"ReturnCode":    ReturnCode,
		"ReturnMsg":     ReturnMsg,
		"PayResult":     PayResult,
		"FacTradeSeq":   FacTradeSeq,
		"PaymentType":   PaymentType,
		"Amount":        Amount,
		"Currency":      Currency,
		"MyCardTradeNo": MyCardTradeNo,
		"MyCardType":    MyCardType,
		"PromoCode":     PromoCode,
		"Hash":          Hash,
		"time":          time,
	})
	//TODO 要接DB
}

// // Returnurl33 3-3API
// func Returnurl33(c *gin.Context) {
// 	//交易結果代碼
// 	payResult := c.PostForm("PayResult")
// 	//ReturnCode 訊息CODE
// 	returnCode := c.PostForm("ReturnCode")
// 	//ReturnMsg 訊息描述
// 	returnMsg := c.PostForm("ReturnMsg")
// 	//廠商交易序號
// 	facTradeSeq := c.PostForm("FacTradeSeq")
// 	//付費方式
// 	paymentType := c.PostForm("PaymentType")
// 	//金額
// 	amount := c.PostForm("Amount")
// 	//幣別
// 	currency := c.PostForm("Currency")
// 	//PaymentType = INGAME 時，傳 MyCard 卡片號碼 , .PaymentType = COSTPOINT 時，傳會員扣點交易序號，格式為 CGM 開頭+數字  其餘 PaymentType 為 Billing 小額付款交易，傳 Billing 交易序號 若 BILLING 為卡片儲值時，此欄位傳 MyCard卡片號碼
// 	myCardTradeNo := c.PostForm("MyCardTradeNo")
// 	//通路代碼
// 	myCardType := c.PostForm("MyCardType")
// 	//活動代碼
// 	promoCode := c.PostForm("PromoCode")
// 	//連續扣款序號
// 	serialID := c.PostForm("SerialId")
// 	time := util.GETNowsqltime()
// 	if returnCode != "" {
// 		c.JSON(http.StatusCreated, gin.H{
// 			"ReturnCode":    returnCode,
// 			"ReturnMsg":     returnMsg,
// 			"PayResult":     payResult,
// 			"FacTradeSeq":   facTradeSeq,
// 			"PaymentType":   paymentType,
// 			"Amount":        amount,
// 			"Currency":      currency,
// 			"MyCardTradeNo": myCardTradeNo,
// 			"MyCardType":    myCardType,
// 			"PromoCode":     promoCode,
// 			"SerialId":      serialID,
// 			"time":          time,
// 		})
// 	} else {
// 		c.JSON(400, gin.H{
// 			"status": "Returnurl33 參數不足",
// 		})
// 	}
// }
