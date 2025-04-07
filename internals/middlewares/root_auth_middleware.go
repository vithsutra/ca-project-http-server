package middlewares

import (
	"errors"
	"log"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func RootMiddleware() echo.MiddlewareFunc {
	return middleware.BasicAuth(
		func(username, password string, ctx echo.Context) (bool, error) {
			rootUserName := os.Getenv("ROOT_USERNAME")
			log.Println("reqquest came ", os.Getenv("ROOT_USERNAME"), os.Getenv("ROOT_PASSWORD"))
			if rootUserName == "" {
				return false, errors.New("internal server error")
			}

			rootPassword := os.Getenv("ROOT_PASSWORD")
			if rootPassword == "" {
				return false, errors.New("internal server error")
			}

			if rootUserName != username {
				return false, errors.New("root username not matching")
			}

			if rootPassword != password {
				return false, errors.New("incorrect root password")
			}

			return true, nil
		},
	)
}
