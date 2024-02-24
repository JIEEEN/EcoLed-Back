package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"golang.org/x/oauth2"

	"github.com/Eco-Led/EcoLed-Back_test/forms"
	"github.com/Eco-Led/EcoLed-Back_test/services"
	"github.com/gin-gonic/gin"
)

type OauthControllers struct{}

var oauthService = new(services.OauthServices)
var redirectURI string

func (ctr OauthControllers) GoogleLogin(c *gin.Context) {
	// Generate state string (service)
	redirectURI = c.Query("redirect-uri")
	state, err := oauthService.GenerateStateString()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate state"})
		return
	}

	// Save state in session
	session := sessions.Default(c)
	session.Set("oauthState", state)
	err = session.Save()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}

	// Redirect to consent page (service)
	url := oauthService.GetLoginURL(state)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (ctr OauthControllers) GoogleCallback(c *gin.Context) {
	// Get state from session
	receivedState := c.Query("state")
	session := sessions.Default(c)
	savedState := session.Get("oauthState")

	// Check state
	if savedState == nil || savedState != receivedState {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid state"})
		return
	}

	// Handle the exchange code to initiate a transport.
	code := c.Query("code")

	// Get token (service)
	oauthToken, err := oauthService.GetOauthToken(code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get userInfo (service)
	userInfo, err := oauthService.GetOauthUser(oauthToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Check whether the user exists (service)
	userExists := oauthService.FindUserExists(userInfo)

	if !userExists { // if user does not exist
		// Save oauthToken and userInfo in session
		tokenJson, err := json.Marshal(oauthToken)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to serialize oauthToken"})
			return
		}
		session.Set("oauthToken", string(tokenJson))

		userInfoJson, err := json.Marshal(userInfo)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to serialize userInfo"})
			return
		}
		session.Set("userInfo", string(userInfoJson))

		err = session.Save()
		if err != nil {
			log.Printf("Session save error: %v", err) // 로깅 추가
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session1"})
			return
		}

		// Redirect to register page
		c.JSON(http.StatusOK, gin.H{"message": "Additional information required", "redirectURL": "/oauth/google/register"})
		return

	} else { // if user exists

		// Save user in DB (service)
		_, token, err := oauthService.SaveOauthUser(oauthToken, userInfo)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Clear state from session
		session.Delete("oauthState")
		err = session.Save()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session2"})
			return
		}

		// Return the user and token
		url := fmt.Sprintf("%s#token=%s", redirectURI, token)
		c.Redirect(http.StatusTemporaryRedirect, url)
	}
}

func (ctr OauthControllers) GoogleRegister(c *gin.Context) {
	// Get oauthToken and userInfo from session
	session := sessions.Default(c)

	tokenJson, ok := session.Get("oauthToken").(string)
	if !ok || tokenJson == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "oauthToken not found in session"})
		return
	}

	var oauthToken *oauth2.Token
	err := json.Unmarshal([]byte(tokenJson), &oauthToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to deserialize oauthToken"})
		return
	}

	userInfoJson, ok := session.Get("userInfo").(string)
	if !ok || userInfoJson == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userInfo not found in session"})
		return
	}

	var userInfo forms.OauthUser
	err = json.Unmarshal([]byte(userInfoJson), &userInfo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to deserialize userInfo"})
		return
	}

	// Get OauthRegisterForm (form)
	var registerForm forms.OauthRegisterForm
	if err := c.ShouldBindJSON(&registerForm); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate with forms.OauthRegisterForm (form)
	if validationError := registerForm.Validate(); validationError != "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": validationError})
		return
	}

	// Register(service)
	user, token, err := oauthService.Register(registerForm, userInfo, oauthToken)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Response with message
	c.JSON(http.StatusOK, gin.H{"message": "Register Success", "user": user, "token": token})
}
