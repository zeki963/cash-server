package controller

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
	200:  "ok",
	1001: "請求參數短少或錯誤",
	1002: "帳號重複",
	1003: "帳號不存在 & 錯誤",
	1005: "該文章不存在",
	2003: "Token  異常錯誤",
	2002: "Token 已超時",
	2001: "Token 生成失敗",
	2000: "Token錯誤",
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
	return MsgFlags[9999]
}
