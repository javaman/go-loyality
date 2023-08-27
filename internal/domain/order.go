package domain

import (
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
	Number     string
	Status     OrderStatus
	Accrual    int64
	UploadedAt time.Time
	Login      string
}

type OrderStoreUsecase interface {
	Store(o *Order) error
}

type OrderRepository interface {
	Insert(o *Order) error
	Select(orderNumber string) (*Order, error)
	// SelectAll(login string) ([]*Order, error)
}
