package controller

import (
	"cash-server/db"
	"cash-server/pkg/encryption"
	"cash-server/pkg/util"
	"cash-server/service"
	"fmt"
	"regexp"

	"github.com/gin-gonic/gin"
)

//AddRespMsg  回復格式
type AddRespMsg struct {
	Account     string
	Token       string
	TokenSecret string
	Time        string
}

//Keyinrule 輸入規則
var Keyinrule, _ = regexp.Compile("^([a-zA-Z0-9]+$)")

// PlatformAdd  新增平台帳號
// @Summary PlatformRegisterServer 新增平台帳號
// @Tags Admin-Platform
// @Produce  json
// @Accept  application/x-www-form-urlencoded
// @Param PlatformName formData string true "平台名稱"
// @Param PlatformAccount formData string true "平台帳號"
// @Param PlatformPassword formData string true "平台密碼"
// @Param PlatformEmail formData string false "聯絡信箱"
// @Param PlatformGroupID formData int false "群組代號 - 預設為1"
// @success 200 {object}  Resp "OK"
// @success 411 {object}  Resp "Error"
// @Router /admin/platform [post]
func PlatformAdd(c *gin.Context) {
	var p db.Platform
	var resmsg AddRespMsg
	c.Bind(&p)
	if p.PlatformGroupID == 0 {
		p.PlatformGroupID = 1
	}
	if p.PlatformAccount != "" && p.PlatformPassword != "" && p.PlatformName != "" {
		if Keyinrule.MatchString(p.PlatformAccount) && Keyinrule.MatchString(p.PlatformPassword) && Keyinrule.MatchString(p.PlatformName) {
			util.Test(fmt.Sprint("註冊 PlatformRegisterServer ->", p.PlatformAccount))
			//p.PlatformTokenSecret = encryption.Rand16()
			//p.PlatformToken = encryption.PswEncrypt(p.PlatformAccount, p.PlatformTokenSecret)
			p.PlatformTokenSecret, p.PlatformToken = tokenCreate(p.PlatformAccount)
			resmsg.Time = util.GETNowsqltime()
			resmsg.Token = p.PlatformToken
			resmsg.TokenSecret = p.PlatformTokenSecret
			resmsg.Account = p.PlatformAccount
			p.Status = "0"
			add := service.PlatformAdd(p)
			if add == true {
				c.JSON(200, resp(200, resmsg))
			} else {
				c.JSON(411, resp(1002, nil))
			}
		} else {
			c.JSON(411, resp(1005, nil))
		}
	} else {
		c.JSON(411, resp(1001, nil))
	}
}

//tokenCreate 產生key
func tokenCreate(Account string) (tokenSe string, token string) {
	for i := 0; i < 10; i++ {
		tokenSe = encryption.Rand16()
		token = encryption.PswEncrypt(Account, tokenSe)
		if Keyinrule.MatchString(token) {
			util.Test(fmt.Sprint(token, " ==> ToKen 給酷接Q 給過 ! "))
			break
		}
	}

	return tokenSe, token
}

//PlatformGet 要帳號拉幹
// @Summary PlatformGet 查詢帳號
// @Tags Admin-Platform
// @Produce  json
// @Accept  application/x-www-form-urlencoded
// @Param PlatformAccount path string true "平台帳號"
// @success 200 {object}  Resp "OK"
// @success 411 {object}  Resp "Error"
// @Router /admin/platform/{PlatformAccount} [get]
func PlatformGet(c *gin.Context) {
	var p db.Platform
	p.PlatformAccount = c.Param("Acc")
	util.Test(fmt.Sprint("查詢 PlatformGet -> ", p.PlatformAccount))
	if service.PlatformQueryOne(p) {
		c.JSON(200, resp(200, service.PlatformFind(p)))
	} else {
		c.JSON(411, resp(1003, nil))
	}
}

//PlatformGetAll 要ALL帳號拉幹
// @Summary PlatformGetAll 查詢ALL帳號
// @Tags Admin-Platform
// @Produce  json
// @Accept  application/x-www-form-urlencoded
// @success 200 {object} []db.Platform
// @Router /admin/platforms [get]
func PlatformGetAll(c *gin.Context) {
	var p db.Platform
	util.Test(fmt.Sprint("查詢 PlatformAll "))
	c.JSON(200, service.PlatformQueryAll(p))
}

//PlatformStatusEnable 啟用帳號
// @Summary PlatformStatusEnable 啟用帳號
// @Tags Admin-Platform
// @Produce  json
// @Accept  application/x-www-form-urlencoded
// @Param PlatformAccount path string true "平台帳號"
// @success 200 {object}  Resp "OK"
// @success 411 {object}  Resp "Error"
// @Router /admin/platform/status/{PlatformAccount} [put]
func PlatformStatusEnable(c *gin.Context) {
	var p db.Platform
	p.PlatformAccount = c.Param("Acc")
	util.Test(fmt.Sprint("[啟用帳號] PlatformStatusEnable ->", p.PlatformAccount))
	if service.PlatformQueryOne(p) {
		if service.PlatformStatusEnable(p) {
			c.JSON(200, resp(1101, p.PlatformAccount))
		} else {
			c.JSON(200, resp(1102, nil))
		}
	} else {
		c.JSON(411, resp(1003, nil))
	}
}

//PlatformStatusDisable 停用帳號
// @Summary PlatformStatusDisable 停用帳號
// @Tags Admin-Platform
// @Produce  json
// @Accept  application/x-www-form-urlencoded
// @Param PlatformAccount path string true "平台帳號"
// @success 200 {object}  Resp "OK"
// @success 411 {object}  Resp "Error"
// @Router /admin/platform/status/{PlatformAccount} [delete]
func PlatformStatusDisable(c *gin.Context) {
	var p db.Platform
	p.PlatformAccount = c.Param("Acc")
	util.Test(fmt.Sprint("[停用帳號] PlatformStatusEnable ->", p.PlatformAccount))
	if service.PlatformQueryOne(p) {
		if service.PlatformStatusDisable(p) {
			c.JSON(200, resp(1102, p.PlatformAccount))
		} else {
			c.JSON(200, resp(1101, nil))
		}
	} else {
		c.JSON(411, resp(1003, nil))
	}
}
