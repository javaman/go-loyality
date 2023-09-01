package http

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func getLogin(c echo.Context) (string, error) {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	return claims.GetSubject()
}

func ExtractLogin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		login, err := getLogin(c)
		if err != nil {
			return c.NoContent(http.StatusUnauthorized)
		}
		c.Set("Login", login)
		return next(c)
	}
}
