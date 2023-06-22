package controllers

import (
	"fmt"
	"go_bedu/dtos"
	"go_bedu/helpers"
	m "go_bedu/middlewares"
	"go_bedu/repositories"
	"go_bedu/usecase"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type AdminController interface {
	LoginAdminController(c echo.Context) error
	LogoutAdminController(c echo.Context) error
	RegisterAdminController(c echo.Context) error
	VerifyEmailAdminController(c echo.Context) error
	VerifyOTPAdminController(c echo.Context) error
	ChangePasswordController(c echo.Context) error
	GetAdminsController(c echo.Context) error
	GetAdminByIdController(c echo.Context) error
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

// Controller for Logout Admin from Cookie
func (c *adminController) LogoutAdminController(ctx echo.Context) error {
	_, err := m.IsAdmin(ctx)
	if err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Please login first",
				helpers.GetErrorData(err),
			),
		)
	}

	_, err = c.adminUsecase.LogoutAdmin(ctx)
	if err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Logout failed",
				helpers.GetErrorData(err),
			),
		)
	}

	return ctx.JSON(
		http.StatusOK,
		helpers.NewResponseMessage(
			http.StatusOK,
			"Success Logout",
		),
	)
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
	ctx.Response().Header().Set("Content-Type", "application/json")
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

	message := "We sent an email with a verification code to " + admin.Email
	return ctx.JSON(http.StatusOK,
		helpers.NewResponse(
			http.StatusOK,
			"Success Create Account",
			message,
		))
}

func (c *adminController) VerifyEmailAdminController(ctx echo.Context) error {
	code := ctx.Param("verificationCode")
	// verification_code := utils.Encode(code)

	res, err := c.adminUsecase.VerifyEmail(code)
	if err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Could not verify email",
				helpers.GetErrorData(err),
			),
		)
	}

	return ctx.JSON(
		http.StatusOK,
		helpers.NewResponse(
			http.StatusOK,
			"Success Verify Email",
			res,
		),
	)
}

// Controller for verify OTP account
func (c *adminController) VerifyOTPAdminController(ctx echo.Context) error {
	req := dtos.ChangePasswordRequest{}

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

	code, err := strconv.Atoi(ctx.Param("otp"))
	if err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Could not verify OTP",
				helpers.GetErrorData(err),
			),
		)
	}

	res, err := c.adminUsecase.UpdateAdminByOTP(code, req)
	if err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"OTP is not valid",
				helpers.GetErrorData(err),
			),
		)
	}

	return ctx.JSON(
		http.StatusOK,
		helpers.NewResponse(
			http.StatusOK,
			"Success Forgot Password",
			res,
		),
	)
}

// Controller for Change Password Admin
func (c *adminController) ChangePasswordController(ctx echo.Context) error {
	req := dtos.ChangePasswordAdminRequest{}

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

	res, err := c.adminUsecase.ChangePassword(uint(id), req)
	if err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Could not change password",
				helpers.GetErrorData(err),
			),
		)
	}

	return ctx.JSON(
		http.StatusOK,
		helpers.NewResponse(
			http.StatusOK,
			"Success Change Password",
			res,
		),
	)
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
				"Field cannot be empty",
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
