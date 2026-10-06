package services

import (
	dto "AuthInGo/Dto"
	db "AuthInGo/db/repositories"
	"AuthInGo/models"
	"AuthInGo/utils"
	"fmt"
	"strconv"

	env "AuthInGo/config/env"

	"github.com/golang-jwt/jwt/v4"
)

type UserService interface {
	GetUserByID(id string) (*models.User, error)
	GetAll() ([]*models.User, error)
	CreateUser(payload *dto.CreateUserRequestDto) (*models.User, error)
	LoginUser(payload *dto.LoginUserRequestDto) (string, error)
	DeleteById(id int) (error)
}

type UserServiceImpl struct {
	userRepository db.UserRepository
}

func NewUserServiceImpl(_userRepository db.UserRepository) UserService {
	return &UserServiceImpl{
		userRepository: _userRepository,
	}
}

func (u *UserServiceImpl) GetUserByID(id string) (*models.User, error) {
	user, err := u.userRepository.GetByID(id)

	if err != nil {
		fmt.Println("Error fetching user by id :", id)
		return nil, err
	}

	if user == nil {
		fmt.Println("No user found with this id")
		return nil, nil
	}

	fmt.Println("User fetched successfully! : ", user)
	
	return user, nil
}

func (u *UserServiceImpl) GetAll() ([]*models.User, error) {
	users, err := u.userRepository.GetAll()

	if err != nil {
		fmt.Println("Error fetching users : ", err)
		return nil, err
	}

	if len(users) == 0 {
		fmt.Println("No users found")
		return nil, nil
	}

	fmt.Println("Users fetched successfully! : ", users)
	
	return users, nil
}

func (u *UserServiceImpl) CreateUser(payload *dto.CreateUserRequestDto) (*models.User, error) {

	user, err := u.userRepository.GetByEmail(payload.Email)

	if err != nil {
		fmt.Println("Error fetching user by email : ", err)
		return nil, err
	}

	if user != nil {
		fmt.Println("User already exists with this email")
		return nil, nil
	}

	username, err := u.userRepository.GetByUsername(payload.Username)

	if err != nil {
		fmt.Println("Error fetching user by username : ", err)
		return nil, err
	}

	if username != nil {
		fmt.Println("User already exists with this username")
		return nil, nil
	}

	hashedPassword, err := utils.HashPassword(payload.Password)

	if err != nil {
		fmt.Println("Error hashing password : ", err)
		return nil, err
	}

	user,err = u.userRepository.Create(payload.Username, payload.Email, hashedPassword)

	if err != nil {
		fmt.Println("Error creating user : ", err)
		return nil, err
	}

	fmt.Println("User created successfully!")
	return user, nil
}

func (u *UserServiceImpl) LoginUser(payload *dto.LoginUserRequestDto) (string, error) {
	user, err := u.userRepository.GetByEmail(payload.Email)

	if err != nil {
		fmt.Println("Error fetching user by email :", payload.Email)
		return "", err
	}

	if user == nil {
		fmt.Println("No user found with this email")
		return "", fmt.Errorf("no user found with this email")
	}

	isPasswordValid := utils.CheckHashedPassword(payload.Password, user.Password)

	if !isPasswordValid {
		fmt.Println("Invalid password")
		return "", fmt.Errorf("invalid password")
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

func (u *UserServiceImpl) DeleteById(id int) (error) {

	user, err := u.userRepository.GetByID(strconv.Itoa(id))

	if err != nil {
		fmt.Println("Error fetching user by id : ", err)
		return err
	}

	if user == nil {
		fmt.Println("User not found with this id")
		return nil
	}

	err = u.userRepository.DeleteById(id)

	if err != nil {
		fmt.Println("Error deleting user : ", err)
		return err
	}

	fmt.Println("User deleted successfully!")
	return nil
}