package repositories

import (
	"go_bedu/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestAdminRepository(t *testing.T) {
	// Setup
	dsn := "r4ha:kmoonkinan@tcp(localhost:3306)/go_bedu?charset=utf8mb4&parseTime=true&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	assert.NoError(t, err)

	adminRepository := NewAdminRepository(db)

	// Running Test Case
	t.Run("Test CreateAdmin", func(t *testing.T) {
		// Create a test admin
		admin := models.Administrator{
			Nama:         "AdminTest",
			Email:        "admin@example.com",
			Password:     "password",
			Verified:     true,
			PhotoProfile: "profile-default.jpg",
			Username:     "admintest",
		}

		// Create the admin in the database
		createdAdmin, err := adminRepository.CreateAdmin(admin)
		assert.NoError(t, err)
		assert.NotZero(t, createdAdmin.ID)

		// Verify that the admin exists in the database
		dbAdmin, err := adminRepository.GetAdminById(uint(createdAdmin.ID))
		assert.NoError(t, err)
		assert.Equal(t, createdAdmin.ID, dbAdmin.ID)
	})

	t.Run("Test GetAdminByEmail", func(t *testing.T) {
		// Get an admin by email
		email := "admin@example.com"
		admin, err := adminRepository.GetAdminByEmail(email)
		assert.NoError(t, err)
		assert.Equal(t, email, admin.Email)
	})

	t.Run("Test UpdateAdmin", func(t *testing.T) {
		// Get an admin to update
		email := "admin@example.com"
		admin, err := adminRepository.GetAdminByEmail(email)
		assert.NoError(t, err)

		// Update the admin's password
		newPassword := "newpassword"
		admin.Password = newPassword
		updatedAdmin, err := adminRepository.UpdateAdmin(admin)
		assert.NoError(t, err)
		assert.Equal(t, newPassword, updatedAdmin.Password)

		// Verify that the admin's password is updated in the database
		dbAdmin, err := adminRepository.GetAdminById(admin.ID)
		assert.NoError(t, err)
		assert.Equal(t, newPassword, dbAdmin.Password)
	})

	t.Run("Test DeleteAdmin", func(t *testing.T) {
		// Get an admin to delete
		email := "admin@example.com"
		admin, err := adminRepository.GetAdminByEmail(email)
		assert.NoError(t, err)

		// Delete admin
		err = adminRepository.DeleteAdmin(admin)
		assert.NoError(t, err)

		// Verify that the admin is deleted from the database
		deletedAdmin, err := adminRepository.GetAdminById(admin.ID)
		assert.Error(t, err)
		assert.Equal(t, models.Administrator{}, deletedAdmin)
	})

}
