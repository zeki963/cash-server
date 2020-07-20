package admin

import (
	"cash-server/model"
	"cash-server/pkg/encryption"
	"cash-server/pkg/util"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// UserRegisterServer  application/x-www-form-url
// @Summary UserRegisterServer
// @Tags Admin
// @Produce  json
// @Accept  application/x-www-form-urlencoded
// @Param name formData string true "Name"
// @Param account formData string true "account"
// @Param password formData string true "password"
// @success 201 {string} string "{"name":"zor","status":"success","time":"2020-07-20 17:25:55","token":"Hqlur4OpUJveSpdy6rm+iw","tokenSecret":"Qvx5CXbGiZ5u1JdF"}"
// @success 400 {string} string "{"status":"FAIL"}"
// @success 411 {string} string "{"status":"FAIL"}"
// @Router /admin/register [post]
func UserRegisterServer(c *gin.Context) {
	account := c.PostForm("account")
	pw := c.PostForm("password")
	name := c.PostForm("name")
	if account != "" && pw != "" && name != "" {
		tokenSecret := encryption.Rand16()
		token := encryption.PswEncrypt(account, tokenSecret)
		time := util.GETNowsqltime()
		log.Println("註冊->", account)
		err := model.UserRegister(account, pw, name, token, tokenSecret, time)
		if err != nil {
			if strings.Contains(err.Error(), "payment_user.user_account_UNIQUE") {
				c.JSON(411, gin.H{
					"status": "fail",
					"msg":    "Account 重複",
				})
			} else {
				c.JSON(411, gin.H{
					"status": "fail",
					"err":    err,
				})
			}
		} else {
			c.JSON(http.StatusCreated, gin.H{
				"status":      "success",
				"account":     account,
				"token":       token,
				"tokenSecret": tokenSecret,
				"time":        time,
			})
		}
	} else {
		c.JSON(400, gin.H{
			"status": "error",
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
