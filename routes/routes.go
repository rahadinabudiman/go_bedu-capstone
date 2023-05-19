package routes

import (
	"go_bedu/constants"
	"go_bedu/controllers"
	m "go_bedu/middlewares"
	"go_bedu/utils"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func New() *echo.Echo {
	e := echo.New()
	// Middleware untuk mengatur CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization, // Menambahkan header Authorization
		},
	}))

	e.Use(m.AllowCORS)

	m.Log(e)
	cv := &utils.CustomValidator{Validators: validator.New()}
	e.Validator = cv

	e.GET("/user", controllers.UserHandler, middleware.JWTWithConfig(middleware.JWTConfig{
		SigningMethod: "HS256",
		SigningKey:    []byte(constants.SECRET_JWT),
		TokenLookup:   "header:Authorization",
	}))

	// Login Routes
	login := e.Group("/login")
	login.POST("", controllers.LoginAdministratorController)

	// Administrator Routes
	administrator := e.Group("/administrator")
	administrator.GET("", controllers.GetAdministratorController, m.VerifyToken)
	administrator.POST("", controllers.CreateAdministratorController)
	administrator.PUT("/:id", controllers.UpdateAdministratorByIdController, m.VerifyToken)
	administrator.GET("/:id", controllers.GetAdministratorByIDController, m.VerifyToken, m.VerifySuperAdmin)
	administrator.DELETE("/:id", controllers.DeleteAdministratorController, m.VerifyToken, m.VerifySuperAdmin)

	// Article Routes
	article := e.Group("/article")
	article.GET("", controllers.GetArticlesControllers)
	article.POST("", controllers.CreateArticleController)
	article.GET("/:id", controllers.GetArticleByIDController)
	article.PUT("/:id", controllers.UpdateArticleByIdController)
	article.DELETE("/:id", controllers.DeleteArticleController)

	return e
}
