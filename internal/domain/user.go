package domain

import "errors"

type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

var (
	ErrorLoginExists error = errors.New("user with same login exists")
)

type UserRegisterUsecase interface {
	Register(u *User) error
}

type UserLoginUsecase interface {
	Login(u *User) (bool, error)
}

type UserRepository interface {
	Insert(u *User) error
	Select(login string) (*User, error)
}
