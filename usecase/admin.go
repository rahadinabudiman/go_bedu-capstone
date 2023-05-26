package usecase

import (
	"go_bedu/dtos"
	"go_bedu/helpers"
	"go_bedu/middlewares"
	"go_bedu/models"
	"go_bedu/repositories"

	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"
)

type AdminUsecase interface {
	LoginAdmin(c echo.Context, req dtos.LoginRequest) (res dtos.LoginResponse, err error)
	GetAdmin() ([]dtos.AdminDetailResponse, error)
	GetAdminById(id uint) (res dtos.AdminProfileResponse, err error)
	UpdateAdmin(id uint, req dtos.UpdateAdminRequest) (res dtos.UpdateAdminResponse, err error)
	CreateAdmin(req *dtos.RegisterAdminRequest) (dtos.AdminDetailResponse, error)
	DeleteAdmin(id uint, req dtos.DeleteAdminRequest) (res helpers.ResponseMessage, err error)
}

type adminUsecase struct {
	adminRepository repositories.AdminRepository
}

func NewAdminUsecase(adminRepository repositories.AdminRepository) *adminUsecase {
	return &adminUsecase{adminRepository}
}

// Logic for get All Admin
func (u *adminUsecase) GetAdmin() ([]dtos.AdminDetailResponse, error) {
	admins, err := u.adminRepository.GetAdmins()
	if err != nil {
		return nil, err
	}

	var adminResponse []dtos.AdminDetailResponse
	for _, admin := range admins {
		adminResponse = append(adminResponse, dtos.AdminDetailResponse{
			AdministratorID: admin.ID,
			Nama:            admin.Nama,
			Email:           admin.Email,
			Role:            admin.Role,
			CreatedAt:       admin.CreatedAt,
			UpdatedAt:       admin.UpdatedAt,
		})
	}

	return adminResponse, nil
}

// Logic for get Admin with Cookie
func (u *adminUsecase) GetAdminById(id uint) (res dtos.AdminProfileResponse, err error) {
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
func (u *adminUsecase) UpdateAdmin(id uint, req dtos.UpdateAdminRequest) (dtos.UpdateAdminResponse, error) {
	var (
		admins models.Administrator
		res    dtos.UpdateAdminResponse
	)

	admins, err := u.adminRepository.GetAdminById(id)
	if err != nil {
		return res, err
	}

	admins.Nama = req.Nama
	admins.Email = req.Email
	admins.Password = req.Password
	admins.Role = req.Role

	// Check Role and save role information from JWT Cookie
	admin, err := u.adminRepository.ReadToken(id)
	if err != nil {
		return res, echo.NewHTTPError(400, "Failed to get Admin")
	}

	// Check if admin is not Super Admin and trying to change role
	if admin.Role != req.Role && admin.Role != "Super Admin" {
		return res, echo.NewHTTPError(401, "You are not allowed to change role")
	}

	// Check if Super Admin changes other Role
	if admin.Role == "Super Admin" {
		if req.Role != "Super Admin" && req.Role != "Admin" {
			return res, echo.NewHTTPError(401, "Tidak bisa mengubah role menjadi selain Super Admin atau Admin")
		}
	}

	admins.ID = uint(id)

	passwordHash, err := helpers.HashPassword(admins.Password)
	if err != nil {
		return res, echo.NewHTTPError(400, "Failed to hash password")
	}

	admins.Password = string(passwordHash)

	admins, err = u.adminRepository.UpdateAdmin(admins)
	if err != nil {
		return res, err
	}

	res.Nama = admins.Nama
	res.Email = admins.Email
	res.Password = admins.Password
	res.Role = admins.Role

	return res, nil
}

// Logic for Create Admin
func (u *adminUsecase) CreateAdmin(req *dtos.RegisterAdminRequest) (dtos.AdminDetailResponse, error) {
	var res dtos.AdminDetailResponse

	CreateAdmin := models.Administrator{
		Nama:     req.Nama,
		Email:    req.Email,
		Password: req.Password,
	}

	// Check apakah email sudah terdaftar atau belum
	_, err := u.adminRepository.GetAdminByEmail(req.Email)
	if err == nil {
		return res, echo.NewHTTPError(400, "Email sudah terdaftar")
	}

	passwordHash, err := helpers.HashPassword(req.Password)
	if err != nil {
		return res, err
	}

	req.Password = string(passwordHash)

	admins, err := u.adminRepository.CreateAdmin(CreateAdmin)

	if err != nil {
		return res, err
	}

	resp := dtos.AdminDetailResponse{
		AdministratorID: admins.ID,
		Nama:            admins.Nama,
		Email:           admins.Email,
		Role:            admins.Role,
		CreatedAt:       admins.CreatedAt,
		UpdatedAt:       admins.UpdatedAt,
	}

	return resp, nil
}

// Logic for Delete Administrator
func (u *adminUsecase) DeleteAdmin(id uint, req dtos.DeleteAdminRequest) (res helpers.ResponseMessage, err error) {
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
func (u *adminUsecase) LoginAdmin(c echo.Context, req dtos.LoginRequest) (res dtos.LoginResponse, err error) {
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
