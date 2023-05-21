package database

import (
	"go_bedu/config"
	"go_bedu/models"

	"gorm.io/gorm"
)

type AdminRepository interface {
	ReadToken(id int) (admin *models.Administrator, err error)
	GetAdmins() (admin []models.Administrator, err error)
	GetAdminById(id int) (admin *models.Administrator, err error)
	GetAdminByEmail(email string) (admin *models.Administrator, err error)
	UpdateAdmin(admin *models.Administrator) error
	CreateAdmin(admin *models.Administrator) error
	DeleteAdmin(admin *models.Administrator) error
}

type adminRepository struct {
	db *gorm.DB
}

func NewAdminRepository(db *gorm.DB) *adminRepository {
	return &adminRepository{db}
}

// Read Token is a function to read token
func (a *adminRepository) ReadToken(id int) (admin *models.Administrator, err error) {
	err = config.DB.Where("id = ?", id).First(&admin).Error

	if err != nil {
		return nil, err
	}

	return admin, nil
}

// Get Admins is a function to get all admins
func (a *adminRepository) GetAdmins() (admin []models.Administrator, err error) {
	if err := config.DB.Preload("Articles").Find(&admin).Error; err != nil {
		return nil, err
	}

	return admin, nil
}

// Get Admin By Id is a function to get admin by id
func (a *adminRepository) GetAdminById(id int) (admin *models.Administrator, err error) {
	err = config.DB.Model(&admin).Preload("Articles").Where("id = ?", id).First(&admin).Error
	if err != nil {
		return nil, err
	}

	return admin, nil
}

// Update Admin is a function to update the admin
func (a *adminRepository) UpdateAdmin(admin *models.Administrator) error {
	if err := config.DB.Updates(&admin).Error; err != nil {
		return err
	}

	return nil
}

// Get Admin By Email is a function to get admin by email
func (a *adminRepository) GetAdminByEmail(email string) (admin *models.Administrator, err error) {
	if err = config.DB.Where("email = ?", email).First(&admin).Error; err != nil {
		return nil, err
	}

	return admin, nil
}

// Create Admin is a function to create the admin
func (a *adminRepository) CreateAdmin(admin *models.Administrator) error {
	if err := config.DB.Create(&admin).Error; err != nil {
		return err
	}

	return nil
}

// Delete Admin is a function to delete the admin
func (a *adminRepository) DeleteAdmin(admin *models.Administrator) error {
	if err := config.DB.Delete(&admin).Error; err != nil {
		return err
	}

	return nil
}
