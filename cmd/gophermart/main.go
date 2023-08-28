package main

import (
	"fmt"

	"github.com/labstack/echo/v4"

	balancehandler "github.com/javaman/go-loyality/internal/balance/delivery/http"
	balancerepo "github.com/javaman/go-loyality/internal/balance/repository/postgres"
	balanceusecases "github.com/javaman/go-loyality/internal/balance/usecase"
	"github.com/javaman/go-loyality/internal/config"
	"github.com/javaman/go-loyality/internal/order/adapters"
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
	br := balancerepo.NewBalanceRepository(cfg.DatabaseURI)

	fmt.Println()

	e := echo.New()

	userhandler.New(e, "iddqd", userusecases.NewUserRegisterUsecase(ur), userusecases.NewUserLoginUsecase(ur))
	orderhandler.New(e, "iddqd", orderusecases.NewOrderStoreUsecase(or, adapters.NewAccrualAdpater(cfg.AccrualSystemAddress)), orderusecases.NewOrderListUsecase(or))
	withdrawhandler.New(e, "iddqd", withdrawusecases.NewWithdrawStoreUsecase(wr), withdrawusecases.NewWithdrawListUsecase(wr))
	balancehandler.New(e, "iddqd", balanceusecases.NewCheckBalanceUsecase(br))

	e.Logger.Fatal(e.Start(cfg.Address))
}
