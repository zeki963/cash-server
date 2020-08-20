package controller

import (
	"cash-server/db"
	"cash-server/pkg/util"
	"cash-server/service"
	"fmt"

	"github.com/gin-gonic/gin"
)

//OrderGet 要交易單拉幹
// @Summary OrderGet 查詢交易單
// @Tags Admin-Order
// @Produce  json
// @Accept  application/x-www-form-urlencoded
// @Param OrderSubID path string true "交易單子單號"
// @success 200 {object}  Resp "OK"
// @success 411 {object}  Resp "Error"
// @Router /admin/order/{OrderSubID} [get]
func OrderGet(c *gin.Context) {
	var o db.Order
	o.OrderSubID = c.Param("Subid")
	util.Test(fmt.Sprint("查詢 OrderGet -> ", o.OrderSubID))
	if service.OrderQueryOne(o) {
		c.JSON(200, resp(200, service.OrderFind(o)))
	} else {
		c.JSON(411, resp(1201, nil))
	}
}

//OrderGetAll 要ALL帳號拉幹
// @Summary OrderGetAll 查詢ALL帳號
// @Tags Admin-Order
// @Produce  json
// @Accept  application/x-www-form-urlencoded
// @success 200 {object} []db.Order
// @Router /admin/orders [get]
func OrderGetAll(c *gin.Context) {
	var o db.Order
	util.Test(fmt.Sprint("查詢 OrderAll "))
	c.JSON(200, service.OrderQueryAll(o))
}
