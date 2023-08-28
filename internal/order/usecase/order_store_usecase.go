package usecase

import (
	"github.com/ShiraazMoollatjie/goluhn"
	"github.com/javaman/go-loyality/internal/domain"
)

type orderStoreUsecase struct {
	orderRepository domain.OrderRepository
	details         domain.OrderDetails
}

func NewOrderStoreUsecase(repo domain.OrderRepository, details domain.OrderDetails) *orderStoreUsecase {
	return &orderStoreUsecase{repo, details}
}

func (uc *orderStoreUsecase) Store(o *domain.Order) error {
	err := goluhn.Validate(o.Number)

	if err != nil {
		return domain.ErrorBadOrderNumber
	}

	storedOrder, err := uc.orderRepository.Select(o.Number)

	if err != nil {
		return err
	}

	if storedOrder != nil {
		if o.Login == storedOrder.Login {
			return domain.ErrorOrderExistsSameUser
		} else {
			return domain.ErrorOrderExistsAnotherUser
		}
	}

	orderDetails, _ := uc.details.Query(o.Number)

	o.Accrual = orderDetails.Accrual
	o.Status = orderDetails.Status

	if len(o.Status) == 0 {
		o.Status = domain.NEW
	}

	err = uc.orderRepository.Insert(o)
	if err != nil {
		return err
	}

	return nil
}
