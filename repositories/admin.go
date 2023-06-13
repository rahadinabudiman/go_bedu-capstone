package repositories

import (
	"go_bedu/models"

	"gorm.io/gorm"
)

type AdminRepository interface {
	LoginAdmin(admin models.Administrator) error
	ReadToken(id uint) (admin models.Administrator, err error)
	GetAdminByVerificationCode(verificationCode any) (admin models.Administrator, err error)
	GetAdminOTP(otp int) (admin models.Administrator, err error)
	GetAdmins() ([]models.Administrator, error)
	GetAdminById(id uint) (models.Administrator, error)
	GetAdminByEmail(email string) (admin models.Administrator, err error)
	GetAdminByUsername(username string) (admin models.Administrator, err error)
	UpdateAdmin(admin models.Administrator) (models.Administrator, error)
	CreateAdmin(admin models.Administrator) (models.Administrator, error)
	DeleteAdmin(admin models.Administrator) error
}

type adminRepository struct {
	db *gorm.DB
}

func NewAdminRepository(db *gorm.DB) *adminRepository {
	return &adminRepository{db}
}

// Get Admin by Verification Code
func (r *adminRepository) GetAdminByVerificationCode(verificationCode any) (models.Administrator, error) {
	var admin models.Administrator

	err := r.db.Where("verification_code = ?", verificationCode).First(&admin).Error

	return admin, err
}

// Get Admin by OTP
func (r *adminRepository) GetAdminOTP(otp int) (admin models.Administrator, err error) {
	var verified = true
	err = r.db.Where("otp = ? AND verified = ?", otp, verified).First(&admin).Error

	return admin, err
}

// Login Administrator from Database
func (r *adminRepository) LoginAdmin(admin models.Administrator) error {
	err := r.db.Where("username = ? AND password = ?", admin.Email, admin.Password).First(&admin).Error

	return err
}

// Read Token is a function to read token
func (r *adminRepository) ReadToken(id uint) (models.Administrator, error) {
	// admin := &models.Administrator{} // Menggunakan objek struct daripada pointer
	var admin models.Administrator

	err := r.db.Where("id = ?", id).First(&admin).Error

	return admin, err
}

// Get Admins is a function to get all admins
func (r *adminRepository) GetAdmins() ([]models.Administrator, error) {
	var admin []models.Administrator
	err := r.db.Preload("Articles").Find(&admin).Error

	if err != nil {
		return admin, err
	}

	return admin, err
}

// Get Admin By Id is a function to get admin by id
func (r *adminRepository) GetAdminById(id uint) (models.Administrator, error) {
	var admin models.Administrator

	err := r.db.Model(&admin).Preload("Articles").Where("id = ?", id).First(&admin).Error
	return admin, err
}

// Update Admin is a function to update the admin
func (r *adminRepository) UpdateAdmin(admin models.Administrator) (models.Administrator, error) {
	err := r.db.Table("administrators").Save(&admin).Error

	return admin, err
}

// Get Admin By Email is a function to get admin by email
func (r *adminRepository) GetAdminByEmail(email string) (models.Administrator, error) {
	var admin models.Administrator

	err := r.db.Where("email = ? AND deleted_at IS NULL", email).First(&admin).Error

	return admin, err
}

// Get Admin By Username is a function to get admin by username
func (r *adminRepository) GetAdminByUsername(username string) (models.Administrator, error) {
	var admin models.Administrator

	err := r.db.Where("Username = ? AND deleted_at IS NULL", username).First(&admin).Error

	return admin, err
}

// Get Admin By Email is a function to get admin by email
func (r *adminRepository) GetAdminByPassword(password string) (models.Administrator, error) {
	// admin := &models.Administrator{} // Menggunakan objek struct daripada pointer
	var admin models.Administrator

	err := r.db.Where("password = ?", password).First(&admin).Error

	return admin, err
}

// Create Admin is a function to create the admin
func (r *adminRepository) CreateAdmin(admin models.Administrator) (models.Administrator, error) {
	err := r.db.Create(&admin).Error

	return admin, err
}

// Delete Admin is a function to delete the admin
func (r *adminRepository) DeleteAdmin(admin models.Administrator) error {
	err := r.db.Delete(&admin).Error

	return err
}
