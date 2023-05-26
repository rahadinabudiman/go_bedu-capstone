package usecase

import (
	"errors"
	"go_bedu/dtos"
	"go_bedu/helpers"
	"go_bedu/middlewares"
	"go_bedu/models"
	"go_bedu/repositories"

	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"
)

type AdminUsecase interface {
	LoginAdmin(c echo.Context, req *dtos.LoginRequest) (res dtos.LoginResponse, err error)
	GetAdmin() ([]models.Administrator, error)
	GetAdminById(id int) (res dtos.AdminProfileResponse, err error)
	UpdateAdmin(id int, req *dtos.UpdateAdminRequest) (res dtos.UpdateAdminResponse, err error)
	CreateAdmin(req *dtos.RegisterAdminRequest) error
	DeleteAdmin(id int, req *dtos.DeleteAdminRequest) (res helpers.ResponseMessage, err error)
}

type adminUsecase struct {
	adminRepository repositories.AdminRepository
}

func NewAdminUsecase(adminRepository repositories.AdminRepository) *adminUsecase {
	return &adminUsecase{adminRepository}
}

// Logic for get All Admin
func (u *adminUsecase) GetAdmin() ([]models.Administrator, error) {
	admin, err := u.adminRepository.GetAdmins()

	if err != nil {
		return nil, err
	}

	return admin, nil
}

// Logic for get Admin with Cookie
func (u *adminUsecase) GetAdminById(id int) (res dtos.AdminProfileResponse, err error) {
	admin, err := u.adminRepository.GetAdminById(id)

	if err != nil {
		echo.NewHTTPError(401, "This routes for admin only")
		return
	}

	res = dtos.AdminProfileResponse{
		ID:    admin.ID,
		Nama:  admin.Nama,
		Email: admin.Email,
		Role:  admin.Role,
	}

	return res, nil
}

// Logic for Update Admin
func (u *adminUsecase) UpdateAdmin(id int, req *dtos.UpdateAdminRequest) (res dtos.UpdateAdminResponse, err error) {
	adminRequest := &models.Administrator{
		Nama:     req.Nama,
		Email:    req.Email,
		Password: req.Password,
		Role:     req.Role,
	}

	// Check Role and save role information from JWT Cookie
	admin, err := u.adminRepository.ReadToken(id)
	if err != nil {
		echo.NewHTTPError(400, "Failed to get Admin")
		return
	}

	// Check if admin is not Super Admin and trying to change role
	if admin.Role != req.Role && admin.Role != "Super Admin" {
		return res, echo.NewHTTPError(401, "You are not allowed to change role")
	}

	// Check if Super Admin Change other Role
	if admin.Role == "Super Admin" {
		if req.Role != "Super Admin" && req.Role != "Admin" {
			return res, echo.NewHTTPError(401, "Tidak bisa mengubah role menjadi selain Super Admin atau Admin")
		}
	}

	adminRequest.ID = uint(id)

	passwordHash, err := helpers.HashPassword(adminRequest.Password)
	if err != nil {
		echo.NewHTTPError(400, "Failed to hash password")
		return
	}

	adminRequest.Password = string(passwordHash)

	err = u.adminRepository.UpdateAdmin(adminRequest)
	if err != nil {
		echo.NewHTTPError(400, "Failed to update Admin")
		return
	}

	res = dtos.UpdateAdminResponse{
		Nama:     adminRequest.Nama,
		Email:    adminRequest.Email,
		Password: adminRequest.Password,
		Role:     adminRequest.Role,
	}

	return
}

// Logic for Create Admin
func (u *adminUsecase) CreateAdmin(req *dtos.RegisterAdminRequest) error {
	adminRequest := &models.Administrator{
		Nama:     req.Nama,
		Email:    req.Email,
		Password: req.Password,
	}

	// Check apakah email sudah terdaftar atau belum
	_, err := u.adminRepository.GetAdminByEmail(adminRequest.Email)
	if err == nil {
		return echo.NewHTTPError(400, "Email sudah terdaftar")
	}

	passwordHash, err := helpers.HashPassword(adminRequest.Password)
	if err != nil {
		return echo.NewHTTPError(400, "Failed to hash password")
	}

	adminRequest.Password = string(passwordHash)

	err = u.adminRepository.CreateAdmin(adminRequest)

	if err != nil {
		return errors.New(err.Error())
	}

	return nil
}

// Logic for Delete Administrator
func (u *adminUsecase) DeleteAdmin(id int, req *dtos.DeleteAdminRequest) (res helpers.ResponseMessage, err error) {
	admin, err := u.adminRepository.ReadToken(id)

	if err != nil {
		echo.NewHTTPError(400, "Failed to get Admin")
		return
	}

	err = helpers.ComparePassword(req.Password, admin.Password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		echo.NewHTTPError(400, err.Error())
		return
	}

	err = u.adminRepository.DeleteAdmin(admin)

	if err != nil {
		return res, echo.NewHTTPError(500, "Failed to delete admin")
	}

	return res, nil
}

// Logic for Login Administrator
func (u *adminUsecase) LoginAdmin(c echo.Context, req *dtos.LoginRequest) (res dtos.LoginResponse, err error) {
	admin, err := u.adminRepository.GetAdminByEmail(req.Email)
	if err != nil {
		echo.NewHTTPError(400, "Email not registered")
		return
	}

	err = helpers.ComparePassword(req.Password, admin.Password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		echo.NewHTTPError(400, err.Error())
		return
	}

	token, err := middlewares.CreateToken(int(admin.ID), admin.Email, admin.Role)
	if err != nil {
		echo.NewHTTPError(400, "Failed to generate token")
		return
	}

	admin.Token = token

	middlewares.CreateCookie(c, token)

	res = dtos.LoginResponse{
		Email: admin.Email,
		Token: admin.Token,
	}

	return
}
