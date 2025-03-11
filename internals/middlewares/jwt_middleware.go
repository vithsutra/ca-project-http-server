package middlewares

import (
	"os"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func JwtMiddleware() echo.MiddlewareFunc {
	secretKey := os.Getenv("JWT_TOKEN_SCRETE_KEY")
	return echojwt.JWT([]byte(secretKey))
}
