package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"

	"github.com/Eco-Led/EcoLed-Back_test/forms"
	"github.com/Eco-Led/EcoLed-Back_test/initializers"
	"github.com/Eco-Led/EcoLed-Back_test/models"
	"golang.org/x/oauth2"
)

type OauthServices struct{}

func (svc OauthServices) GenerateStateString() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func (svc OauthServices) GetLoginURL(state string) (url string) {
	GoogleOauthConfig := forms.GetGoogleOauthConfig()
	url = GoogleOauthConfig.AuthCodeURL(state)
	return url
}

func (svc OauthServices) GetOauthToken(code string) (*oauth2.Token, error) {
	GoogleOauthConfig := forms.GetGoogleOauthConfig()
	token, err := GoogleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (svc OauthServices) GetOauthUser(token *oauth2.Token) (forms.OauthUser, error) {
	GoogleOauthConfig := forms.GetGoogleOauthConfig()
	client := GoogleOauthConfig.Client(context.Background(), token)
	response, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return forms.OauthUser{}, err
	}
	defer response.Body.Close() // 응답 본문을 함수 종료 시 자동으로 닫도록 함

	userInfo := forms.OauthUser{}
	if err := json.NewDecoder(response.Body).Decode(&userInfo); err != nil {
		return forms.OauthUser{}, err
	}

	return userInfo, nil
}

func (svc OauthServices) FindUserExists(userInfo forms.OauthUser) bool {
	var userModel = models.Users{}
	initializers.DB.First(&userModel, "email=?", userInfo.Email)

	if userModel.ID != 0 {
		return true
	} else {
		return false
	}
}

func (svc OauthServices) SaveOauthUser(oauthToken *oauth2.Token, userInfo forms.OauthUser) (forms.UserReturnForm, forms.Token, error) {
	var userModel = models.Users{}
	var oauthModel = models.OAuth{}
	var profileModel = models.Profiles{}

	// get user from DB
	initializers.DB.First(&userModel, "email=?", userInfo.Email)

	// update oauth from DB
	initializers.DB.First(&oauthModel, "user_id=?", userModel.ID)
	oauthModel.Access_token = oauthToken.AccessToken
	oauthModel.Refresh_token = oauthToken.RefreshToken
	oauthModel.Expiry = oauthToken.Expiry
	initializers.DB.Save(&oauthModel)

	// get profile from DB
	initializers.DB.First(&profileModel, "user_id=?", userModel.ID)

	// return user
	user := forms.UserReturnForm{
		Email:     userModel.Email,
		Nickname:  profileModel.Nickname,
		CreatedAt: userModel.CreatedAt.String(),
	}

	// Create token
	var tokenService = TokenServices{}
	td, err := tokenService.CreateToken(int64(userModel.ID))
	if err != nil {
		return user, forms.Token{}, err
	}

	// Save token
	err = tokenService.SaveToken(int64(userModel.ID), td)
	if err != nil {
		return user, forms.Token{}, err
	}

	// Set return value (token)
	token := forms.Token{
		AccessToken:  td.AccessToken,
		RefreshToken: td.RefreshToken,
	}

	return user, token, nil
}

func (svc OauthServices) Register(registerForm forms.OauthRegisterForm, userInfo forms.OauthUser, oauthToken *oauth2.Token) (forms.UserReturnForm, forms.Token, error) {
	var userModel = models.Users{}
	var oauthModel = models.OAuth{}
	var profileModel = models.Profiles{}
	var accountModel = models.Accounts{}

	// Create user
	userModel.Email = userInfo.Email
	userModel.Password = ""
	userModel.Oauth_provider = "Google"
	userModel.Oauth_id = userInfo.ID
	initializers.DB.Create(&userModel)

	// Create oauth
	oauthModel.User_id = userModel.ID
	oauthModel.Provider = "Google"
	oauthModel.Provider_userid = userInfo.ID
	oauthModel.Access_token = oauthToken.AccessToken
	oauthModel.Refresh_token = oauthToken.RefreshToken
	oauthModel.Expiry = oauthToken.Expiry
	initializers.DB.Create(&oauthModel)

	// Create profile
	profileModel.User_id = userModel.ID
	profileModel.Nickname = registerForm.Nickname
	initializers.DB.Create(&profileModel)

	// Create account
	accountModel.User_id = userModel.ID
	accountModel.Name = registerForm.Accountname
	initializers.DB.Create(&accountModel)

	// Return user
	user := forms.UserReturnForm{
		Email:     userModel.Email,
		Nickname:  profileModel.Nickname,
		CreatedAt: userModel.CreatedAt.String(),
	}

	// Create token
	var tokenService = TokenServices{}
	td, err := tokenService.CreateToken(int64(userModel.ID))
	if err != nil {
		return user, forms.Token{}, err
	}

	// Save token
	err = tokenService.SaveToken(int64(userModel.ID), td)
	if err != nil {
		return user, forms.Token{}, err
	}

	// Set return value (token)
	token := forms.Token{
		AccessToken:  td.AccessToken,
		RefreshToken: td.RefreshToken,
	}

	return user, token, nil
}
