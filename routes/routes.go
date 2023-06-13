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
	m.Log(e)
	e.Pre(mid.RemoveTrailingSlash())

	adminRepository := repositories.NewAdminRepository(db)
	adminUsecase := usecase.NewAdminUsecase(adminRepository)
	adminController := controllers.NewAdminController(adminUsecase, adminRepository)

	articleRepository := repositories.NewArticleRepository(db)
	articleUsecase := usecase.NewArticleUsecase(articleRepository)
	articleController := controllers.NewArticleController(articleUsecase)

	userRepository := repositories.NewUserRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepository)
	userController := controllers.NewUserControllers(userUsecase, userRepository)

	authUsecase := usecase.NewAuthUsecase(adminRepository, userRepository)
	authControllers := controllers.NewAuthControllers(authUsecase)

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

	// Main API
	api := e.Group("/api/v1")

	// AUTH API
	api.POST("/admin/register", adminController.RegisterAdminController)
	api.POST("/admin/login", adminController.LoginAdminController)
	api.POST("/login", userController.LoginUserController)
	api.POST("/register", userController.RegisterUserController)

	// Utils API
	api.POST("/change-password/:otp", userController.VerifyOTPUserController)
	api.POST("/admin/change-password/:otp", adminController.VerifyOTPAdminController)
	api.GET("/verifyemail/:verificationCode", userController.VerifyEmailUserController)
	api.GET("/admin/verifyemail/:verificationCode", adminController.VerifyEmailAdminController)

	api.POST("/forgot-password", authControllers.ForgotPasswordControllers)

	article := api.Group("/article")
	article.GET("", articleController.GetAllArticles)
	article.GET("/:id", articleController.GetArticleById)

	// User Only
	user := api.Group("/user")
	user.Use(m.VerifyToken)
	user.GET("", userController.GetAllUserController)
	user.GET("/profile", userController.GetUserController)
	user.PUT("", userController.UpdateUserController)
	user.DELETE("", userController.DeleteUserController)
	user.POST("/change-password", userController.ChangePasswordController)
	user.GET("/logout", userController.LogoutUserController)

	// Admin Only
	admin := api.Group("/admin")
	admin.Use(m.VerifyToken)
	admin.GET("", adminController.GetAdminsController)
	admin.GET("/profile", adminController.GetAdminByIdController)
	admin.PUT("", adminController.UpdateAdminController)
	admin.DELETE("", adminController.DeleteAdminController)
	admin.POST("/change-password", adminController.ChangePasswordController)
	admin.GET("/logout", adminController.LogoutAdminController)

	// Article Admin Routes
	admin.POST("/article", articleController.CreateArticle)
	admin.PUT("/article/:id", articleController.UpdateArticle)
	admin.DELETE("/article/:id", articleController.DeleteArticle)
}
