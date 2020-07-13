package asana

import (
	"cash-server/pkg/util"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"gopkg.in/danilopolani/gocialite.v1"
)

const (
	userkey = "user"
)

var (
	clientID     = "1171871475113226"
	clientSecret = "4f490e81f1a51f3f94975705d4890dd0"
	redirectURL  = "https://127.0.0.1:8443/auth/asana/callback"
)

//UserType ASANA retun info
type UserType struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Data        struct {
		ID    int64  `json:"id"`
		Gid   string `json:"gid"`
		Name  string `json:"name"`
		Email string `json:"email"`
	} `json:"data"`
	RefreshToken string `json:"refresh_token"`
}

// RedirectHandler Redirect to correct oAuth URL
func RedirectHandler(c *gin.Context) {
	authURL, err := gocialite.NewDispatcher().New().
		Driver("asana").
		Scopes([]string{}).
		Redirect(
			clientID,     // Client ID
			clientSecret, // Client Secret
			redirectURL,  // Redirect URL
		)

	// Check for errors (usually driver not valid)
	if err != nil {
		c.Writer.Write([]byte("Error: " + err.Error()))
		return
	}

	// Redirect with authURL
	c.Redirect(http.StatusFound, authURL) // Redirect with 302 HTTP code
}

//CallbackHandler Handle callback of provider
func CallbackHandler(c *gin.Context) {

	// Retrieve query params for state and code
	state := c.Query("state")
	code := c.Query("code")
	//2次進攻redirectURL
	resp, err := http.Post("https://app.asana.com/-/oauth_token",
		"application/x-www-form-urlencoded",
		strings.NewReader("grant_type=authorization_code&client_id="+clientID+"&client_secret="+clientSecret+"&redirect_uri="+redirectURL+"&state="+state+"&code="+code))
	if err != nil {
		util.Error(err.Error())
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		util.Error("resp error")
	}

	Serverslice1 := UserType{}
	e := json.Unmarshal([]byte(body), &Serverslice1)
	if e != nil {
		util.Error(e.Error())
	}

	// Save the username in the session
	//session.Set(userkey, Serverslice1.Data.Name)

	//fmt.Println(body)

	//rsp回來的資料
	util.Info(string(" > User "+Serverslice1.Data.Name) + "  login ! ")
	c.Writer.Write([]byte("Hi, " + string(Serverslice1.Data.Name)))
}

// Login is a handler that parses a form and checks for specific data
func Login(c *gin.Context) {
	session := sessions.Default(c)
	username := c.PostForm("username")
	password := c.PostForm("password")

	// Validate form input
	if strings.Trim(username, " ") == "" || strings.Trim(password, " ") == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parameters can't be empty"})
		return
	}

	// Check for username and password match, usually from a database
	if username != "id" || password != "pw" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
		util.Error("Authentication failed")
		return
	}
	// Save the username in the session
	session.Set(userkey, username)
	// In real world usage you'd set this to the users ID
	if err := session.Save(); err != nil {
		util.Error("Failed to save session")
		util.Error(err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully authenticated user"})
	util.Info("Successfully to session")
}

//Logout logoutttttttt
func Logout(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(userkey)
	if user == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session token"})
		return
	}
	session.Delete(userkey)
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}

//Me lookme
func Me(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(userkey)
	c.JSON(http.StatusOK, gin.H{"user": user})
}

//AuthRequired is a simple middleware to check the session
func AuthRequired(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(userkey)
	if user == nil {
		// Abort the request with the appropriate error code
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	// Continue down the chain to handler etc
	c.Next()
}
