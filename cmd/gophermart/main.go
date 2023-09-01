package main

import (
	"github.com/labstack/echo/v4"

	balancehandler "github.com/javaman/go-loyality/internal/balance/delivery/http"
	balancerepo "github.com/javaman/go-loyality/internal/balance/repository/postgres"
	balanceusecases "github.com/javaman/go-loyality/internal/balance/usecase"
	"github.com/javaman/go-loyality/internal/config"
	adapters "github.com/javaman/go-loyality/internal/order/adapters/accrual"
	orderhandler "github.com/javaman/go-loyality/internal/order/delivery/http"
	orderrepo "github.com/javaman/go-loyality/internal/order/repository/postgres"
	orderusecases "github.com/javaman/go-loyality/internal/order/usecase"
	userhandler "github.com/javaman/go-loyality/internal/user/delivery/http"
	userrepo "github.com/javaman/go-loyality/internal/user/repository/postgres"
	userusecases "github.com/javaman/go-loyality/internal/user/usecase"
	withdrawlhandler "github.com/javaman/go-loyality/internal/withdraw/delivery/http"
	withdrawlrepo "github.com/javaman/go-loyality/internal/withdraw/repository/postgres"
	withdrawlusecases "github.com/javaman/go-loyality/internal/withdraw/usecase"
)

func main() {
	cfg := config.Configure()

	ur := userrepo.NewUserRepository(cfg.DatabaseURI)
	or := orderrepo.NewOrderRepository(cfg.DatabaseURI)
	wr := withdrawlrepo.NewWithdrawRepository(cfg.DatabaseURI)
	br := balancerepo.NewBalanceRepository(cfg.DatabaseURI)

	e := echo.New()

	userhandler.New(e, cfg.Secret, userusecases.NewUserRegisterUsecase(ur), userusecases.NewUserLoginUsecase(ur))
	orderhandler.New(e, cfg.Secret, orderusecases.NewOrderStoreUsecase(or, adapters.NewAccrualAdapter(cfg.AccrualSystemAddress)), orderusecases.NewOrderListUsecase(or))
	withdrawlhandler.New(e, cfg.Secret, withdrawlusecases.NewWithdrawStoreUsecase(wr), withdrawlusecases.NewWithdrawListUsecase(wr))
	balancehandler.New(e, cfg.Secret, balanceusecases.NewCheckBalanceUsecase(br))

	e.Logger.Fatal(e.Start(cfg.Address))
}
