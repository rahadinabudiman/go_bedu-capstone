package controllers

import (
	"fmt"
	"go_bedu/dtos"
	"go_bedu/helpers"
	m "go_bedu/middlewares"
	"go_bedu/repositories"
	"go_bedu/usecase"

	"github.com/labstack/echo"
)

type AdminController interface {
	LoginAdminController(c echo.Context) error
	RegisterAdminController(c echo.Context) error
	GetAdminsController(c echo.Context) error
	GetAdminByIdController(c echo.Context) error
	CreateAdminController(c echo.Context) error
	UpdateAdminController(c echo.Context) error
	DeleteAdminController(c echo.Context) error
}

type adminController struct {
	adminUsecase    usecase.AdminUsecase
	adminRepository repositories.AdminRepository
}

func NewAdminController(adminUsecase usecase.AdminUsecase, adminRepository repositories.AdminRepository) *adminController {
	return &adminController{adminUsecase, adminRepository}
}

// Controller for Login Admin from DB
func (a *adminController) LoginAdminController(c echo.Context) error {
	req := dtos.LoginRequest{}

	c.Bind(&req)
	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(400, "Field cannot be empty")
	}

	res, err := a.adminUsecase.LoginAdmin(c, &req)
	if err != nil {
		return echo.NewHTTPError(400, "Invalid Email or Password")
	}

	return c.JSON(200, helpers.Response{
		Message: "Success Login",
		Data:    res,
	})
}

// Controller For Register Admin From DB
func (a *adminController) RegisterAdminController(c echo.Context) error {
	req := dtos.RegisterAdminRequest{}

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

// Controller for Get All Admins from DB
func (a *adminController) GetAdminsController(c echo.Context) error {
	admins, err := a.adminUsecase.GetAdmin()

	if err != nil {
		return c.JSON(500, map[string]interface{}{
			"message": "Internal Server Error",
		})
	}

	return c.JSON(200, helpers.Response{
		Message: "Success Get All Admins",
		Data:    admins,
	})
}

// Controller for Get Admin by ID from DB
func (a *adminController) GetAdminByIdController(c echo.Context) error {
	id, err := m.IsAdmin(c)
	if err != nil {
		return echo.NewHTTPError(401, "This routes for admin only")
	}

	res, err := a.adminUsecase.GetAdminById(id)
	if err != nil {
		return echo.NewHTTPError(400, err.Error())
	}

	return c.JSON(200, helpers.Response{
		Message: fmt.Sprintf("Welcome %s", res.Nama),
		Data:    res,
	})
}

// Controller for Update Admin by ID from DB
func (a *adminController) UpdateAdminController(c echo.Context) error {
	req := dtos.UpdateAdminRequest{}

	id, err := m.IsAdmin(c)
	if err != nil {
		return echo.NewHTTPError(401, "This routes for admin only")
	}

	c.Bind(&req)

	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(400, "Field cannot be empty")
	}

	res, err := a.adminUsecase.UpdateAdmin(id, &req)

	if err != nil {
		return echo.NewHTTPError(400, err.Error())
	}

	return c.JSON(200, helpers.Response{
		Message: "Success update admin",
		Data:    res,
	})
}

// Controller for Delete Admin by ID from DB
func (a *adminController) DeleteAdminController(c echo.Context) error {
	id, err := m.IsAdmin(c)
	if err != nil {
		return echo.NewHTTPError(401, "this routes for admin only")
	}

	req := dtos.DeleteAdminRequest{}

	c.Bind(&req)
	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(400, "Field cannot be empty")
	}

	fmt.Printf(req.Password)

	_, err = a.adminUsecase.DeleteAdmin(id, &req)
	if err != nil {
		return echo.NewHTTPError(400, err.Error())
	}

	return c.JSON(200, helpers.ResponseMessage{
		Message: "Delete Admin Sukses",
	})
}
