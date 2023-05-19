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
	administrator.POST("", controllers.CreateAdministratorController)
	administrator.GET("/:id", controllers.GetAdministratorByIDController)
	administrator.PUT("/:id", controllers.UpdateAdministratorByIdController)
	administrator.DELETE("/:id", controllers.DeleteAdministratorController)

	// Article Routes
	article := e.Group("/article")
	article.GET("", controllers.GetArticlesControllers)
	article.POST("", controllers.CreateArticleController)
	article.GET("/:id", controllers.GetArticleByIDController)
	article.PUT("/:id", controllers.UpdateArticleByIdController)
	article.DELETE("/:id", controllers.DeleteArticleController)

	return e
}
