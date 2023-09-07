package domain

import "encoding/json"

type Balance struct {
	Current   json.Number `json:"current"`
	Withdrawn json.Number `json:"withdrawn"`
}

type CheckBalanceUsecase interface {
	Check(login string) (Balance, error)
}

type BalanceRepository interface {
	Select(login string) (Balance, error)
}
