package usecase

import (
	"github.com/ShiraazMoollatjie/goluhn"
	"github.com/javaman/go-loyality/internal/domain"
)

type withdrawStoreUsecase struct {
	repository domain.WithdrawRepository
}

func NewWithdrawStoreUsecase(r domain.WithdrawRepository) *withdrawStoreUsecase {
	return &withdrawStoreUsecase{repository: r}
}

func (u *withdrawStoreUsecase) Store(w *domain.Withdraw) error {
	err := goluhn.Validate(w.Order)

	if err != nil {
		return domain.ErrorBadOrderNumber
	}

	return u.repository.Insert(w)
}
