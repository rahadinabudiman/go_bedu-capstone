package middlewares

import (
	"go_bedu/constants"
	"go_bedu/models"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// Create Token JWT from constants
func CreateToken(id int, email, role string) (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["id"] = id
	claims["email"] = email
	claims["role"] = role
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(constants.SECRET_JWT))
}

// Check Login
var IsLoggedIn = middleware.JWTWithConfig(middleware.JWTConfig{
	SigningMethod: "HS256",
	SigningKey:    []byte(constants.SECRET_JWT),
	TokenLookup:   "cookie:bEDUCookie",
	ErrorHandler: func(err error) error {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	},
})

// Middleware untuk verifikasi JWT token
func VerifyToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Mendapatkan token dari header Authorization
		authHeader := c.Request().Header.Get("Authorization")
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Memeriksa apakah token ada
		if tokenString == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Missing token")
		}

		// Memverifikasi token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Menggunakan secret key yang sama dengan saat membuat token
			return []byte(constants.SECRET_JWT), nil
		})

		// Menangani error saat verifikasi token
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
		}

		// Memeriksa apakah token valid
		if !token.Valid {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
		}

		// Menyimpan informasi user dari token di dalam konteks
		claims := token.Claims.(jwt.MapClaims)
		userID := int(claims["id"].(float64))
		email := claims["email"].(string)
		role := claims["role"].(string)

		// Menyimpan informasi user di dalam konteks
		c.Set("userID", userID)
		c.Set("email", email)
		c.Set("role", role)

		// Melanjutkan ke handler selanjutnya
		return next(c)
	}
}

// Verifikasi Super Admin Middleware
func VerifySuperAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Mendapatkan token dari header Authorization
		authHeader := c.Request().Header.Get("Authorization")
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Memeriksa apakah token ada
		if tokenString == "" {
			return c.JSON(http.StatusUnauthorized, models.ResponseMessage{
				Message: "Missing token",
			})
		}

		// Memverifikasi token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Menggunakan secret key yang sama dengan saat membuat token
			return []byte(constants.SECRET_JWT), nil
		})

		// Menangani error saat verifikasi token
		if err != nil {
			return c.JSON(http.StatusUnauthorized, models.ResponseMessage{
				Message: "Invalid token",
			})
		}

		// Memeriksa apakah token valid
		if !token.Valid {
			return c.JSON(http.StatusUnauthorized, models.ResponseMessage{
				Message: "Invalid token",
			})
		}

		// Menyimpan informasi user dari token di dalam konteks
		claims := token.Claims.(jwt.MapClaims)
		userID := int(claims["id"].(float64))
		email := claims["email"].(string)
		role := claims["role"].(string)

		// Memeriksa apakah role adalah super admin
		if role != "Super Admin" {
			return c.JSON(http.StatusForbidden, models.ResponseMessage{
				Message: "Unauthorized access",
			})
		}

		// Menyimpan informasi user di dalam konteks
		c.Set("userID", userID)
		c.Set("email", email)
		c.Set("role", role)

		// Melanjutkan ke handler selanjutnya
		return next(c)
	}
}
