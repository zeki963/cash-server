package controller

import (
	"cash-server/pkg/util"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

//Demopage Demopage
func Demopage(c *gin.Context) {
	c.HTML(200, "demo.html", gin.H{
		"title": "DEMO page",
	})
}

//Homepage Homepage
func Homepage(c *gin.Context) {
	c.HTML(200, "index.html", gin.H{
		"title": "Home page",
	})
	util.Logger().WithFields(logrus.Fields{
		"name": "Info",
	}).Info("有人連進根目錄了", "Info")
}

//NoResponse 不存在，返回404
func NoResponse(c *gin.Context) {
	c.JSON(404, res(0000, nil))
}
