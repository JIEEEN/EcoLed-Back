package services

import (
	"errors"
	"time"

	"github.com/Eco-Led/EcoLed-Back_test/forms"
	"github.com/Eco-Led/EcoLed-Back_test/initializers"
	"github.com/Eco-Led/EcoLed-Back_test/models"
)

type PaylogServices struct{}


func (svc PaylogServices) CreatePaylog(userID uint, paylog forms.PaylogForm) (err error) {
	// Get account
	var account models.Accounts
	result := initializers.DB.First(&account, "user_id=?", userID)
	if result.Error != nil {
		err := errors.New("failed to get account")
		return err
	}

	// Create paylog
	result = initializers.DB.Create(&models.Paylogs{
		Date:       paylog.Date,
		Time:       paylog.Time,
		Content:    paylog.Content,
		Cost:       paylog.Cost,
		Name:       paylog.Name,
		Place:      paylog.Place,
		Material:   paylog.Material,
		Ecoscore:   paylog.Ecoscore,
		Account_id: account.ID,
	})
	if result.Error != nil {
		err := errors.New("failed to create paylog")
		return err
	}

	//Update account
	account.Total_ecoscore += paylog.Ecoscore
	account.Balance -= paylog.Cost
	result = initializers.DB.Save(&account)
	if result.Error != nil {
		err := errors.New("failed to update account")
		return err
	}

	return nil

}


func (svc PaylogServices) GetPaylog(userID uint, paylogID uint) (paylog models.Paylogs, err error) {
	// Get account
	var account models.Accounts
	result := initializers.DB.First(&account, "user_id=?", userID)
	if result.Error != nil {
		err := errors.New("failed to get account")
		return paylog, err
	}

	// Get paylog
	result = initializers.DB.First(&paylog, "id=?", paylogID)
	if result.Error != nil {
		err := errors.New("failed to get paylog")
		return paylog, err
	}

	//Check if paylog is owned by user
	if paylog.Account_id != account.ID {
		err := errors.New("paylog is not owned by user")
		return paylog, err
	}

	return paylog, nil

}


func (svc PaylogServices) GetPaylogs(accountID uint, page int) (paylogs []models.Paylogs, err error) {
	// Setting page
	endDate := time.Now().AddDate(0, 0, -page*31)
	startDate := endDate.AddDate(0, 0, -31)
	endDateStr := endDate.Format("200601021504")
	startDateStr := startDate.Format("200601021504")

	// Get paylogs
	result := initializers.DB.Where("account_id=?", accountID).
		Where("deleted_at IS NULL").
		Where("CONCAT(Date, Time) BETWEEN ? AND ?", startDateStr, endDateStr).
		Order("CONCAT(Date, Time) DESC").
		Find(&paylogs)

	if result.Error != nil {
		err := errors.New("failed to get paylogs")
		return paylogs, err
	}
	if result.RowsAffected == 0 {
		err := errors.New("no paylogs found")
		return paylogs, err
	}

	return paylogs, nil

}

func (svc PaylogServices) UpdatePaylog(userID uint, paylogID uint, paylog forms.PaylogForm) (err error) {
	// Get account
	var account models.Accounts
	result := initializers.DB.First(&account, "user_id=?", userID)
	if result.Error != nil {
		err := errors.New("failed to get account")
		return err
	}

	// Get paylog
	var paylogModel models.Paylogs
	result = initializers.DB.First(&paylogModel, "id=?", paylogID)
	if result.Error != nil {
		err := errors.New("failed to get paylog")
		return err
	}

	//Check if paylog is owned by user
	if paylogModel.Account_id != account.ID {
		err := errors.New("paylog is not owned by user")
		return err
	}

	// Save last ecoscore and cost
	var lastEcoscore float64 = paylogModel.Ecoscore
	var lastCost int64 = paylogModel.Cost

	// Update paylog
	paylogModel.Date = paylog.Date
	paylogModel.Time = paylog.Time
	paylogModel.Content = paylog.Content
	paylogModel.Cost = paylog.Cost
	paylogModel.Name = paylog.Name
	paylogModel.Place = paylog.Place
	paylogModel.Material = paylog.Material
	paylogModel.Ecoscore = paylog.Ecoscore
	result = initializers.DB.Save(&paylogModel)
	if result.Error != nil {
		err := errors.New("failed to update paylog")
		return err
	}

	//Update account
	account.Total_ecoscore -= lastEcoscore
	account.Balance += lastCost

	account.Total_ecoscore += paylog.Ecoscore
	account.Balance -= paylog.Cost

	result = initializers.DB.Save(&account)
	if result.Error != nil {
		err := errors.New("failed to update account")
		return err
	}

	return nil

}

func (svc PaylogServices) DeletePaylog(userID uint, paylogID uint) (err error) {
	// Get account
	var account models.Accounts
	result := initializers.DB.First(&account, "user_id=?", userID)
	if result.Error != nil {
		err := errors.New("failed to get account")
		return err
	}

	// Get paylog
	var paylog models.Paylogs
	result = initializers.DB.First(&paylog, "id=?", paylogID)
	if result.Error != nil {
		err := errors.New("failed to get paylog")
		return err
	}

	//Check if paylog is owned by user
	if paylog.Account_id != account.ID {
		return errors.New("paylog is not owned by user")
	}
	// Delete paylog
	result = initializers.DB.Delete(&paylog)
	if result.Error != nil {
		err := errors.New("failed to delete paylog")
		return err
	}

	//Update account
	account.Total_ecoscore -= paylog.Ecoscore
	account.Balance += paylog.Cost
	result = initializers.DB.Save(&account)
	if result.Error != nil {
		err := errors.New("failed to update account")
		return err
	}

	return nil

}
