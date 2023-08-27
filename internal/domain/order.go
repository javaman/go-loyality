package domain

import (
	"encoding/json"
	"errors"
	"time"
)

var (
	ErrorBadOrderNumber         error = errors.New("bad order number format")
	ErrorOrderExistsSameUser    error = errors.New("order exists")
	ErrorOrderExistsAnotherUser error = errors.New("order exists but different user")
)

type OrderStatus string

const (
	Registered OrderStatus = "REGISTERED"
	Invalid    OrderStatus = "INVALID"
	Processing OrderStatus = "PROCESSING"
	Processed  OrderStatus = "PROCESSED"
)

type Order struct {
	Number     string      `json:"number"`
	Status     OrderStatus `json:"status"`
	Accrual    int64       `json:"accrual"`
	UploadedAt time.Time   `json:"uploaded_at"`
	Login      string      `json:"-"`
}

func (o *Order) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Number     string      `json:"number"`
		Status     OrderStatus `json:"status"`
		Accrual    int64       `json:"accrual"`
		UploadedAt string      `json:"uploaded_at"`
	}{
		Number:     o.Number,
		Status:     o.Status,
		Accrual:    o.Accrual,
		UploadedAt: o.UploadedAt.Format(time.RFC3339),
	})
}

type OrderStoreUsecase interface {
	Store(o *Order) error
}

type OrderListUsecase interface {
	List(login string) ([]*Order, error)
}

type OrderRepository interface {
	Insert(o *Order) error
	Select(orderNumber string) (*Order, error)
	SelectAll(login string) ([]*Order, error)
}
