package adapters

import (
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/javaman/go-loyality/internal/domain"
)

type adapter struct {
	*resty.Client
}

func NewAccrualAdapter(endpoint string) *adapter {
	adapter := &adapter{resty.New()}
	adapter.SetBaseURL(endpoint)
	return adapter
}

func (a *adapter) Query(number string) (*domain.Order, error) {
	var order domain.Order
	r, err := a.R().SetResult(&order).Get("/api/orders/" + number)
	if err != nil {
		return nil, err
	}

	if r.StatusCode() == http.StatusNoContent {
		return nil, domain.ErrorOrderNotFound
	}

	if r.StatusCode() == http.StatusTooManyRequests {
		return nil, domain.ErrorTooFast
	}

	return &order, nil
}
