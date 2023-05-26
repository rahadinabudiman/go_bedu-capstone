package repositories

import (
	"go_bedu/models"

	"github.com/jinzhu/gorm"
)

type AdminRepository interface {
	LoginAdmin(admin *models.Administrator) error
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

// Login Administrator from Database
func (r *adminRepository) LoginAdmin(admin *models.Administrator) error {
	if err := r.db.Where("email = ? AND password = ?", admin.Email, admin.Password).First(&admin).Error; err != nil {
		return err
	}

	return nil
}

// Read Token is a function to read token
func (r *adminRepository) ReadToken(id int) (*models.Administrator, error) {
	admin := &models.Administrator{} // Menggunakan objek struct daripada pointer

	err := r.db.Where("id = ?", id).First(admin).Error

	if err != nil {
		return nil, err
	}

	return admin, nil
}

// Get Admins is a function to get all admins
func (r *adminRepository) GetAdmins() (admin []models.Administrator, err error) {
	if err := r.db.Preload("Articles").Find(&admin).Error; err != nil {
		return nil, err
	}

	return admin, nil
}

// Get Admin By Id is a function to get admin by id
func (r *adminRepository) GetAdminById(id int) (*models.Administrator, error) {
	admin := &models.Administrator{}

	err := r.db.Model(&admin).Preload("Articles").Where("id = ?", id).First(&admin).Error
	if err != nil {
		return nil, err
	}

	return admin, nil
}

// Update Admin is a function to update the admin
func (r *adminRepository) UpdateAdmin(admin *models.Administrator) error {
	if err := r.db.Table("administrators").Save(admin).Error; err != nil {
		return err
	}

	return nil
}

// Get Admin By Email is a function to get admin by email
func (r *adminRepository) GetAdminByEmail(email string) (*models.Administrator, error) {
	admin := &models.Administrator{} // Menggunakan objek struct daripada pointer

	if err := r.db.Where("email = ?", email).First(admin).Error; err != nil {
		return nil, err
	}

	return admin, nil
}

// Get Admin By Email is a function to get admin by email
func (r *adminRepository) GetAdminByPassword(password string) (*models.Administrator, error) {
	admin := &models.Administrator{} // Menggunakan objek struct daripada pointer

	if err := r.db.Where("password = ?", password).First(admin).Error; err != nil {
		return nil, err
	}

	return admin, nil
}

// Create Admin is a function to create the admin
func (r *adminRepository) CreateAdmin(admin *models.Administrator) error {
	if err := r.db.Create(&admin).Error; err != nil {
		return err
	}

	return nil
}

// Delete Admin is a function to delete the admin
func (r *adminRepository) DeleteAdmin(admin *models.Administrator) error {
	if err := r.db.Delete(&admin).Error; err != nil {
		return err
	}

	return nil
}
