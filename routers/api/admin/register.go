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
// @Param name formData string true "平台名稱"
// @Param account formData string true "平台帳號"
// @Param password formData string true "平台密碼"
// @Param email formData string false "聯絡信箱"
// @Param groupid formData string false "群組代號 - 預設為1"
// @success 200 {string} string "{"account":"zxcc","status":"success","time":"2020-07-30 17:18:05","token":"18zM3WAF3HKa0Ve4Iss+HA","tokenSecret":"BNFk4iJHx"}"
// @success 400 {string} string "{"status":"FAIL",    "msg": "錯誤訊息"}"
// @success 411 {string} string "{"status":"FAIL",    "msg": "錯誤訊息"}"
// @Router /admin/register [post]
func PlatformRegisterServer(c *gin.Context) {
	account := c.PostForm("account")
	pw := c.PostForm("password")
	name := c.PostForm("name")
	email := c.PostForm("email")
	groupid := c.PostForm("groupid")
	if groupid == "" {
		groupid = "1"
	}
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
