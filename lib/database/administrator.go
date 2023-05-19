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

// Get Admin by ID from DB
func GetAdministratorByID(id int) (admin models.Administrator, err error) {
	err = config.DB.Where("id = ?", id).Find(&admin).Error

	if err != nil {
		return models.Administrator{}, err
	}
	return admin, nil
}

// Update Admin by ID from DB
func UpdateAdministrator(admin models.Administrator, id int) (models.Administrator, error) {
	err := config.DB.Table("administrators").Where("id = ?", id).Updates(&admin).Error

	if err != nil {
		return models.Administrator{}, err
	}

	return admin, nil
}

// Delete Admin by ID from DB
func DeleteAdministrator(id int) (interface{}, error) {
	err := config.DB.Where("id = ?", id).Delete(&models.Administrator{}).Error

	if err != nil {
		return nil, err
	}

	return "Administrator behasil dihapus", nil
}

// Login Admin by Email and Password from DB
func LoginAdministrator(admin models.Administrator) (models.Administrator, error) {
	err := config.DB.Where("email = ? AND password = ?", admin.Email, admin.Password).First(&admin).Error

	if err != nil {
		return models.Administrator{}, err
	}

	return admin, nil
}
