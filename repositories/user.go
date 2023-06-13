package repositories

import (
	"go_bedu/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	LoginUser(user models.User) error
	ReadToken(id uint) (models.User, error)
	GetUserByVerificationCode(verificationCode any) (models.User, error)
	GetUserOTP(otp int) (user models.User, err error)
	GetUserById(id uint) (models.User, error)
	GetUserByEmail(email string) (user models.User, err error)
	GetUserByUsername(username string) (user models.User, err error)
	GetUsers() ([]models.User, error)
	CreateUser(user models.User) (models.User, error)
	UpdateUser(user models.User) (models.User, error)
	DeleteUser(user models.User) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{db}
}

func (r *userRepository) LoginUser(user models.User) error {
	err := r.db.Where("username = ? AND password = ?", user.Username, user.Password).Error

	return err
}

func (r *userRepository) ReadToken(id uint) (models.User, error) {
	var user models.User

	err := r.db.Where("id = ?", user.ID).Error

	return user, err
}

func (r *userRepository) GetUserByVerificationCode(verificationCode any) (models.User, error) {
	var user models.User

	err := r.db.Where("verification_code = ?", verificationCode).First(&user).Error

	return user, err
}

func (r *userRepository) GetUserOTP(otp int) (user models.User, err error) {
	var verified = true
	err = r.db.Where("otp = ? AND verified = ?", otp, verified).First(&user).Error

	return user, err
}

func (r *userRepository) GetUserById(id uint) (models.User, error) {
	var user models.User

	err := r.db.Model(&user).Where("id = ?", id).First(&user).Error

	return user, err
}

func (r *userRepository) GetUserByEmail(email string) (user models.User, err error) {
	err = r.db.Model(&user).Where("email = ? AND deleted_at IS NULL", email).First(&user).Error

	return user, err
}

func (r *userRepository) GetUserByUsername(username string) (user models.User, err error) {
	err = r.db.Model(&user).Where("username = ? AND deleted_at IS NULL", username).First(&user).Error

	return user, err
}

func (r *userRepository) GetUsers() ([]models.User, error) {
	var user []models.User
	err := r.db.Find(&user).Error

	if err != nil {
		return user, err
	}

	return user, err
}

func (r *userRepository) CreateUser(user models.User) (models.User, error) {
	err := r.db.Create(&user).Error

	return user, err
}

func (r *userRepository) UpdateUser(user models.User) (models.User, error) {
	err := r.db.Table("users").Save(&user).Error

	return user, err
}

func (r *userRepository) DeleteUser(user models.User) error {
	err := r.db.Delete(&user).Error

	return err
}
