package controllers

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"golang.org/x/oauth2"

	"github.com/Eco-Led/EcoLed-Back_test/forms"
	"github.com/Eco-Led/EcoLed-Back_test/services"
	"github.com/gin-gonic/gin"
)

type OauthControllers struct{}

var oauthService = new(services.OauthServices)

func (ctr OauthControllers) GoogleLogin(c *gin.Context) {
	// Generate state string (service)
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

	if !userExists{ // if user does not exist
		// Save oauthToken and userInfo in session
		session.Set("oauthToken", oauthToken)
		session.Set("userInfo", userInfo)
		err = session.Save()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
			return
		}
		
		// Redirect to register page
		c.Redirect(http.StatusFound, "api/v1/oauth/google/register")
        return

	} else{ // if user exists

		// Save user in DB (service)
		user, token, err := oauthService.SaveOauthUser(oauthToken, userInfo)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Clear state from session
		session.Delete("oauthState")
		err = session.Save()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
			return
		}

		// Return the user and token
		c.JSON(http.StatusOK, gin.H{"user": user, "oauthToken": oauthToken, "token": token})
	}

}

func (ctr OauthControllers) GoogleRegister(c *gin.Context) {
	// Get oauthToken and userInfo from session
	session := sessions.Default(c)
	oauthToken := session.Get("oauthToken").(*oauth2.Token)
	userInfo := session.Get("userInfo").(forms.OauthUser)

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