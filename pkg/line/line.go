package line

import (
	"cash-server/pkg/util"

	"github.com/line/line-bot-sdk-go/linebot"
)

//Linepost 傳送訊息
func Linepost(text string) {
	bot, err := linebot.New("e9be603e5e7af5232f763b6e849719f3", "tOL98RcL6Boy2fL05G8gxYEaOAV3S9QI+JdgEy8oInbOYz2+hHHm+g5AhES/0xiMQj0ozCm3HHmXiPG1xmSD/L5BIOYtshm3fQjm6Aj7ONe/evE0PT9vjLmJdSHNA8SkDS9WDvMTxVJWAw9nszGg8AdB04t89/1O/w1cDnyilFU=")
	if err != nil {
		util.Error(err.Error())
		util.Error("Linebot error 01")
	}
	if _, err := bot.PushMessage("Uf89452be180a8a97a67276775807a1bf", linebot.NewTextMessage(text)).Do(); err != nil {
		util.Error(err.Error())
		println("Linebot error 02")
	}
	util.Info(text)
}

//SayHello test
func SayHello() {
	Linepost("Hello,I am Line Bot  :) ")
}
