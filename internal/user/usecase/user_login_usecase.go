package usecase

import "github.com/javaman/go-loyality/internal/domain"

type userLoginUsecase struct {
	userRepository domain.UserRepository
}

func NewUserLoginUsecase(ur domain.UserRepository) *userLoginUsecase {
	return &userLoginUsecase{ur}
}

func (uc *userLoginUsecase) Login(u *domain.User) (bool, error) {
	found, err := uc.userRepository.Select(u.Login)
	if err != nil {
		return false, err
	}
	if found == nil {
		return false, nil
	}
	return u.Login == found.Login, nil
}
