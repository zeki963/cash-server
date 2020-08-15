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
		platformStatus := model.PlatformQueryStatusUseToken(token)
		//平台確認
		if platformStatus == "1" {
			//id  確認
			c.JSON(200, resp(200, nil))
		} else {
			c.JSON(400, resp(1001, nil))
		}
	} else {
		c.JSON(400, resp(1001, nil))
	}
}
