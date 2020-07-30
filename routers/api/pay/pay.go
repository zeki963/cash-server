package pay

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
			c.JSON(400, gin.H{
				"status": "權限不足",
			})
		}
	} else {
		c.JSON(400, gin.H{
			"status": "token 錯誤或參數不足",
		})
	}
}
