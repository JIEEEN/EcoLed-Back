package forms

import (
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/go-playground/validator/v10"

)

var GoogleClientID = os.Getenv("GOOGLE_CLIENT_ID")
var GoogleClientSecret = os.Getenv("GOOGLE_CLIENT_SECRET")
var GoogleRedirectURL = os.Getenv("GOOGLE_REDIRECT_URL")

var (
	GoogleOauthConfig = &oauth2.Config{
		RedirectURL:  GoogleRedirectURL,
		ClientID:     GoogleClientID,
		ClientSecret: GoogleClientSecret,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
)

type OauthToken struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	Expiry       string `json:"expiry"`
}

type OauthUser struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

type OauthRegisterForm struct {
	Nickname 	string `json:"nickname" binding:"required"`
	Accountname string `json:"accountname" binding:"required"`
}

func (f OauthRegisterForm) Validate() string {
    validate := validator.New()
    err := validate.Struct(f)

    if err != nil {
        for _, err := range err.(validator.ValidationErrors) {
            switch err.Field() {
            case "Nickname":
                return f.NicknameError(err.Tag())
            case "Accountname":
                return f.AccountnameError(err.Tag())
            }
        }
    }
    return ""
}

// Custom validation error messages for Nickname field
func (f OauthRegisterForm) NicknameError(tag string) string {
    switch tag {
    case "required":
        return "Nickname is required"
    default:
        return "Invalid nickname"
    }
}

// Custom validation error messages for AccountName field
func (f OauthRegisterForm) AccountnameError(tag string) string {
    switch tag {
    case "required":
        return "Accountname is required"
    default:
        return "Invalid Accountname"
    }
}
