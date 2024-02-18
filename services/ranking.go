package services

import (
	"errors"

	"github.com/Eco-Led/EcoLed-Back_test/forms"
	"github.com/Eco-Led/EcoLed-Back_test/initializers"
	"github.com/Eco-Led/EcoLed-Back_test/models"
)

type RankingServices struct{}

func (svc RankingServices) GetRanking() (ranking []forms.RankingReturnForm, err error) {
	// Sort accounts by total_ecoscore
	var accounts []models.Accounts
	result := initializers.DB.Order("total_ecoscore desc").Find(&accounts)
	if result.Error != nil {
		err := errors.New("failed to get accounts")
		return ranking, err
	}

	// Get rankingReturnForm date from sorted accounts
	var currentRank int
	var lastEcoscore float64 = -1
	for i, account := range accounts {
		var profile models.Profiles
		result := initializers.DB.Where("user_id = ?", account.User_id).First(&profile)

		if result.Error != nil {
			// err := errors.New("failed to get profile")
			continue
		}

		// To distinguish same ecoscore
		if account.Total_ecoscore != lastEcoscore {
			currentRank = i + 1
			lastEcoscore = account.Total_ecoscore
		}

		// Append rankingReturnForm
		ranking = append(ranking, forms.RankingReturnForm{
			Nickname:      profile.Nickname,
			Age:           int64(profile.Age),
			Profile_image:  profile.Profile_image,
			TotalEcoscore: account.Total_ecoscore,
			Rank:          int64(currentRank),
		})

	}

	return ranking, nil

}
