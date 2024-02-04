package services

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/Eco-Led/EcoLed-Back_test/forms"
	"github.com/Eco-Led/EcoLed-Back_test/initializers"

	jwt "github.com/golang-jwt/jwt/v4"
	uuid "github.com/google/uuid"
)

type TokenServices struct{}

func (svc TokenServices) CreateToken(userID int64) (*forms.TokenDetails, error) {
	td := &forms.TokenDetails{}
	td.AtExpires = time.Now().Add(time.Minute * 15).Unix()
	td.AccessUUID = uuid.New().String()

	td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	td.RefreshUUID = uuid.New().String()

	var err error
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = td.AccessUUID
	atClaims["user_id"] = userID
	atClaims["exp"] = td.AtExpires

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return nil, err
	}
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUUID
	rtClaims["user_id"] = userID
	rtClaims["exp"] = td.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
	if err != nil {
		return nil, err
	}
	return td, nil
}

func (svc TokenServices) SaveToken(userid int64, td *forms.TokenDetails) (err error) {
	at := time.Unix(td.AtExpires, 0)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	errAccess := initializers.Redis.Set(td.AccessUUID, strconv.Itoa(int(userid)), at.Sub(now)).Err()
	if errAccess != nil {
		return errAccess
	}
	errRefresh := initializers.Redis.Set(td.RefreshUUID, strconv.Itoa(int(userid)), rt.Sub(now)).Err()
	if errRefresh != nil {
		return errRefresh
	}
	return nil
}

func (svc TokenServices) ExtractTokenID(tokenString string) (int64, error) {
	// Parse token
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})

	// Check token
	if err != nil {
		return 0, err
	}

	// Get userID from token
	userIDValue, ok := claims["user_id"]
	if !ok {
		return 0, errors.New("user_id not found in token")
	}

	// Change type to float64
	userIDFloat, ok := userIDValue.(float64)
	if !ok {
		return 0, errors.New("user_id is not a float64")
	}

	return int64(userIDFloat), nil
}

func (svc TokenServices) DeleteToken(accessTokenString string, RefreshTokenString string) (int64, error) {
	// Parse AccessToken
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(accessTokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})

	// Check token
	if err != nil {
		return 0, err
	}

	// Get accessUUID from token
	accessUUIDValue, ok := claims["access_uuid"]
	if !ok {
		return 0, errors.New("access_uuid not found in token")
	}

	// Change type to string
	accessUUIDString, ok := accessUUIDValue.(string)
	if !ok {
		return 0, errors.New("access_uuid is not a string")
	}

	// Delete Accesstoken
	deleted, err := initializers.Redis.Del(accessUUIDString).Result()
	if err != nil {
		return 0, err
	}

	// Parse RefreshToken
	claims = jwt.MapClaims{}
	_, err = jwt.ParseWithClaims(RefreshTokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("REFRESH_SECRET")), nil
	})

	// Check token
	if err != nil {
		return 0, err
	}

	// Get refreshUUID from token
	refreshUUIDValue, ok := claims["refresh_uuid"]
	if !ok {
		return 0, errors.New("refresh_uuid not found in token")
	}

	// Change type to string
	refreshUUIDString, ok := refreshUUIDValue.(string)
	if !ok {
		return 0, errors.New("refresh_uuid is not a string")
	}

	// Delete Refreshtoken
	deleted, err = initializers.Redis.Del(refreshUUIDString).Result()
	if err != nil {
		return 0, err
	}


	// Return deleted
	return deleted, nil
}
