package services

import (
	db "AuthInGo/db/repositories"
	"fmt"
)

type UserService interface {
	GetUserByID() error
}

type UserServiceImpl struct {
	userRepository db.UserRepository
}

func NewUserServiceImpl(_userRepository db.UserRepository) UserService {
	return &UserServiceImpl{
		userRepository: _userRepository,
	}
}

func (u *UserServiceImpl) GetUserByID() error {
	fmt.Println("Getting user in services");

	u.userRepository.GetByID()
	return nil
}