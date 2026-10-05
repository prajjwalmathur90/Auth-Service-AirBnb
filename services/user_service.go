package services

import (
	dto "AuthInGo/Dto"
	db "AuthInGo/db/repositories"
	"AuthInGo/utils"
	"fmt"

	env "AuthInGo/config/env"

	"github.com/golang-jwt/jwt/v4"
)

type UserService interface {
	GetUserByID(id int) error
	CreateUser() error
	LoginUser(payload *dto.LoginUserRequestDto) (string, error)
}

type UserServiceImpl struct {
	userRepository db.UserRepository
}

func NewUserServiceImpl(_userRepository db.UserRepository) UserService {
	return &UserServiceImpl{
		userRepository: _userRepository,
	}
}

func (u *UserServiceImpl) GetUserByID(id int) error {

	

	u.userRepository.GetByID(id)
	return nil
}

func (u *UserServiceImpl) CreateUser() error {
	return nil
}

func (u *UserServiceImpl) LoginUser(payload *dto.LoginUserRequestDto) (string, error) {
	user, err := u.userRepository.GetByEmail(payload.Email)

	if err != nil {
		fmt.Println("Error fetching user by email :", payload.Email)
		return "", err
	}

	if user == nil {
		fmt.Println("No user found with this email")
		return "", err
	}

	isPasswordValid := utils.CheckHashedPassword(payload.Password, user.Password)

	if !isPasswordValid {
		fmt.Println("Invalid password")
		return "", nil
	}

	jwtPayload := jwt.MapClaims{
		"id": user.Id,
		"email": user.Email,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtPayload)

	tokenString, err := token.SignedString([]byte(env.GetString("JWT_SECRET", "TOKEN")))

	if err != nil {
		fmt.Println("Error signing token : ", err)
		return "", err
	}

	fmt.Println("JWT Token :", tokenString)
	
	return tokenString, nil
}

