package http

import (
	"io"
	"net/http"

	mwr "github.com/javaman/go-loyality/internal/delivery/http"
	"github.com/javaman/go-loyality/internal/domain"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

type orderHandler struct {
	orderStoreUsecase domain.OrderStoreUsecase
	orderListUsecase  domain.OrderListUsecase
	secret            string
}

func New(e *echo.Echo, secret string, orderStoreUsecase domain.OrderStoreUsecase, orderListUsecase domain.OrderListUsecase) {
	handler := &orderHandler{
		secret:            secret,
		orderStoreUsecase: orderStoreUsecase,
		orderListUsecase:  orderListUsecase,
	}

	config := echojwt.Config{
		SigningKey: []byte(secret),
	}

	r := e.Group("/api/user/orders")

	r.Use(echojwt.WithConfig(config))
	r.Use(mwr.ExtractLogin)
	r.POST("", handler.StoreOrder)
	r.GET("", handler.List)
}

func (h *orderHandler) List(c echo.Context) error {
	login := c.Get("Login").(string)

	list, err := h.orderListUsecase.List(login)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	if len(list) == 0 {
		return c.NoContent(http.StatusNoContent)
	}

	return c.JSON(http.StatusOK, list)
}

func (h *orderHandler) StoreOrder(c echo.Context) error {
	b, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	orderNumber := string(b)

	login := c.Get("Login").(string)

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
