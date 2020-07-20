package mycard

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// PayIndex application/json  application/x-www-form-urlencoded
// @Summary JSONtestPost
// @Produce  json
// @Param name query string true "Name"
// @Param message query string true "Message"
// @Param num query int true "Num"
// @success 200 {string} string "{"code":"200","message":"","name":"","num":"","status":"SUCCESS"}"
// @Router /json_post [post]
func PayIndex(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"msg":    "success",
	})
}
