package services

import (
	"errors"
	"time"

	"github.com/Eco-Led/EcoLed-Back_test/forms"
	"github.com/Eco-Led/EcoLed-Back_test/initializers"
	"github.com/Eco-Led/EcoLed-Back_test/models"

	"golang.org/x/crypto/bcrypt"
)

type UserServices struct{}

func (svc UserServices) Login(loginForm forms.LoginForm) (user forms.UserReturnForm, token forms.Token, err error) {
	//call by value (not call by reference)
	var userModel = models.Users{}
	var profileModel = models.Profiles{}

	//call by reference
	var tokenService = new(TokenServices)

	//From controller, binding value is received. So, check whether the value is valid.
	initializers.DB.First(&userModel, "email=?", loginForm.Email)
	initializers.DB.First(&profileModel, "user_id=?", userModel.ID)

	// If the value is not valid, return error
	if userModel.ID == 0 || profileModel.ID == 0 {
		err := errors.New("data does not exist in db")
		return user, token, err
	}

	// Set return value (user)
	user = forms.UserReturnForm{
		Email:     userModel.Email,
		Nickname:  profileModel.Nickname,
		CreatedAt: userModel.CreatedAt.String(),
	}

	// Compare password
	err = bcrypt.CompareHashAndPassword([]byte(userModel.Password), []byte(loginForm.Password))
	if err != nil {
		err := errors.New("invalid password")
		return user, token, err
	}

	// Create token
	td, err := tokenService.CreateToken(int64(userModel.ID))
	if err != nil {
		return user, token, err
	}

	// Save token
	err = tokenService.SaveToken(int64(userModel.ID), td)
	if err != nil {
		return user, token, err
	}

	// Set return value (token)
	token = forms.Token{
		AccessToken:  td.AccessToken,
		RefreshToken: td.RefreshToken,
	}

	// Return user, token
	return user, token, err

}

func (svc UserServices) Register(registerForm forms.RegisterForm) error {
	//call by value (not call by reference)
	var userModel = models.Users{}
	var profileModel = models.Profiles{}

	//Start transaction
	tx := initializers.DB.Begin()

	// Check whether the email is unique
	tx.First(&userModel, "email=?", registerForm.Email)
	if userModel.ID != 0 {
		err := errors.New("email already exists")
		return err
	}

	//Check whether the nickname is unique
	tx.First(&profileModel, "nickname=?", registerForm.Nickname)
	if profileModel.ID != 0 {
		err := errors.New("nickname already exists")
		return err
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerForm.Password), bcrypt.DefaultCost)
	if err != nil {
		err := errors.New("failed to hash password")
		return err
	}

	// Create user
	user := models.Users{
		Email:    registerForm.Email,
		Password: string(hashedPassword),
	}
	result := tx.Create(&user)
	if result.Error != nil {
		tx.Rollback()
		err := errors.New("failed to create user")
		return err
	}

	// Create profile
	result = tx.Create(&models.Profiles{
		Nickname: registerForm.Nickname,
		User_id:  user.ID,
	})
	if result.Error != nil {
		tx.Rollback()
		err := errors.New("failed to create profile")
		return err
	}

	// Create account
	result = tx.Create(&models.Accounts{
		Name:    registerForm.Accountname,
		User_id: user.ID,
	})
	if result.Error != nil {
		tx.Rollback()
		err := errors.New("failed to create account")
		return err
	}

	return tx.Commit().Error

}

func (svc UserServices) Logout(accessToken string, refreshToken string) (err error) {
	//call by reference
	var tokenService = new(TokenServices)

	// Delete token
	_, err = tokenService.DeleteToken(accessToken, refreshToken)
	if err != nil {
		return err
	}

	// Return error
	return nil
}

// for verifying email error
var ErrEmailNotFound = errors.New("email not found")

func (svc UserServices) FindEmail(email string) error {
	//call by value (not call by reference)
	var userModel = models.Users{}

	// Check whether the email exists
	initializers.DB.First(&userModel, "email=?", email)
	if userModel.ID != 0 {
		return nil
	}

	// Return error
	return ErrEmailNotFound
}

func (svc UserServices) SaveVerificationCode(email string, code string) error {
	// Save verification code
	err := initializers.Redis.Set("authCode:"+email, code, 5*time.Minute).Err()
	if err != nil {
		return err
	}

	return nil
}

func (svc UserServices) VerifyCode(email string, code string) error {
	// Get verification code
	savedCode, err := initializers.Redis.Get("authCode:" + email).Result()
	if err != nil {
		return err
	}

	// Compare verification code
	if savedCode != code {
		return errors.New("invalid verification code")
	}

	// Return error
	return nil
}

func (svc UserServices) UpdatePassword(email string, password string) error {
	//call by value (not call by reference)
	var userModel = models.Users{}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		err := errors.New("failed to hash password")
		return err
	}

	// Update password
	err = initializers.DB.Model(&userModel).Where("email=?", email).Update("password", hashedPassword).Error
	if err != nil {
		return err
	}

	// Return error
	return nil
}