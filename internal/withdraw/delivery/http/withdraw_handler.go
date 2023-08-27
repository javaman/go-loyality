package http

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/javaman/go-loyality/internal/domain"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

type withdrawHandler struct {
	storeWithdrawUsecase domain.StoreWithdrawUsecase
	listWithdrawsUsecase domain.ListWithdrawsUsecase
	secret               string
}

func New(e *echo.Echo, secret string, storeWithdrawUsecase domain.StoreWithdrawUsecase, listWithdrawsUsecase domain.ListWithdrawsUsecase) {
	handler := &withdrawHandler{
		secret:               secret,
		storeWithdrawUsecase: storeWithdrawUsecase,
		listWithdrawsUsecase: listWithdrawsUsecase,
	}

	config := echojwt.Config{
		SigningKey: []byte(secret),
	}

	r1 := e.Group("/api/user/balance/withdraw")
	r1.Use(echojwt.WithConfig(config))
	r1.POST("", handler.StoreWithdraw)

	r2 := e.Group("/api/user/withdrawls")
	r2.Use(echojwt.WithConfig(config))
	r2.GET("", handler.ListWithdrawls)
}

func getLogin(c echo.Context) (string, error) {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	return claims.GetSubject()
}

func (h *withdrawHandler) StoreWithdraw(c echo.Context) error {
	withdraw := new(domain.Withdraw)

	if err := c.Bind(withdraw); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	login, err := getLogin(c)

	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	withdraw.Login = login

	err = h.storeWithdrawUsecase.Store(withdraw)
	switch err {
	case nil:
		return c.NoContent(http.StatusOK)
	case domain.ErrorBadOrderNumber:
		return c.NoContent(http.StatusUnprocessableEntity)
	default:
		return c.NoContent(http.StatusInternalServerError)
	}
}

func (h *withdrawHandler) ListWithdrawls(c echo.Context) error {
	login, err := getLogin(c)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	result, err := h.listWithdrawsUsecase.List(login)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	if len(result) == 0 {
		return c.NoContent(http.StatusNoContent)
	}

	return c.JSON(http.StatusOK, result)
}
