package controller

import (
	"cash-server/db"
	"cash-server/pkg/encryption"
	"cash-server/pkg/util"
	"cash-server/service"

	"github.com/gin-gonic/gin"
)

//AddRespMsg  回復格式
type AddRespMsg struct {
	Account     string
	Token       string
	TokenSecret string
	Time        string
}

// PlatformRegisterServerAdd  新增平台帳號
// @Summary PlatformRegisterServer 新增平台帳號
// @Tags Admin
// @Produce  json
// @Accept  application/x-www-form-urlencoded
// @Param PlatformName formData string true "平台名稱"
// @Param PlatformAccount formData string true "平台帳號"
// @Param PlatformPassword formData string true "平台密碼"
// @Param PlatformEmail formData string false "聯絡信箱"
// @Param PlatformGroupID formData string false "群組代號 - 預設為1"
// @success 200 {string} string "{"ecode":200,"msg":"ok","data":{"Account":"test123","Token":"7z7c***Cqw","TokenSecret":"wq0o***XC","Time":"2020-08-06 15:53:36"}}"
// @success 411 {string} string "{"ecode":1001,"msg":"請求參數短少或錯誤","data":null}"
// @Router /admin/platform [post]
func PlatformRegisterServerAdd(c *gin.Context) {
	var p db.PaymentPlatform
	var resmsg AddRespMsg
	c.Bind(&p)
	if p.PlatformGroupID == 0 {
		p.PlatformGroupID = 1
	}
	if p.PlatformAccount != "" && p.PlatformPassword != "" && p.PlatformName != "" {
		util.Info("註冊 PlatformRegisterServer ->", p.PlatformAccount)
		p.PlatformTokenSecret = encryption.Rand16()
		p.PlatformToken = encryption.PswEncrypt(p.PlatformAccount, p.PlatformTokenSecret)
		resmsg.Time = util.GETNowsqltime()
		resmsg.Token = p.PlatformToken
		resmsg.TokenSecret = p.PlatformTokenSecret
		resmsg.Account = p.PlatformAccount
		p.Status = "0"
		add := service.PaymentPlatformAdd(p)
		println(add)
		if add == true {
			c.JSON(200, resp(200, resmsg))
		} else {
			c.JSON(411, resp(1002, nil))
		}
	} else {
		c.JSON(411, resp(1001, nil))
	}
}
