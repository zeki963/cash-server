package controller

import "cash-server/pkg/util"

//Resp 回覆結構
type Resp struct {
	Ecode int         `json:"ecode"`
	Msg   string      `json:"msg"`
	Data  interface{} `json:"data"`
}

func resp(code int, data interface{}) Resp {
	var resp Resp
	resp.Ecode = code
	resp.Msg = getMsg(code)
	resp.Data = data
	return resp
}

//MsgFlags errcode to string
var MsgFlags = map[int]string{
	//OK
	200: "ok",
	//1000 基本錯誤
	1001: "請求參數短少或錯誤",
	1002: "Platform 帳號註冊失敗",
	1003: "Platform 資料不存在",
	1004: "Platform 帳號未啟用相關服務",
	1005: "資料不符合規定格式",
	1101: "Platform 帳號已啟用",
	1102: "Platform 帳號已停用",
	1201: "Order 資料不存在",
	1202: "Order ",
	//2000 GameGroup端錯誤
	2001: "Casino 無此帳號或建單異常",
	//3000 第3方錯誤
	3001: "MyCard Auth 錯誤",
	//其他
	8000: "資料庫錯誤",
	9999: "未預期錯誤",
	0000: "page not exists!,你想幹嘛ヽ(`Д´)ノ  ",
}

//getMsg 訊息
func getMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}
	util.Test("缺少正確錯誤碼")
	return MsgFlags[9999]
}
