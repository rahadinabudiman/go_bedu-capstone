package middlewares

import (
	"go_bedu/constants"
	"go_bedu/models"
	"net/http"
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
},
)

// Check if Role is Super Admin
func IsSuperAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, ok := c.Get("user").(*jwt.Token)
		if !ok {
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid or missing jwt token")
		}
		claims, ok := user.Claims.(jwt.MapClaims)
		if !ok {
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid jwt claims")
		}
		if role, ok := claims["role"].(string); !ok || role != "Super Admin" {
			return echo.NewHTTPError(http.StatusUnauthorized, "user is not an Super Admin")
		}
		return next(c)
	}
}

// JWT Validator if Cookie is Login
func JWTValidator(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("bEDUCookie")
		if err != nil {
			return c.JSON(http.StatusUnauthorized, models.ResponseMessage{
				Message: "Unauthorized",
			})
		}
		token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
			return []byte(constants.SECRET_JWT), nil
		})
		if err != nil {
			return c.JSON(http.StatusUnauthorized, models.ResponseMessage{
				Message: "Unauthorized",
			})
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			return c.JSON(http.StatusUnauthorized, models.ResponseMessage{
				Message: "Unauthorized",
			})
		}
		id, ok := claims["id"].(float64)
		if !ok {
			return c.JSON(http.StatusUnauthorized, models.ResponseMessage{
				Message: "Unauthorized",
			})
		}
		c.Set("id", int(id))
		return next(c)
	}
}

// JWT Validator if Admin is Login
func JWTValidatorAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("bEDUCookie")
		if err != nil {
			return c.JSON(http.StatusUnauthorized, models.ResponseMessage{
				Message: "Unauthorized",
			})
		}
		token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
			return []byte(constants.SECRET_JWT), nil
		})
		if err != nil {
			return c.JSON(http.StatusUnauthorized, models.ResponseMessage{
				Message: "Unauthorized",
			})
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			return c.JSON(http.StatusUnauthorized, models.ResponseMessage{
				Message: "Unauthorized",
			})
		}
		role, ok := claims["role"].(string)
		if !ok {
			return c.JSON(http.StatusUnauthorized, models.ResponseMessage{
				Message: "Unauthorized",
			})
		}
		c.Set("role", string(role))
		return next(c)
	}
}
