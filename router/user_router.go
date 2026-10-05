package router

import (
	dto "AuthInGo/Dto"
	"AuthInGo/controllers"
	"AuthInGo/middleware"

	"github.com/go-chi/chi/v5"
)

type UserRouter struct {
	userController *controllers.UserController
}

func NewUserRouter(_userController *controllers.UserController) Router {
	return &UserRouter{
		userController: _userController,
	}
}

func (ur *UserRouter) Register(r chi.Router) {
	r.Get("/{id}", ur.userController.GetUserByID)
	r.Get("/", ur.userController.GetAll)
	r.With(middleware.ValidateBody(func() any { return &dto.CreateUserRequestDto{} })).Post("/signup", ur.userController.CreateUser)
	r.With(middleware.ValidateBody(func() any { return &dto.LoginUserRequestDto{} })).Post("/login", ur.userController.LoginUser)
	r.Delete("/{id}", ur.userController.DeleteById)
}
