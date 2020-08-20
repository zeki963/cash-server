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
	util.Logger().WithFields(logrus.Fields{
		"name": "Info",
	}).Info("有人連進測試畫面", "Info")
}

//Homepage Homepage
func Homepage(c *gin.Context) {
	c.HTML(200, "index.html", gin.H{
		"title": "Home page",
	})
}

//NoResponse 不存在，返回404
func NoResponse(c *gin.Context) {
	c.JSON(404, resp(0000, nil))
}
