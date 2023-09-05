package domain

import (
	"encoding/json"
	"time"
)

type OrderStatus string

const (
	Registered OrderStatus = "REGISTERED"
	Invalid    OrderStatus = "INVALID"
	Processing OrderStatus = "PROCESSING"
	Processed  OrderStatus = "PROCESSED"
	NEW        OrderStatus = "NEW"
)

type Order struct {
	Number     string      `json:"order"`
	Status     OrderStatus `json:"status"`
	Accrual    json.Number `json:"accrual"`
	UploadedAt time.Time   `json:"uploaded_at"`
	Login      string      `json:"-"`
}

func (o *Order) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Number     string      `json:"number"`
		Status     OrderStatus `json:"status"`
		Accrual    json.Number `json:"accrual"`
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
	Update(o *Order, version int) (bool, error)
	SelectTenOrders() ([]*Order, error)
}

type OrderDetails interface {
	Query(number string) (*Order, error)
}
