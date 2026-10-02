package controllers

import (
	"AuthInGo/services"
	"net/http"
)

// we are not writing interfaces for controllers because
// controllers do not have multiple implementation in our
// project that's why we are calling controllers from router
// via concrete class

type UserController struct {
	UserService services.UserService
}

func NewUserController(_userService services.UserService) *UserController {
	return &UserController{
		UserService: _userService,
	}
}

func (uc *UserController) RegisterUser(w http.ResponseWriter, r *http.Request) {
	uc.UserService.CreateUser()
	w.Write([]byte("User registration endpoint"))
}