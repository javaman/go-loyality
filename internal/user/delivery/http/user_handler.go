package http

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/javaman/go-loyality/internal/domain"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userRegisterUsecase domain.UserRegisterUsecase
	userLoginUsecase    domain.UserLoginUsecase
	secret              string
}

func New(e *echo.Echo, secret string, userRegisterUsecase domain.UserRegisterUsecase, userLoginUsecase domain.UserLoginUsecase) {
	handler := &UserHandler{
		userRegisterUsecase: userRegisterUsecase,
		userLoginUsecase:    userLoginUsecase,
		secret:              secret,
	}
	e.POST("/api/user/register", handler.Register)
	e.POST("/api/user/login", handler.Login)
}

func (h *UserHandler) Register(c echo.Context) error {
	u := new(domain.User)

	if err := c.Bind(u); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	switch err := h.userRegisterUsecase.Register(u); err {
	case domain.ErrorLoginExists:
		return c.NoContent(http.StatusConflict)
	case nil:
		if e := h.setAuthorizationHeader(c, u.Login); e != nil {
			return c.NoContent(http.StatusInternalServerError)
		}
		return c.NoContent(http.StatusOK)
	default:
		return c.NoContent(http.StatusInternalServerError)
	}
}

func (h *UserHandler) Login(c echo.Context) error {
	u := new(domain.User)

	if err := c.Bind(u); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	ok, err := h.userLoginUsecase.Login(u)

	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	if !ok {
		return c.NoContent(http.StatusUnauthorized)
	}

	err = h.setAuthorizationHeader(c, u.Login)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusOK)
}

func (h *UserHandler) setAuthorizationHeader(c echo.Context, login string) error {

	claims := jwt.RegisteredClaims{
		Subject:   login,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 1)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(h.secret))

	if err != nil {
		return err
	}

	c.Response().Header().Set("Authorization", "Bearer "+t)
	return nil
}
