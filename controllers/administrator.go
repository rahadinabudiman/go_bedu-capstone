package controllers

import (
	"go_bedu/lib/database"
	"go_bedu/models"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

// Get All Admins from DB
func GetAdministratorController(c echo.Context) error {
	admins, err := database.GetAdministrators()

	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ResponseMessage{
			Message: err.Error(),
		})
	}

	if len(admins) == 0 {
		return c.JSON(http.StatusOK, models.ResponseMessage{
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

// Create Admin to DB
func CreateAdministratorController(c echo.Context) error {
	admin := models.Administrator{}
	c.Bind(&admin)

	if err := c.Validate(&admin); err != nil {
		return c.JSON(http.StatusBadRequest, models.ResponseMessage{
			Message: err.Error(),
		})
	}

	admins, err := database.CreateAdministrator(admin)

	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ResponseMessage{
			Message: err.Error(),
		})
	}

	AdminCreateRES := models.AdminCreateRES{
		Nama:     admins.Nama,
		Email:    admins.Email,
		Password: admins.Password,
	}

	return c.JSON(http.StatusOK, models.Response{
		Message: "success create admin",
		Data:    AdminCreateRES,
	})
}

// Get Administrator by ID From DB
func GetAdministratorByIDController(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ResponseMessage{
			Message: err.Error(),
		})
	}

	admin, err := database.GetAdministratorByID(id)

	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ResponseMessage{
			Message: err.Error(),
		})
	}

	AdminsResponse := models.AdminsResponse{
		Nama:  admin.Nama,
		Email: admin.Email,
	}

	return c.JSON(http.StatusOK, models.Response{
		Message: "success get admin",
		Data:    AdminsResponse,
	})
}

// Updates Administrator by ID from DB
func UpdateAdministratorByIdController(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ResponseMessage{
			Message: err.Error(),
		})
	}

	admin := models.Administrator{}
	c.Bind(&admin)

	if err := c.Validate(&admin); err != nil {
		return c.JSON(http.StatusBadRequest, models.ResponseMessage{
			Message: err.Error(),
		})
	}

	admin, err = database.UpdateAdministrator(admin, id)

	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ResponseMessage{
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, models.Response{
		Message: "success update administrator",
		Data:    admin,
	})
}

// Delete Administrator by ID from DB
func DeleteAdministratorController(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ResponseMessage{
			Message: err.Error(),
		})
	}

	// Check Apakah ID ada di DB
	_, err = database.GetAdministratorByID(id)

	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ResponseMessage{
			Message: err.Error(),
		})
	}

	_, err = database.DeleteAdministrator(id)

	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ResponseMessage{
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, models.ResponseMessage{
		Message: "success delete administrator",
	})
}
