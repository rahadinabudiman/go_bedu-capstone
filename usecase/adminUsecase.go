package usecase

import (
	"errors"
	"go_bedu/models"
	"go_bedu/models/payload"
	"go_bedu/repository/database"

	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"
)

type AdminUsecase interface {
	GetAdmin() ([]models.Administrator, error)
	GetAdminById(id int) (res payload.AdminProfileResponse, err error)
	UpdateAdmin(id int, req *payload.UpdateAdminRequest) (res payload.UpdateAdminResponse, err error)
	CreateAdmin(req *payload.RegisterAdminRequest) error
	DeleteAdmin(id int, password string) error
}

type adminUsecase struct {
	adminRepository database.AdminRepository
}

func NewAdminUsecase(adminRepository database.AdminRepository) *adminUsecase {
	return &adminUsecase{adminRepository}
}

// Logic for get All Admin
func (a *adminUsecase) GetAdmin() ([]models.Administrator, error) {
	admin, err := a.adminRepository.GetAdmins()

	if err != nil {
		return nil, err
	}

	return admin, nil
}

// Logic for get Admin with Cookie
func (a *adminUsecase) GetAdminById(id int) (res payload.AdminProfileResponse, err error) {
	admin, err := a.adminRepository.GetAdminById(id)

	if err != nil {
		echo.NewHTTPError(401, "This routes for admin only")
		return
	}

	res = payload.AdminProfileResponse{
		ID:    admin.ID,
		Nama:  admin.Nama,
		Email: admin.Email,
		Role:  admin.Role,
	}

	return res, nil
}

// Logic for Update Admin
func (a *adminUsecase) UpdateAdmin(id int, req *payload.UpdateAdminRequest) (res payload.UpdateAdminResponse, err error) {
	adminRequest := &models.Administrator{
		Nama:     req.Nama,
		Email:    req.Email,
		Password: req.Password,
	}

	adminRequest.ID = uint(id)

	err = a.adminRepository.UpdateAdmin(adminRequest)

	if err != nil {
		echo.NewHTTPError(400, "Failed to update Admin")
		return
	}

	res = payload.UpdateAdminResponse{
		Nama:     adminRequest.Nama,
		Email:    adminRequest.Email,
		Password: adminRequest.Password,
	}

	return
}

// Logic for Create Admin
func (a *adminUsecase) CreateAdmin(req *payload.RegisterAdminRequest) error {
	adminRequest := &models.Administrator{
		Nama:     req.Nama,
		Email:    req.Email,
		Password: req.Password,
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(adminRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		return echo.NewHTTPError(400, "Failed to hash password")
	}

	adminRequest.Password = string(passwordHash)

	err = a.adminRepository.CreateAdmin(adminRequest)

	if err != nil {
		return errors.New(err.Error())
	}

	return nil
}

// Logic for Delete Administrator
func (a *adminUsecase) DeleteAdmin(id int, password string) error {
	admin, err := a.adminRepository.ReadToken(id)

	if err != nil {
		return echo.NewHTTPError(400, "Failed to get Admin")
	}

	err = bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return echo.NewHTTPError(400, err.Error())
	}

	err = a.adminRepository.DeleteAdmin(admin)

	if err != nil {
		return err
	}

	return nil
}
