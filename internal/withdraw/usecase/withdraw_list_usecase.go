package usecase

import "github.com/javaman/go-loyality/internal/domain"

type withdrawListUsecase struct {
	repo domain.WithdrawRepository
}

func NewWithdrawListUsecase(r domain.WithdrawRepository) *withdrawListUsecase {
	return &withdrawListUsecase{repo: r}
}

func (u *withdrawListUsecase) List(login string) ([]*domain.Withdraw, error) {
	return u.repo.SelectAll(login)
}
