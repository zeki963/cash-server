package util

import (
	"cash-server/db"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

//Logger Logger
func Logger() *logrus.Logger {
	now := time.Now()
	logFilePath := ""
	if dir, err := os.Getwd(); err == nil {
		logFilePath = dir + "/logs/"
	}
	if err := os.MkdirAll(logFilePath, 0777); err != nil {
		fmt.Println(err.Error())
	}
	logFileName := now.Format("2006-01-02") + ".log"
	//日誌文件
	fileName := path.Join(logFilePath, logFileName)
	if _, err := os.Stat(fileName); err != nil {
		if _, err := os.Create(fileName); err != nil {
			fmt.Println(err.Error())
		}
	}
	//寫入文件
	src, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Println("err", err)
	}
	//實例化
	logger := logrus.New()
	//設置輸出
	logger.Out = src
	//設置日誌級別
	logger.SetLevel(logrus.DebugLevel)
	//設置日誌格式
	logger.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	return logger
}

//LoggerToFile 紀錄到檔案
func LoggerToFile() gin.HandlerFunc {
	logger := Logger()
	return func(c *gin.Context) {
		// 開始時間
		startTime := time.Now()
		// 處理請求
		c.Next()
		// 結束時間
		endTime := time.Now()
		// 執行時間
		latencyTime := endTime.Sub(startTime)
		// 請求方式
		reqMethod := c.Request.Method
		// 請求路由
		reqURL := c.Request.RequestURI
		// 狀態碼
		statusCode := c.Writer.Status()
		// 請求IP
		clientIP := c.ClientIP()
		//日誌格式
		logger.Infof("| %3d | %13v | %15s | %s | %s |",
			statusCode,
			latencyTime,
			clientIP,
			reqMethod,
			reqURL,
		)
	}
}

//LoggerToDB 紀錄到資料庫
func LoggerToDB() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 開始時間
		startTime := time.Now()
		// 處理請求
		c.Next()
		// 結束時間
		endTime := time.Now()
		// 執行時間
		latencyTime := endTime.Sub(startTime)
		// 請求方式
		reqMethod := c.Request.Method
		// 請求路由
		reqURL := c.Request.RequestURI
		// 狀態碼
		statusCode := c.Writer.Status()
		// 請求IP
		clientIP := c.ClientIP()
		//時間
		time := GETNowsqltime()

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
		reqheader := fmt.Sprintf("%v", c.Request.Header)
		//寫入紀錄
		insertLogConnect(statusCode, latencyTime, clientIP, reqMethod, reqURL, time, reqbody, reqheader)
	}
}

//InsertLogConnect 寫入紀錄
func insertLogConnect(statusCode int, latencyTime time.Duration, clientIP, reqMethod string, reqURL string, time string, reqbody, reqheader string) error {
	if _, err := db.SqlDB.Exec(
		"INSERT INTO log_connect (statusCode, latencyTime, clientIP, reqMethod, reqURL, create_date,update_date,reqbody,reqheader) VALUES (?,?,?,?,?,?,?,?,?)",
		statusCode, latencyTime, clientIP, reqMethod, reqURL, time, time, reqbody, reqheader); err != nil {
		Error("DB INSERT fail")
		Error(err.Error())
		return err
	}
	return nil
}
