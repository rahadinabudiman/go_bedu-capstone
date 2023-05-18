package controllers

import (
	"go_bedu/lib/database"
	"go_bedu/models"
	"net/http"

	"github.com/labstack/echo"
)

func GetAdministratorController(c echo.Context) error {
	admins, err := database.GetAdministrators()

	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Message: err.Error(),
		})
	}

	if len(admins) == 0 {
		return c.JSON(http.StatusOK, models.Response{
			Message: "No administrator found",
		})
	}

	alladminres := make([]models.AdminsResponse, len(admins))
	for i, admin := range admins {
		alladminres[i] = models.AdminsResponse{
			Nama:  admin.Nama,
			Email: admin.Email,
		}
	}

	return c.JSON(http.StatusOK, models.Response{
		Message: "success get all admin",
		Data:    alladminres,
	})
}
