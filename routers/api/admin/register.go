package admin

import (
	"cash-server/model"
	"cash-server/pkg/encryption"
	"cash-server/pkg/util"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// PlatformRegisterServer  application/x-www-form-url
// @Summary PlatformRegisterServer
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
func PlatformRegisterServer(c *gin.Context) {
	account := c.PostForm("account")
	pw := c.PostForm("password")
	name := c.PostForm("name")
	email := c.PostForm("email")
	groupid := c.PostForm("groupid")
	if account != "" && pw != "" && name != "" {
		tokenSecret := encryption.Rand16()
		token := encryption.PswEncrypt(account, tokenSecret)
		time := util.GETNowsqltime()
		util.Info("註冊 PlatformRegisterServer ->", account)
		err := model.PlatformRegister(account, pw, name, groupid, email, token, tokenSecret, time)
		if err != nil {
			if strings.Contains(err.Error(), "payment_platform.platform_account_UNIQUE") {
				c.JSON(411, gin.H{
					"status": "fail",
					"msg":    "Account 重複",
				})
			} else {
				c.JSON(411, gin.H{
					"status": "error",
					"msg":    err,
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
			"msg":    "error",
		})
	}
}
