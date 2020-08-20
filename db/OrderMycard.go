package db

//OrderMycard OrderMycard
type OrderMycard struct {
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
