package controllers

import (
	"go_bedu/models/payload"
	"go_bedu/repository/database"
	"go_bedu/usecase"

	"github.com/labstack/echo"
)

type AuthController interface {
	LoginAdminController(c echo.Context) error
	RegisterAdminController(c echo.Context) error
}

type authController struct {
	authUsecase    usecase.AuthUsecase
	authRepository database.AuthRepository
	adminUsecase   usecase.AdminUsecase
}

func NewAuthController(
	authUsecase usecase.AuthUsecase,
	authRepository database.AuthRepository,
	adminUsecase usecase.AdminUsecase) *authController {
	return &authController{authUsecase, authRepository, adminUsecase}
}

// Controller for Login Admin from DB
func (a *authController) LoginAdminController(c echo.Context) error {
	req := payload.LoginRequest{}

	c.Bind(&req)
	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(400, "Field cannot be empty")
	}

	res, err := a.authUsecase.LoginUser(c, &req)
	if err != nil {
		return echo.NewHTTPError(400, "Invalid Email or Password")
	}

	return c.JSON(200, payload.Response{
		Message: "Success Login",
		Data:    res,
	})
}

// Controller For Register Admin From DB
func (a *authController) RegisterAdminController(c echo.Context) error {
	req := payload.RegisterAdminRequest{}

	c.Bind(&req)
	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(400, "Field cannot be empty or Password must be 6 character")
	}

	err := a.adminUsecase.CreateAdmin(&req)

	if err != nil {
		return echo.NewHTTPError(400, err.Error())
	}

	return c.JSON(200, map[string]interface{}{
		"message": "Success Register",
	})
}
