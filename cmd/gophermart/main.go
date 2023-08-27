package main

import (
	"fmt"

	"github.com/labstack/echo/v4"

	"github.com/javaman/go-loyality/internal/config"
	orderhandler "github.com/javaman/go-loyality/internal/order/delivery/http"
	orderrepo "github.com/javaman/go-loyality/internal/order/repository/postgres"
	orderusecases "github.com/javaman/go-loyality/internal/order/usecase"
	userhandler "github.com/javaman/go-loyality/internal/user/delivery/http"
	userrepo "github.com/javaman/go-loyality/internal/user/repository/postgres"
	userusecases "github.com/javaman/go-loyality/internal/user/usecase"
	withdrawhandler "github.com/javaman/go-loyality/internal/withdraw/delivery/http"
	withdrawrepo "github.com/javaman/go-loyality/internal/withdraw/repository/postgres"
	withdrawusecases "github.com/javaman/go-loyality/internal/withdraw/usecase"
)

func main() {
	cfg := config.Configure()

	ur := userrepo.NewUserRepository(cfg.DatabaseURI)
	or := orderrepo.NewOrderRepository(cfg.DatabaseURI)
	wr := withdrawrepo.NewWithdrawRepository(cfg.DatabaseURI)

	fmt.Println(wr)

	e := echo.New()

	userhandler.New(e, "iddqd", userusecases.NewUserRegisterUsecase(ur), userusecases.NewUserLoginUsecase(ur))
	orderhandler.New(e, "iddqd", orderusecases.NewOrderStoreUsecase(or), orderusecases.NewOrderListUsecase(or))
	withdrawhandler.New(e, "iddqd", withdrawusecases.NewWithdrawStoreUsecase(wr), withdrawusecases.NewWithdrawListUsecase(wr))

	e.Logger.Fatal(e.Start(cfg.Address))
}
