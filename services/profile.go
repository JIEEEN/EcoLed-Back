package services

import (
	"errors"

	"github.com/Eco-Led/EcoLed-Back_test/forms"
	"github.com/Eco-Led/EcoLed-Back_test/initializers"
	"github.com/Eco-Led/EcoLed-Back_test/models"
)

type ProfileServices struct{}

func (svc ProfileServices) UpdateProfile(userID uint, profileForm forms.ProfileForm) (err error) {
	// Update profile
	result := initializers.DB.Model(&models.Profiles{}).Where("user_id=?", userID).Updates(models.Profiles{
		Nickname:  profileForm.Nickname,
		Age:       profileForm.Age,
		Introduce: profileForm.Introduce,
	})
	if result.Error != nil {
		err := errors.New("failed to update profile")
		return err
	}

	return nil

}

func (svc ProfileServices) GetProfile(userID uint) (profile models.Profiles, err error) {
	// Get profile
	result := initializers.DB.Where("user_id=?", userID).First(&profile)
	if result.Error != nil {
		err := errors.New("failed to get profile")
		return profile, err
	}

	return profile, nil

}

func (svc ProfileServices) GetProfileByNickname(nickname string) (profile models.Profiles, err error) {
	// Get profile
	result := initializers.DB.Where("nickname=?", nickname).First(&profile)
	if result.Error != nil {
		err := errors.New("failed to get profile")
		return profile, err
	}

	return profile, nil

}
