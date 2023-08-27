package http

import (
	"io"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/javaman/go-loyality/internal/domain"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

type orderHandler struct {
	orderStoreUsecase domain.OrderStoreUsecase
	secret            string
}

func New(e *echo.Echo, secret string, orderStoreUsecase domain.OrderStoreUsecase) {
	handler := &orderHandler{
		secret:            secret,
		orderStoreUsecase: orderStoreUsecase,
	}

	config := echojwt.Config{
		SigningKey: []byte(secret),
	}

	r := e.Group("/api/user/orders")
	r.Use(echojwt.WithConfig(config))
	r.POST("", handler.StoreOrder)
}

func (h *orderHandler) StoreOrder(c echo.Context) error {
	b, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	orderNumber := string(b)

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	login, err := claims.GetSubject()

	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	o := &domain.Order{
		Number: orderNumber,
		Login:  login,
	}

	switch e := h.orderStoreUsecase.Store(o); e {
	case nil:
		return c.NoContent(http.StatusAccepted)
	case domain.ErrorBadOrderNumber:
		return c.NoContent(http.StatusUnprocessableEntity)
	case domain.ErrorOrderExistsSameUser:
		return c.NoContent(http.StatusOK)
	case domain.ErrorOrderExistsAnotherUser:
		return c.NoContent(http.StatusConflict)
	default:
		return c.NoContent(http.StatusInternalServerError)
	}
}
