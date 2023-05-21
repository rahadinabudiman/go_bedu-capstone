package routes

import (
	"go_bedu/constants"
	"go_bedu/controllers"
	m "go_bedu/middlewares"
	"go_bedu/repository/database"
	"go_bedu/usecase"
	"go_bedu/utils"

	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	mid "github.com/labstack/echo/middleware"
)

func NewRoute(e *echo.Echo, db *gorm.DB) {
	m.Log(e)
	e.Pre(mid.RemoveTrailingSlash())

	adminRepository := database.NewAdminRepository(db)
	adminUsecase := usecase.NewAdminUsecase(adminRepository)
	adminController := controllers.NewAdminController(adminUsecase, adminRepository)

	authRepository := database.NewAuthRepository(db)
	authUsecase := usecase.NewAuthUsecase(authRepository, adminRepository)
	authController := controllers.NewAuthController(authUsecase, authRepository, adminUsecase)

	// Middleware untuk mengatur CORS
	e.Use(mid.CORSWithConfig(mid.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization, // Menambahkan header Authorization
		},
	}))

	e.Use(m.AllowCORS)

	cv := &utils.CustomValidator{Validators: validator.New()}
	e.Validator = cv

	e.GET("/user", controllers.UserHandler, mid.JWTWithConfig(mid.JWTConfig{
		SigningMethod: "HS256",
		SigningKey:    []byte(constants.SECRET_JWT),
		TokenLookup:   "header:Authorization",
	}))

	e.POST("/register", authController.RegisterAdminController)
	e.POST("/login", authController.LoginAdminController)

	// Routes Baru
	admin := e.Group("/admin")
	admin.GET("", adminController.GetAdminsController, m.VerifyToken)
	admin.GET("/profile", adminController.GetAdminByIdController, m.VerifyToken)
	admin.PUT("", adminController.UpdateAdminController, m.VerifyToken)
	admin.DELETE("", adminController.DeleteAdminController, m.VerifyToken)

	// Article Routes
	article := e.Group("/article")
	article.GET("", controllers.GetArticlesControllers)
	article.POST("", controllers.CreateArticleController)
	article.GET("/:id", controllers.GetArticleByIDController)
	article.PUT("/:id", controllers.UpdateArticleByIdController)
	article.DELETE("/:id", controllers.DeleteArticleController)
}
