package usecase

import (
	"github.com/javaman/go-loyality/internal/domain"
)

type orderListUsecase struct {
	orderRepository domain.OrderRepository
}

func NewOrderListUsecase(repo domain.OrderRepository) *orderListUsecase {
	return &orderListUsecase{repo}
}

func (uc *orderListUsecase) List(login string) ([]*domain.Order, error) {
	return uc.orderRepository.SelectAll(login)
}
