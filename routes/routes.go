package routes

import (
	"go_bedu/controllers"
	m "go_bedu/middlewares"
	"go_bedu/utils"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo"
)

func New() *echo.Echo {
	e := echo.New()

	m.Log(e)
	cv := &utils.CustomValidator{Validators: validator.New()}
	e.Validator = cv

	// Administrator Routes
	administrator := e.Group("/administrator")
	administrator.GET("", controllers.GetAdministratorController)
	return e
}
