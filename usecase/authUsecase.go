package usecase

import (
	"go_bedu/helpers"
	"go_bedu/middlewares"
	"go_bedu/models/payload"
	"go_bedu/repositories"

	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase interface {
	LoginUser(c echo.Context, req *payload.LoginRequest) (res payload.LoginResponse, err error)
}

type authUsecase struct {
	authRepository  repositories.AuthRepository
	adminRepository repositories.AdminRepository
}

func NewAuthUsecase(authRepository repositories.AuthRepository, adminRepository repositories.AdminRepository) *authUsecase {
	return &authUsecase{authRepository, adminRepository}
}

// Logic for Login Administrator
func (a *authUsecase) LoginUser(c echo.Context, req *payload.LoginRequest) (res payload.LoginResponse, err error) {
	admin, err := a.adminRepository.GetAdminByEmail(req.Email)
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

	res = payload.LoginResponse{
		Email: admin.Email,
		Token: admin.Token,
	}

	return
}
