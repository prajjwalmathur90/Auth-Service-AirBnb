package router

import (
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
	// public routes
	r.With(middleware.ValidateCreateUserRequest).Post("/signup", ur.userController.CreateUser)
	r.With(middleware.ValidateLoginRequest).Post("/login", ur.userController.LoginUser)
	r.Post("/logout", ur.userController.LogoutUser)

	// protected routes (require a valid JWT cookie)
	r.Group(func(pr chi.Router) {
		pr.Use(middleware.JWTAuthMiddleware)

		pr.Get("/{id}", ur.userController.GetUserByID)
		pr.Get("/", ur.userController.GetAll)
		pr.Delete("/{id}", ur.userController.DeleteById)
	})
}
