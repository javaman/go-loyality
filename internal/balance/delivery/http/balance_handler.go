package http

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/javaman/go-loyality/internal/domain"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

type balanceHandler struct {
	checkBalanceUsecase domain.CheckBalanceUsecase
	secret              string
}

func New(e *echo.Echo, secret string, checkBalanceUsecase domain.CheckBalanceUsecase) *balanceHandler {
	handler := &balanceHandler{
		secret:              secret,
		checkBalanceUsecase: checkBalanceUsecase,
	}

	config := echojwt.Config{
		SigningKey: []byte(secret),
	}

	r1 := e.Group("/api/user/balance")
	r1.Use(echojwt.WithConfig(config))
	r1.GET("", handler.Check)

	return handler
}

func getLogin(c echo.Context) (string, error) {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	return claims.GetSubject()
}

func (h *balanceHandler) Check(c echo.Context) error {

	login, err := getLogin(c)

	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	balance, err := h.checkBalanceUsecase.Check(login)

	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, &balance)
}
