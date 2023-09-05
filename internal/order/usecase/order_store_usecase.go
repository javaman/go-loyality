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

	err = uc.orderRepository.Insert(o)

	switch err {
	case nil:
		return uc.orderJustCreated(o)
	case domain.ErrorOrderExists:
		return uc.orderExists(o)
	default:
		return err
	}

}

func (uc *orderStoreUsecase) orderExists(o *domain.Order) error {
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

	return nil
}

func (uc *orderStoreUsecase) orderJustCreated(o *domain.Order) error {
	go func() {
		queryResult, err := uc.details.Query(o.Number)
		if err != nil {
			queryResult.Number = o.Number
			uc.orderRepository.Update(queryResult, 0)
		}
	}()
	return nil
}
