package adapters

import (
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/javaman/go-loyality/internal/domain"
)

type adapter struct {
	*resty.Client
}

func NewAccrualAdpater(endpoint string) *adapter {
	adapter := &adapter{resty.New()}
	adapter.SetBaseURL(endpoint)
	return adapter
}

func (a *adapter) Query(number string) (*domain.Order, error) {
	var order domain.Order
	r, err := a.R().SetResult(&order).Get("/api/orders/" + number)

	if err != nil {
		panic(err)
	}

	fmt.Println("====" + r.Status())
	return &order, nil
}
