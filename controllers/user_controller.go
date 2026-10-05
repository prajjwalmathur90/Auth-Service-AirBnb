package controllers

import (
	dto "AuthInGo/Dto"
	"AuthInGo/services"
	"AuthInGo/utils"
	"fmt"
	"net/http"
	"strconv"
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
	id, err := strconv.Atoi(r.PathValue("id"))

	if err != nil {
		fmt.Println("Error converting id to int : ", err)
		return
	}
	
	uc.UserService.GetUserByID(id)
	w.Write([]byte("Getting user by ID endpoint"))
}

func (uc *UserController) LoginUser(w http.ResponseWriter, r *http.Request) {
	var payload dto.LoginUserRequestDto;

	if jsonErr := utils.ReadJsonBody(r, &payload); jsonErr != nil {
		utils.WriteJsonErrorResponse(w, http.StatusBadRequest, "Invalid Json Body", jsonErr)
		return
	}

	if validationErr := utils.Validator.Struct(payload); validationErr != nil {
		utils.WriteJsonErrorResponse(w, http.StatusBadRequest, "Invalid Data!!", validationErr)
		return
	}

	jwtToken, err := uc.UserService.LoginUser(&payload)

	if err != nil {
		utils.WriteJsonErrorResponse(w, http.StatusInternalServerError, "Failed to login user", err)
		return
	}

	utils.WriteJsonSuccessResponse(w, http.StatusOK, "User logged in successfully!", jwtToken)
}