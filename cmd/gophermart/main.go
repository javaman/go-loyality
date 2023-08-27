package main

import (
	"github.com/labstack/echo/v4"

	"github.com/javaman/go-loyality/internal/config"
	orderhandler "github.com/javaman/go-loyality/internal/order/delivery/http"
	orderrepo "github.com/javaman/go-loyality/internal/order/repository/postgres"
	orderusecases "github.com/javaman/go-loyality/internal/order/usecase"
	userhandler "github.com/javaman/go-loyality/internal/user/delivery/http"
	userrepo "github.com/javaman/go-loyality/internal/user/repository/postgres"
	userusecases "github.com/javaman/go-loyality/internal/user/usecase"
)

func main() {
	cfg := config.Configure()

	ur := userrepo.NewUserRepository(cfg.DatabaseURI)
	or := orderrepo.NewOrderRepository(cfg.DatabaseURI)

	e := echo.New()

	userhandler.New(e, "iddqd", userusecases.NewUserRegisterUsecase(ur), userusecases.NewUserLoginUsecase(ur))
	orderhandler.New(e, "iddqd", orderusecases.NewOrderStoreUsecase(or))

	e.Logger.Fatal(e.Start(cfg.Address))
}
