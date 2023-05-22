package repositories

import (
	"go_bedu/config"
	"go_bedu/models"

	"github.com/jinzhu/gorm"
)

type AuthRepository interface {
	LoginAdmin(admin *models.Administrator) error
}

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *authRepository {
	return &authRepository{db}
}

func (a *authRepository) LoginAdmin(admin *models.Administrator) error {
	if err := config.DB.Where("email = ? AND password = ?", admin.Email, admin.Password).First(&admin).Error; err != nil {
		return err
	}

	return nil
}
