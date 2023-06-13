package controllers

import (
	"go_bedu/dtos"
	"go_bedu/helpers"
	"go_bedu/usecase"
	"net/http"

	"github.com/labstack/echo/v4"
)

type AuthControllers interface {
	ForgotPasswordControllers(c echo.Context) error
}

type authControllers struct {
	authUsecase usecase.AuthUsecase
}

func NewAuthControllers(authUsecase usecase.AuthUsecase) AuthControllers {
	return &authControllers{
		authUsecase: authUsecase,
	}
}

func (c *authControllers) ForgotPasswordControllers(ctx echo.Context) error {
	req := dtos.ForgotPasswordRequest{}
	ctx.Bind(&req)
	if err := ctx.Validate(&req); err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Field cannot be empty or Password must be 6 characters",
				helpers.GetErrorData(err),
			),
		)
	}

	res, err := c.authUsecase.ForgotPassword(req)
	if err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Could not forgot password",
				helpers.GetErrorData(err),
			),
		)
	}

	return ctx.JSON(
		http.StatusOK,
		helpers.NewResponse(
			http.StatusOK,
			"Success Reset Password",
			res,
		),
	)
}
