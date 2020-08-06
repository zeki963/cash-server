package controller

import (
	"cash-server/model"

	"github.com/gin-gonic/gin"
)

// Checkrouter application/json  application/x-www-form-urlencoded
func Checkrouter(c *gin.Context) {
	id := c.PostForm("id")
	itemid := c.PostForm("itemid")
	itemprice := c.PostForm("itemprice")
	token := c.PostForm("token")
	//time := util.GETNowsqltime()
	if id != "" && itemid != "" && itemprice != "" {
		platform := model.PlatformQueryInfodata(token)
		//平台確認
		if platform.Status == "1" {
			//id  確認

		} else {
			c.JSON(400, res(1001, nil))
		}
	} else {
		c.JSON(400, res(1001, nil))
	}
}
