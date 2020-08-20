package controller

import (
	casinogrpc "cash-server/grpc/casino"

	"github.com/gin-gonic/gin"
)

//CasinoGetItem 查詢casino商品清單
// @Summary GetItem 查詢casino商品清單
// @Tags Casino
// @Produce  json
// @Accept  application/x-www-form-urlencoded
// @success 200 {object}  casinogrpc.CasinoItem "OK"
// @Router /casino/shop/ [get]
func CasinoGetItem(c *gin.Context) {
	c.JSON(200, casinogrpc.GetItem())
}

//CasinoGetUser 查詢casino用戶
// @Summary GetItem 查詢casino用戶
// @Tags Casino
// @Produce  json
// @Accept  application/x-www-form-urlencoded
// @Param UserID path string true "CASINO角色帳號"
// @success 200 {object}  Resp "OK"
// @success 411 {object}  Resp "Error"
// @Router /casino/user/{UserID} [get]
func CasinoGetUser(c *gin.Context) {
	userID := c.Param("Acc")
	if name := casinogrpc.VetifyUserName(userID); name != "Error" {
		c.JSON(200, resp(200, name))
	} else {
		c.JSON(200, resp(2001, nil))
	}

}
