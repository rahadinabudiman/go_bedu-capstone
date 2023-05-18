package database

import (
	"go_bedu/config"
	"go_bedu/models"
)

// Create Administrator to DB
func CreateAdministrator(admin models.Administrator) (models.Administrator, error) {
	err := config.DB.Create(&admin).Error

	if err != nil {
		return models.Administrator{}, err
	}

	return admin, nil
}

// Get All Admins from DB
func GetAdministrators() (admin []models.Administrator, err error) {
	err = config.DB.Find(&admin).Error

	if err != nil {
		return []models.Administrator{}, err
	}
	return admin, nil
}
