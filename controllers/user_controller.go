package controllers

import (
	dto "AuthInGo/Dto"
	"AuthInGo/middleware"
	"AuthInGo/services"
	"AuthInGo/utils"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
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

func (uc *UserController) GetUserByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))

	if err != nil {
		fmt.Println("Error converting id to int : ", err)
		return
	}
	
	user, err := uc.UserService.GetUserByID(id)

	if err != nil {
		utils.WriteJsonErrorResponse(w, http.StatusInternalServerError, "Failed to get user by id", err)
		return
	}

	if user == nil {
		utils.WriteJsonErrorResponse(w, http.StatusNotFound, "User not found with this id", nil)
		return
	}

	utils.WriteJsonSuccessResponse(w, http.StatusOK, "User fetched successfully!", user)
}

func (uc *UserController) GetAll(w http.ResponseWriter, r *http.Request) {
	users, err := uc.UserService.GetAll()

	if err != nil {
		utils.WriteJsonErrorResponse(w, http.StatusInternalServerError, "Failed to get all users", err)
		return
	}

	if len(users) == 0 {
		utils.WriteJsonErrorResponse(w, http.StatusNotFound, "No users found", nil)
		return
	}

	utils.WriteJsonSuccessResponse(w, http.StatusOK, "Users fetched successfully!", users)
}

func (uc *UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
	payload := r.Context().Value(middleware.PayloadKey).(*dto.CreateUserRequestDto)

	user, err := uc.UserService.CreateUser(payload)

	if err != nil {
		utils.WriteJsonErrorResponse(w, http.StatusInternalServerError, "Failed to create user", err)
		return
	}

	if user == nil {
		utils.WriteJsonErrorResponse(w, http.StatusConflict, "User already exists with this email or username", nil)
		return
	}

	utils.WriteJsonSuccessResponse(w, http.StatusOK, "User created successfully!", user)
}


func (uc *UserController) LoginUser(w http.ResponseWriter, r *http.Request) {
	payload := r.Context().Value(middleware.PayloadKey).(*dto.LoginUserRequestDto)

	jwtToken, err := uc.UserService.LoginUser(payload)

	if err != nil {
		utils.WriteJsonErrorResponse(w, http.StatusInternalServerError, "Failed to login user", err)
		return
	}

	utils.WriteJsonSuccessResponse(w, http.StatusOK, "User logged in successfully!", jwtToken)
}

func (uc *UserController) DeleteById(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))

	if err != nil {
		utils.WriteJsonErrorResponse(w, http.StatusBadRequest, "Invalid Id", err)
		return
	}

	err = uc.UserService.DeleteById(id)

	if err != nil {
		utils.WriteJsonErrorResponse(w, http.StatusInternalServerError, "Failed to delete user", err)
		return
	}

	utils.WriteJsonSuccessResponse(w, http.StatusOK, "User deleted successfully!", nil)
}
