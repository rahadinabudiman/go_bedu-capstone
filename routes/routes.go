package routes

import (
	"go_bedu/controllers"
	m "go_bedu/middlewares"
	"go_bedu/repositories"
	"go_bedu/usecase"
	"go_bedu/utils"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	mid "github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"
)

func NewRoute(e *echo.Echo, db *gorm.DB) {
	adminRepository := repositories.NewAdminRepository(db)
	adminUsecase := usecase.NewAdminUsecase(adminRepository)
	adminController := controllers.NewAdminController(adminUsecase, adminRepository)

	articleRepository := repositories.NewArticleRepository(db)
	articleUsecase := usecase.NewArticleUsecase(articleRepository)
	articleController := controllers.NewArticleController(articleUsecase)

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

	// Mengatur folder untuk file gambar
	e.Static("/public", "public")

	api := e.Group("/api/v1")

	api.POST("/register", adminController.RegisterAdminController)
	api.POST("/login", adminController.LoginAdminController)

	// Admin Only
	admin := api.Group("/admin")
	admin.Use(m.VerifyToken)
	admin.GET("", adminController.GetAdminsController)
	admin.GET("/profile", adminController.GetAdminByIdController)
	admin.PUT("", adminController.UpdateAdminController)
	admin.DELETE("", adminController.DeleteAdminController)

	// Article Routes

	admin.GET("/article", articleController.GetAllArticles)
	admin.GET("/article/:id", articleController.GetArticleById)
	admin.POST("/article", articleController.CreateArticle)
	admin.PUT("/article/:id", articleController.UpdateArticle)
	admin.DELETE("/article/:id", articleController.DeleteArticle)
}
