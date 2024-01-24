package controllers

import (
	"net/http"

	"github.com/Eco-Led/EcoLed-Back_test/forms"
	"github.com/Eco-Led/EcoLed-Back_test/services"

	"github.com/gin-gonic/gin"
)

type UserControllers struct{}

var userService = new(services.UserServices)

func (ctr UserControllers) Login(c *gin.Context) {
	// Bind JSON with forms.LoginForm (form)
	var loginForm forms.LoginForm
	if err := c.ShouldBindJSON(&loginForm); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate with forms.UserForm (form)
	userForm := forms.UserForm{}
	if validationError := userForm.Login(loginForm); validationError != "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": validationError,
		})
		return
	}

	// Login(service)
	user, token, err := userService.Login(loginForm)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Response user, token with message
	c.JSON(http.StatusOK, gin.H{
		"message": "Login Success",
		"user":    user,
		"token":   token,
	})

}


func (ctr UserControllers) Register(c *gin.Context) {
	// Bind JSON with forms.RegisterForm (form)
	var registerForm forms.RegisterForm
	if err := c.ShouldBindJSON(&registerForm); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate with forms.UserForm (form)
	userForm := forms.UserForm{}
	if validationError := userForm.Register(registerForm); validationError != "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": validationError,
		})
		return
	}

	// Register(service)
	err := userService.Register(registerForm)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Response with message
	c.JSON(http.StatusOK, gin.H{"message": "Register Success"})

}


func (ctr UserControllers) Logout(c *gin.Context) {
}
