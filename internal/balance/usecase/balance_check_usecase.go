package usecase

import "github.com/javaman/go-loyality/internal/domain"

type checkBalanceUsecase struct {
	repo domain.BalanceRepository
}

func NewCheckBalanceUsecase(r domain.BalanceRepository) *checkBalanceUsecase {
	return &checkBalanceUsecase{repo: r}
}

func (u *checkBalanceUsecase) Check(login string) (domain.Balance, error) {
	return u.repo.Select(login)
}
