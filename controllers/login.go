package controllers

import (
	"go_bedu/lib/database"
	"go_bedu/middlewares"
	"go_bedu/models"
	"net/http"

	"github.com/labstack/echo"
)

func LoginAdministratorController(c echo.Context) error {
	admin := models.Administrator{}
	c.Bind(&admin)

	if err := c.Validate(&admin); err != nil {
		return c.JSON(http.StatusBadRequest, models.ResponseMessage{
			Message: err.Error(),
		})
	}

	admin, err := database.LoginAdministrator(admin)

	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ResponseMessage{
			Message: "Email atau Password tidak valid",
		})
	}

	token, err := middlewares.CreateToken(int(admin.ID), admin.Email, admin.Role)

	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ResponseMessage{
			Message: err.Error(),
		})
	}

	middlewares.CreateCookie(c, token)

	responseLogin := models.AdminsJWTRES{
		ID:    admin.ID,
		Nama:  admin.Nama,
		Email: admin.Email,
		Token: token,
	}

	return c.JSON(http.StatusOK, models.Response{
		Message: "success login admin",
		Data:    responseLogin,
	})

}
