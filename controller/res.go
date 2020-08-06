package controller

//Res 回覆結構
type Res struct {
	Ret  int         `json:"ret"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func res(code int, data interface{}) Res {
	var res Res
	res.Ret = code
	res.Msg = GetMsg(code)
	res.Data = data
	return res
}

//MsgFlags errcode to string
var MsgFlags = map[int]string{
	200:  "ok",
	1001: "請求參數短少或錯誤",
	1002: "帳號重複",
	1003: "該標籤不存在",
	1005: "該文章不存在",
	2003: "Token鑒權失敗",
	2002: "Token已超時",
	2001: "Token生成失敗",
	2000: "Token錯誤",
	8000: "資料庫錯誤",
	9999: "未預期錯誤",
	0000: "page not exists!,你想幹嘛ヽ(`Д´)ノ  ",
}

//GetMsg 訊息
func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}
	return MsgFlags[9999]
}
