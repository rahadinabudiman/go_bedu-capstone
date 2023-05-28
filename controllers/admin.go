package controllers

import (
	"fmt"
	"go_bedu/dtos"
	"go_bedu/helpers"
	m "go_bedu/middlewares"
	"go_bedu/repositories"
	"go_bedu/usecase"
	"net/http"

	"github.com/labstack/echo/v4"
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
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Field cannot be empty or Password must be 6 character",
				helpers.GetErrorData(err),
			),
		)
	}

	admin, err := c.adminUsecase.CreateAdmin(&req)
	if err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Could not create admin",
				helpers.GetErrorData(err),
			),
		)
	}

	return ctx.JSON(http.StatusOK,
		helpers.NewResponse(
			http.StatusOK,
			"Success Create Admin",
			admin,
		))
}

// Controller for Get All Admins from DB
func (c *adminController) GetAdminsController(ctx echo.Context) error {
	admins, err := c.adminUsecase.GetAdmin()
	if err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Could not get admin",
				helpers.GetErrorData(err),
			),
		)
	}

	return ctx.JSON(
		http.StatusOK,
		helpers.NewResponse(
			http.StatusOK,
			"Success Get Admin",
			admins,
		),
	)
}

// Controller for Get Admin by ID from DB
func (c *adminController) GetAdminByIdController(ctx echo.Context) error {
	id, err := m.IsAdmin(ctx)
	if err != nil {
		return ctx.JSON(
			http.StatusUnauthorized,
			helpers.NewErrorResponse(
				http.StatusUnauthorized,
				"Routes for Admin Only",
				helpers.GetErrorData(err),
			),
		)
	}

	res, err := c.adminUsecase.GetAdminById(uint(id))
	if err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Could not get admin",
				helpers.GetErrorData(err),
			),
		)
	}

	return ctx.JSON(http.StatusOK, helpers.Response{
		Message: fmt.Sprintf("Welcome %s", res.Nama),
		Data:    res,
	})
}

// Controller for Update Admin by ID from DB
func (c *adminController) UpdateAdminController(ctx echo.Context) error {
	req := dtos.UpdateAdminRequest{}

	id, err := m.IsAdmin(ctx)
	if err != nil {
		return ctx.JSON(
			http.StatusUnauthorized,
			helpers.NewErrorResponse(
				http.StatusUnauthorized,
				"Routes for Admin Only",
				helpers.GetErrorData(err),
			),
		)
	}

	ctx.Bind(&req)
	if err := ctx.Validate(&req); err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Field cannot be empty or Password must be 6 character",
				helpers.GetErrorData(err),
			),
		)
	}

	res, err := c.adminUsecase.UpdateAdmin(uint(id), req)
	if err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Could not update admin",
				helpers.GetErrorData(err),
			),
		)
	}

	return ctx.JSON(
		http.StatusOK,
		helpers.NewResponse(
			http.StatusOK,
			"Success Update Admin",
			res,
		),
	)
}

// Controller for Delete Admin by ID from DB
func (c *adminController) DeleteAdminController(ctx echo.Context) error {
	id, err := m.IsAdmin(ctx)
	if err != nil {
		return ctx.JSON(
			http.StatusUnauthorized,
			helpers.NewErrorResponse(
				http.StatusUnauthorized,
				"Routes for Admin Only",
				helpers.GetErrorData(err),
			),
		)
	}

	req := dtos.DeleteAdminRequest{}

	ctx.Bind(&req)
	if err := ctx.Validate(&req); err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Field cannot be empty or Password must be 6 character",
				helpers.GetErrorData(err),
			),
		)
	}

	_, err = c.adminUsecase.DeleteAdmin(uint(id), req)
	if err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Could not Delete admin",
				helpers.GetErrorData(err),
			),
		)
	}

	return ctx.JSON(
		http.StatusOK,
		helpers.NewResponseMessage(
			http.StatusOK,
			"Success Delete Admin",
		),
	)
}
