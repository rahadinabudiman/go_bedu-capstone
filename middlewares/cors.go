package middlewares

import "github.com/labstack/echo/v4"

// Middleware untuk mengizinkan CORS
func AllowCORS(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderAccessControlAllowOrigin, "*")
		c.Response().Header().Set(echo.HeaderAccessControlAllowHeaders, "Content-Type, Authorization")
		c.Response().Header().Set(echo.HeaderAccessControlAllowMethods, "GET, POST, PUT, DELETE")
		return next(c)
	}
}
