package domain

import (
	"encoding/json"
	"time"
)

type Withdraw struct {
	Order       string `json:"order"`
	Login       string
	Sum         int64 `json:"sum"`
	ProcessedAt time.Time
}

func (w *Withdraw) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Order       string `json:"order"`
		Sum         int64  `json:"sum"`
		ProcessedAt string `json:"processed_at"`
	}{
		Order:       w.Order,
		Sum:         w.Sum,
		ProcessedAt: w.ProcessedAt.Format(time.RFC3339),
	})
}

type StoreWithdrawUsecase interface {
	Store(w *Withdraw) error
}

type ListWithdrawsUsecase interface {
	List(login string) ([]*Withdraw, error)
}

type WithdrawRepository interface {
	Insert(w *Withdraw) error
	SelectAll(login string) ([]*Withdraw, error)
}
