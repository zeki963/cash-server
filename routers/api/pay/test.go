package pay

import (
	madmin "cash-server/model"
	"cash-server/pkg/encryption"
	"cash-server/pkg/util"
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
// @Tags Test
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
func sadf() {

}

// TestRegisterServer  application/x-www-form-url
// @Summary TestRegisterServer
// @Tags Test
// @Produce  json
// @Accept  application/x-www-form-urlencoded
// @Param name formData string true "Name"
// @Param cqikey formData string true "cqikey"
// @success 200 {string} string "{"name":"test","token":"123456","status":"SUCCESS"}"
// @success 400 {string} string "{"status":"FAIL"}"
// @Router /test/A [post]
func TestRegisterServer(c *gin.Context) {
	name := c.PostForm("name")
	cqikey := c.PostForm("cqikey")
	if name != "" && cqikey == "cqig7777" {
		token := encryption.Md5(name + "123")

		time := util.GETNowsqltime()
		log.Println(name, token, time)
		err := madmin.InsertServer(name, token, time)
		if err != nil {
			c.JSON(400, gin.H{
				"status": err,
			})

		} else {
			c.JSON(http.StatusOK, gin.H{
				"status": "success",
				"name":   name,
				"token":  token,
				"time":   time,
			})
		}
	} else {
		c.JSON(400, gin.H{
			"status": "FAIL",
		})
	}
}
