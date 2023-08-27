package usecase

import (
	"github.com/ShiraazMoollatjie/goluhn"
	"github.com/javaman/go-loyality/internal/domain"
)

type orderStoreUsecase struct {
	orderRepository domain.OrderRepository
}

func NewOrderStoreUsecase(repo domain.OrderRepository) *orderStoreUsecase {
	return &orderStoreUsecase{repo}
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

	//TODO replace call with another system
	o.Accrual = 500
	o.Status = domain.Processed

	err = uc.orderRepository.Insert(o)
	if err != nil {
		return err
	}

	return nil
}
