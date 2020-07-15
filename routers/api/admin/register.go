package admin

import (
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
// @Router /register [post]
func RegisterServer(c *gin.Context) {
	name := c.PostForm("name")
	cqikey := c.PostForm("cqikey")
	token := "123456"
	if name != "" && cqikey == "cqig7777" {
		log.Println(name, token)
		c.JSON(http.StatusOK, gin.H{
			"status": "SUCCESS",
			"name":   name,
			"token":  token,
		})
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
