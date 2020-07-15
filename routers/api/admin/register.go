package admin

import (
	madmin "cash-server/model"
	"cash-server/pkg/encryption"
	"cash-server/pkg/util"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RegisterServer  application/x-www-form-url
// @Summary RegisterServer
// @Tags Admin
// @Produce  json
// @Accept  application/x-www-form-urlencoded
// @Param name formData string true "Name"
// @Param cqikey formData string true "cqikey"
// @success 200 {string} string "{"name":"test","token":"123456","status":"SUCCESS"}"
// @success 400 {string} string "{"status":"FAIL"}"
// @Router /admin/register [post]
func RegisterServer(c *gin.Context) {
	name := c.PostForm("name")
	cqikey := c.PostForm("cqikey")
	if name != "" && cqikey == "cqig7777" {
		token := encryption.Md5(name + "123")

		time := util.GETNowsqltime()
		log.Println(name, token, time)
		err := madmin.InsertServer(name, token, time)
		if err != nil {
			c.JSON(400, gin.H{
				"status": err,
			})

		} else {
			c.JSON(http.StatusOK, gin.H{
				"status": "success",
				"name":   name,
				"token":  token,
				"time":   time,
			})
		}
	} else {
		c.JSON(400, gin.H{
			"status": "FAIL",
		})
	}
}

// ListServer application/json  application/x-www-form-urlencoded
// @Summary ListServer
// @Tags Admin
// @Produce  json
// @Accept  application/x-www-form-urlencoded
// @Param cqikey formData string true "cqikey"
// @success 200 {string} string "{"status":"SUCCESS"}"
// @success 400 {string} string "{"status":"FAIL"}"
// @Router /list [post]
func ListServer(c *gin.Context) {
	cqikey := c.PostForm("cqikey")
	if cqikey == "cqig7777" {
		c.JSON(http.StatusOK, gin.H{
			"status": "SUCCESS",
		})
	} else {
		c.JSON(400, gin.H{
			"status": "FAIL",
		})
	}
}
