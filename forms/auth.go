package forms

import ()

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUUID   string
	RefreshUUID  string
	AtExpires    int64
	RtExpires    int64
}

type AccessDetails struct {
	AccessUUID string
	UserID     int64
}

type AuthForm struct {}

type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshToken struct {
	// AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (f AuthForm) RefreshToken(refreshToken RefreshToken) string {
	if refreshToken.RefreshToken == "" {
		return "Refresh token is required"
	}
	return ""
}
