package router

import (
	"fmt"
	"io"
	"os"

	"github.com/gin-gonic/gin"
)

//ginlogmode  server logs
func ginlogmode() {
	// Disable log's color
	gin.DisableConsoleColor()
	// Force log's color
	//gin.ForceConsoleColor()
	f, err := os.Create("gin.log")
	if err != nil {
		fmt.Println("Open Log File Failed", err)
	}
	gin.DefaultWriter = io.MultiWriter(f)
}

//Cors Cors
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

//MycardCallAuth Mycard call 白名單
func MycardCallAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		whiteList := []string{
			"127.0.0.1",
			"218.32.37.148",
			"220.130.127.125",
		}
		flag := false
		for _, host := range whiteList {
			if c.ClientIP() == host {
				flag = true
				break
			}
		}
		if !flag {
			c.JSON(511, "your ip is not trusted")
			c.Abort()
		}
	}
}
