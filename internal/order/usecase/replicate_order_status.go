package usecase

import (
	"sync"
	"time"

	"github.com/javaman/go-loyality/internal/domain"
	"go.uber.org/ratelimit"
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

			var wg sync.WaitGroup
			ratelimit := ratelimit.New(10) // 10 timeslots per second

			for _, order := range tenOrders {

				wg.Add(1)
				go uc.updateOrder(order, ratelimit, &wg)
			}

			wg.Wait()
		}
	}()
}

func (uc *replicateOrderStatusUsecase) updateOrder(o *domain.Order, rl ratelimit.Limiter, wg *sync.WaitGroup) {
	defer func() { wg.Done() }()

	rl.Take()
	updatedOrder, err := uc.orderDetails.Query(o.Number)

	switch err {
	case nil:
		o.Accrual = updatedOrder.Accrual
		o.Status = updatedOrder.Status
		uc.orderRepository.Update(o, -1)
	case domain.ErrorTooFast:
		for i := 0; i < 60000; i++ {
			rl.Take()
		}
	}
}
