package controllers

import (
	"fmt"
	"go_bedu/dtos"
	"go_bedu/helpers"
	m "go_bedu/middlewares"
	"go_bedu/repositories"
	"go_bedu/usecase"
	"net/http"

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
func (c *adminController) LoginAdminController(ctx echo.Context) error {
	req := dtos.LoginRequest{}

	ctx.Bind(&req)
	if err := ctx.Validate(&req); err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Filed Cannot Be Empty",
				helpers.GetErrorData(err),
			),
		)
	}

	res, err := c.adminUsecase.LoginAdmin(ctx, req)
	if err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Could not login",
				helpers.GetErrorData(err),
			),
		)
	}

	return ctx.JSON(
		http.StatusOK,
		helpers.NewResponse(
			http.StatusOK,
			"Success Login",
			res,
		),
	)
}

// Controller For Register Admin From DB
func (c *adminController) RegisterAdminController(ctx echo.Context) error {
	req := dtos.RegisterAdminRequest{}

	ctx.Bind(&req)
	if err := ctx.Validate(&req); err != nil {
		return echo.NewHTTPError(400, "Field cannot be empty or Password must be 6 character")
	}

	admin, err := c.adminUsecase.CreateAdmin(&req)

	if err != nil {
		return echo.NewHTTPError(400, err.Error())
	}

	return ctx.JSON(200, map[string]interface{}{
		"message": "Success Register",
		"data":    admin,
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
func (c *adminController) GetAdminByIdController(ctx echo.Context) error {
	id, err := m.IsAdmin(ctx)
	if err != nil {
		return echo.NewHTTPError(401, "This routes for admin only")
	}

	res, err := c.adminUsecase.GetAdminById(uint(id))
	if err != nil {
		return echo.NewHTTPError(400, err.Error())
	}

	return ctx.JSON(200, helpers.Response{
		Message: fmt.Sprintf("Welcome %s", res.Nama),
		Data:    res,
	})
}

// Controller for Update Admin by ID from DB
func (c *adminController) UpdateAdminController(ctx echo.Context) error {
	req := dtos.UpdateAdminRequest{}

	id, err := m.IsAdmin(ctx)
	if err != nil {
		return echo.NewHTTPError(401, "This routes for admin only")
	}

	ctx.Bind(&req)

	if err := ctx.Validate(&req); err != nil {
		return echo.NewHTTPError(400, "Field cannot be empty")
	}

	res, err := c.adminUsecase.UpdateAdmin(uint(id), req)

	if err != nil {
		return echo.NewHTTPError(400, err.Error())
	}

	return ctx.JSON(200, helpers.Response{
		Message: "Success update admin",
		Data:    res,
	})
}

// Controller for Delete Admin by ID from DB
func (c *adminController) DeleteAdminController(ctx echo.Context) error {
	id, err := m.IsAdmin(ctx)
	if err != nil {
		return echo.NewHTTPError(401, "this routes for admin only")
	}

	req := dtos.DeleteAdminRequest{}

	ctx.Bind(&req)
	if err := ctx.Validate(&req); err != nil {
		return echo.NewHTTPError(400, "Field cannot be empty")
	}

	fmt.Printf(req.Password)

	_, err = c.adminUsecase.DeleteAdmin(uint(id), req)
	if err != nil {
		return echo.NewHTTPError(400, err.Error())
	}

	return ctx.JSON(200, helpers.ResponseMessage{
		Message: "Delete Admin Sukses",
	})
}
