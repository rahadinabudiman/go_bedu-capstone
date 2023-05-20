package controllers

import "github.com/labstack/echo"

type AdminController interface {
	GetAdminsController(c echo.Context) error
	CreateAdminController(c echo.Context) error
}
