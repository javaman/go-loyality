package usecase

import (
	"time"

	"github.com/javaman/go-loyality/internal/domain"
)

type replicateOrderStatusUsecase struct {
	orderRepository domain.OrderRepository
	orderDetails    domain.OrderDetails
}

func NewReplicateOrderStatusUsecase(repo domain.OrderRepository, details domain.OrderDetails) *replicateOrderStatusUsecase {
	return &replicateOrderStatusUsecase{repo, details}
}

func (uc *replicateOrderStatusUsecase) Run() {
	go func() {
		for {
			time.Sleep(300 * time.Millisecond)
			tenOrders, err := uc.orderRepository.SelectTenOrders()
			if err != nil {
				continue
			}
			for _, order := range tenOrders {
				updatedOrder, err := uc.orderDetails.Query(order.Number)
				if err != nil {
					continue
				}
				order.Accrual = updatedOrder.Accrual
				order.Status = updatedOrder.Status
				uc.orderRepository.Update(order, -1)
			}
		}
	}()
}
