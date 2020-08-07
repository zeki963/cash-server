package controller

import (
	"cash-server/db"
	"cash-server/pkg/util"
	"cash-server/service"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

//LoggerToDB 紀錄到資料庫
func LoggerToDB() gin.HandlerFunc {
	return func(c *gin.Context) {
		var log db.LogConnect
		// 開始時間
		startTime := time.Now()
		// 處理請求
		c.Next()
		// 結束時間
		endTime := time.Now()
		// 執行時間
		latencyTime := endTime.Sub(startTime)
		log.LatencyTime = util.ShortDur(latencyTime)
		// 請求方式
		log.ReqMethod = c.Request.Method
		// 請求路由
		log.ReqURL = c.Request.RequestURI
		// 狀態碼
		log.StatusCode = c.Writer.Status()
		// 請求IP
		log.ClientIP = c.ClientIP()
		//參數 reqbody  reqheader
		reqbody := "{"
		//手刻JSON 抱歉我就爛不會寫 以後再改
		for k, v := range c.Request.Form {
			//fmt.Printf(`"%v:%v"`, k, v)
			v2 := strings.Join(v, "")
			reqbody += `"`
			reqbody += k
			reqbody += `":"`
			reqbody += v2
			reqbody += `",`
		}
		reqbody += "}"
		log.Reqbody = reqbody
		log.Reqheader = fmt.Sprintf("%v", c.Request.Header)
		//LogConnectAdd 寫入紀錄
		service.LogConnectAdd(log)
	}
}
