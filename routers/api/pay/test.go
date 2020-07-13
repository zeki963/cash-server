package pay

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	name    string
	message string
	num     string
)

//UrlencodedPost 1
func UrlencodedPost(c *gin.Context) {

	name := c.PostForm("name")
	message := c.PostForm("message")
	num := c.DefaultPostForm("num", "1231412")

	log.Println(name, message, num)
	c.JSON(http.StatusOK, gin.H{
		"status":  "SUCCESS",
		"name":    name,
		"message": message,
		"num":     num,
	})
}

//User 123
type User struct {
	Name    string `json:"name" form:"name"`
	Message string `json:"message" form:"message"`
	Num     string `json:"num" form:"num"`
}

// JSONtestPost application/json  application/x-www-form-urlencoded
// @Summary JSONtestPost
// @Produce  json
// @Param name query string true "Name 名稱"
// @Param message query string true "Message"
// @Param num query int true "Num"
// @success 200 {string} string "{"code":"200","message":"","name":"","num":"","status":"SUCCESS"}"
// @Router /json_post [post]
func JSONtestPost(c *gin.Context) {
	var user User

	c.Bind(&user)

	log.Println(user.Name, user.Message, user.Num)

	c.JSON(http.StatusOK, gin.H{
		"code":    "200",
		"status":  "SUCCESS",
		"name":    user.Name,
		"message": user.Message,
		"num":     user.Num,
	})
}
