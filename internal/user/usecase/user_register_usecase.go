package usecase

import "github.com/javaman/go-loyality/internal/domain"

type userRegisterUsecase struct {
	userRepository domain.UserRepository
}

func NewUserRegisterUsecase(ur domain.UserRepository) *userRegisterUsecase {
	return &userRegisterUsecase{ur}
}

func (uc *userRegisterUsecase) Register(u *domain.User) error {
	err := uc.userRepository.Insert(u)
	return err
}
