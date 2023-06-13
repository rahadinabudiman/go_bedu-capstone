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

type UserControllers interface {
	LoginUserController(c echo.Context) error
	LogoutUserController(c echo.Context) error
	RegisterUserController(c echo.Context) error
	VerifyEmailUserController(c echo.Context) error
	VerifyOTPUserController(c echo.Context) error
	ChangePasswordController(c echo.Context) error
	GetAllUserController(c echo.Context) error
	GetUserController(c echo.Context) error
	UpdateUserController(c echo.Context) error
	DeleteUserController(c echo.Context) error
}

type userControllers struct {
	userUsecase    usecase.UserUsecase
	userRepository repositories.UserRepository
}

func NewUserControllers(userUsecase usecase.UserUsecase, userRepository repositories.UserRepository) UserControllers {
	return &userControllers{
		userUsecase:    userUsecase,
		userRepository: userRepository,
	}
}

// Controller for Login User
func (c *userControllers) LoginUserController(ctx echo.Context) error {
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

	res, err := c.userUsecase.LoginUser(ctx, req)
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

// Controller for Logout User from Cookie
func (c *userControllers) LogoutUserController(ctx echo.Context) error {
	_, err := m.IsUser(ctx)
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

	_, err = c.userUsecase.LogoutUser(ctx)
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

func (c *userControllers) RegisterUserController(ctx echo.Context) error {
	req := dtos.RegisterUserRequest{}

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

	user, err := c.userUsecase.CreateUser(&req)
	if err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Could not create user",
				helpers.GetErrorData(err),
			),
		)
	}

	message := "We sent an email with a verification code to " + user.Email
	return ctx.JSON(http.StatusOK,
		helpers.NewResponse(
			http.StatusOK,
			"Success Create Account",
			message,
		))
}

func (c *userControllers) VerifyEmailUserController(ctx echo.Context) error {
	code := ctx.Param("verificationCode")
	// verification_code := utils.Encode(code)

	res, err := c.userUsecase.VerifyEmail(code)
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

func (c *userControllers) VerifyOTPUserController(ctx echo.Context) error {
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

	res, err := c.userUsecase.UpdateUserByOTP(code, req)
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

func (c *userControllers) ChangePasswordController(ctx echo.Context) error {
	req := dtos.ChangePasswordUserRequest{}

	id, err := m.IsUser(ctx)
	if err != nil {
		return ctx.JSON(
			http.StatusUnauthorized,
			helpers.NewErrorResponse(
				http.StatusUnauthorized,
				"Routes for User Only",
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

	res, err := c.userUsecase.ChangePassword(uint(id), req)
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

func (c *userControllers) GetAllUserController(ctx echo.Context) error {
	users, err := c.userUsecase.GetUsers()
	if err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Could not get users",
				helpers.GetErrorData(err),
			),
		)
	}

	return ctx.JSON(
		http.StatusOK,
		helpers.NewResponse(
			http.StatusOK,
			"Success Get Users",
			users,
		),
	)
}

func (c *userControllers) GetUserController(ctx echo.Context) error {
	id, err := m.IsUser(ctx)
	if err != nil {
		return ctx.JSON(
			http.StatusUnauthorized,
			helpers.NewErrorResponse(
				http.StatusUnauthorized,
				"Routes for User Only",
				helpers.GetErrorData(err),
			),
		)
	}

	res, err := c.userUsecase.GetUserById(uint(id))
	if err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Could not get user",
				helpers.GetErrorData(err),
			),
		)
	}

	return ctx.JSON(http.StatusOK, helpers.Response{
		Message: fmt.Sprintf("Welcome %s", res.Nama),
		Data:    res,
	})
}

func (c *userControllers) UpdateUserController(ctx echo.Context) error {
	req := dtos.UpdateUserRequest{}

	id, err := m.IsUser(ctx)
	if err != nil {
		return ctx.JSON(
			http.StatusUnauthorized,
			helpers.NewErrorResponse(
				http.StatusUnauthorized,
				"Routes for User Only",
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

	res, err := c.userUsecase.UpdateUser(uint(id), req)
	if err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Could not update user",
				helpers.GetErrorData(err),
			),
		)
	}

	return ctx.JSON(
		http.StatusOK,
		helpers.NewResponse(
			http.StatusOK,
			"Success Update User",
			res,
		),
	)
}

func (c *userControllers) DeleteUserController(ctx echo.Context) error {
	id, err := m.IsUser(ctx)
	if err != nil {
		return ctx.JSON(
			http.StatusUnauthorized,
			helpers.NewErrorResponse(
				http.StatusUnauthorized,
				"Routes for User Only",
				helpers.GetErrorData(err),
			),
		)
	}

	req := dtos.DeleteUserRequest{}

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

	_, err = c.userUsecase.DeleteUser(uint(id), req)
	if err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Could not Delete user",
				helpers.GetErrorData(err),
			),
		)
	}

	return ctx.JSON(
		http.StatusOK,
		helpers.NewResponseMessage(
			http.StatusOK,
			"Success Delete User",
		),
	)
}
