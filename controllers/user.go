package controllers

import (
	"errors"
	"net/http"

	"github.com/Eco-Led/EcoLed-Back_test/forms"
	"github.com/Eco-Led/EcoLed-Back_test/services"

	"github.com/gin-gonic/gin"
)

type UserControllers struct{}

var userService = new(services.UserServices)
var emailService = new(services.EmailServices)

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
	// Get accesstoken from header
	accessToken := c.Request.Header.Get("Authorization")
	if accessToken == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "token not found"})
		return
	}

	// Get refresh token from body
	refreshToken := forms.RefreshToken{}
	if err := c.ShouldBindJSON(&refreshToken); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate with forms.UserForm (form)
	authForm := forms.AuthForm{}
	if validationError := authForm.RefreshToken(refreshToken); validationError != "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": validationError,
		})
		return
	}

	// Logout(service)
	err := userService.Logout(accessToken, refreshToken.RefreshToken)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Response with message
	c.JSON(http.StatusOK, gin.H{"message": "Logout Success"})

}

func (ctr UserControllers) FindEmail(c *gin.Context) {
	// Bind JSON with forms.FindEmailForm (form)
	var findUserInfoForm forms.FindUserInfoForm
	if err := c.ShouldBindJSON(&findUserInfoForm); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate with forms.UserForm (form)
	userForm := forms.UserForm{}
	if validation := userForm.FindUserInfo(findUserInfoForm); validation != "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": validation})
		return
	}

	// FindEmail(service)
	err := userService.FindEmail(findUserInfoForm.Email)
	if err != nil {
		// Email is not found in DB
		if errors.Is(err, services.ErrEmailNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Email not found", "redirect": "/findEmailPage1"})
		} else {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}

	// Email is found in DB
	c.JSON(http.StatusOK, gin.H{"message": "Email found", "redirect": "/findEmailPage2"})

}

func (ctr UserControllers) FindPassword(c *gin.Context) {
	// Bind JSON with forms.FindEmailForm (form)
	var findUserInfoForm forms.FindUserInfoForm
	if err := c.ShouldBindJSON(&findUserInfoForm); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate with forms.UserForm (form)
	userForm := forms.UserForm{}
	if validation := userForm.FindUserInfo(findUserInfoForm); validation != "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": validation})
		return
	}

	// FindEmail(service)
	err := userService.FindEmail(findUserInfoForm.Email)
	if err != nil {
		// Email is not found in DB
		if errors.Is(err, services.ErrEmailNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Email not found", "redirect": "/findEmailPage1"})
		} else {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}

	// Email is found in DB
	// SendEmail(service)
	verificationCode, err := emailService.SendVerifyingEmail("EcoLed 비밀번호 찾기 인증 이메일 입니다.", "./verifyingEmail.html", []string{findUserInfoForm.Email})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = userService.SaveVerificationCode(findUserInfoForm.Email, verificationCode)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Response with message
	c.JSON(http.StatusOK, gin.H{"message": "Send Email Success", "redirect": "/verifycode"})

}

func (ctr UserControllers) VerifyCode(c *gin.Context) {
	// Bind JSON with forms.VerifyCodeForm (form)
	var verifyCodeForm forms.VerifyCodeForm
	if err := c.ShouldBindJSON(&verifyCodeForm); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate with forms.UserForm (form)
	userForm := forms.UserForm{}
	if validation := userForm.VerifyCode(verifyCodeForm); validation != "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": validation})
		return
	}

	// VerifyCode(service)
	err := userService.VerifyCode(verifyCodeForm.Email, verifyCodeForm.Code)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Response with message
	c.JSON(http.StatusOK, gin.H{"message": "Verify Code Success", "redirect": "/password"})

}

func (ctr UserControllers) UpdatePassword(c *gin.Context) {
	// Bind JSON with forms.UpdatePasswordForm (form)
	var loginForm forms.LoginForm
	if err := c.ShouldBindJSON(&loginForm); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate with forms.UserForm (form)
	userForm := forms.UserForm{}
	if validation := userForm.Login(loginForm); validation != "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": validation})
		return
	}

	// UpdatePassword(service)
	err := userService.UpdatePassword(loginForm.Email, loginForm.Password)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Response with message
	c.JSON(http.StatusOK, gin.H{"message": "Update Password Success", "redirect": "/login"})

}
