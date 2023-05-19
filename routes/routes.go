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

	// Login Routes
	login := e.Group("/login")
	login.POST("", controllers.LoginAdministratorController)

	// Administrator Routes
	administrator := e.Group("/administrator")
	administrator.GET("", controllers.GetAdministratorController, m.IsLoggedIn, m.IsSuperAdmin)
	administrator.POST("", controllers.CreateAdministratorController)
	administrator.PUT("/:id", controllers.UpdateAdministratorByIdController)
	administrator.GET("/:id", controllers.GetAdministratorByIDController, m.IsLoggedIn, m.IsSuperAdmin)
	administrator.DELETE("/:id", controllers.DeleteAdministratorController, m.IsLoggedIn, m.IsSuperAdmin)

	// Article Routes
	article := e.Group("/article")
	article.GET("", controllers.GetArticlesControllers)
	article.POST("", controllers.CreateArticleController)
	article.GET("/:id", controllers.GetArticleByIDController)
	article.PUT("/:id", controllers.UpdateArticleByIdController)
	article.DELETE("/:id", controllers.DeleteArticleController)

	return e
}
