package controller

import (
	"cash-server/db"
	"cash-server/model"
	"cash-server/pkg/encryption"
	"cash-server/pkg/util"
	"strconv"

	"github.com/gin-gonic/gin"
)

//Addrep  回復格式
type Addrep struct {
	Account     string
	Token       string
	TokenSecret string
	Time        string
}

// PlatformRegisterServerAdd  新增平台帳號
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
func PlatformRegisterServerAdd(c *gin.Context) {
	var rep Addrep
	rep.Account = c.PostForm("account")
	pw := c.PostForm("password")
	name := c.PostForm("name")
	email := c.PostForm("email")
	groupid := c.PostForm("groupid")
	if groupid == "" {
		groupid = "1"
	}
	if rep.Account != "" && pw != "" && name != "" {
		if model.PaymentPlatformCheck(rep.Account) {
			c.JSON(411, res(1002, nil))
		} else {
			//PaymentPlatformCheck(rep.Account)
			rep.TokenSecret = encryption.Rand16()
			rep.Token = encryption.PswEncrypt(rep.Account, rep.TokenSecret)
			rep.Time = util.GETNowsqltime()
			util.Info("註冊 PlatformRegisterServer ->", rep.Account)
			//err := model.PlatformRegister(rep.Account, pw, name, groupid, email, rep.Token, rep.TokenSecret, rep.Time)
			var p db.PaymentPlatform
			p.PlatformAccount = rep.Account
			p.PlatformGroupID, _ = strconv.Atoi(groupid)
			p.PlatformName = name
			p.PlatformPassword = pw
			p.PlatformToken = rep.Token
			p.PlatformTokenSecret = rep.TokenSecret
			p.PlatformEmail = email
			p.Status = "0"
			err := model.PaymentPlatformAdd(p)
			if err != nil {
				c.JSON(411, res(8000, rep))
			} else {
				c.JSON(200, res(200, rep))
			}
		}
	} else {
		c.JSON(411, res(1001, nil))
	}
}
