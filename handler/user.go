package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-backer/auth"
	"go-backer/helper"
	"go-backer/user"
	"net/http"
)

type userHandler struct {
	userService user.Service
	authService auth.Service
}

func NewUserHandler(userService user.Service, authService auth.Service) *userHandler {
	return &userHandler{userService, authService}
}

func (h *userHandler) RegisterUser(c *gin.Context) {
	//	get input from user
	var input user.RegisterUserInput
	//mapping from input to Register User Input
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Register Account Failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	//	call user service with mapped input as parameter
	newUser, err := h.userService.RegisterUser(input)
	if err != nil {
		response := helper.APIResponse("Register Account Failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	token, err := h.authService.GenerateToken(newUser.ID)
	if err != nil {
		response := helper.APIResponse("Register Account Failed: Generate Token Failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	formatter := user.FormatUser(newUser, token)

	response := helper.APIResponse("Account has been registered", http.StatusOK, "success", formatter)
	c.JSON(http.StatusOK, response)
}

func (h *userHandler) Login(c *gin.Context) {
	//	user input password
	var input user.LoginInput

	//	input processed by handler
	//	map from input into input struct
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Login Failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	//	call service with input struct as parameter
	loggedInUser, err := h.userService.Login(input)
	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}

		response := helper.APIResponse("Login Failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	token, err := h.authService.GenerateToken(loggedInUser.ID)
	if err != nil {
		response := helper.APIResponse("Login Failed: Generate Token Failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := user.FormatUser(loggedInUser, token)
	response := helper.APIResponse("Successfully loggedIn", http.StatusOK, "success", formatter)
	c.JSON(http.StatusOK, response)
}

func (h *userHandler) CheckEmailAvailability(c *gin.Context) {
	//	input email from user
	var input user.CheckEmailInput

	//	map input to struct input + validation
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Check Email Failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	//	call service with struct input as parameter
	//  service will call appropriate repository
	//	repository will determine whether email is already exist or not
	isEmailAvailable, err := h.userService.IsEmailAvailable(input)
	if err != nil {
		errorMessage := gin.H{"errors": "Server error"}

		response := helper.APIResponse("Check Email Failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	//gin.H used for mapping
	data := gin.H{"is_available": isEmailAvailable}

	message := "Email has been registered!"
	if isEmailAvailable {
		message = "Email is available"
	}

	response := helper.APIResponse(message, http.StatusOK, "succeed", data)
	c.JSON(http.StatusOK, response)

}

func (h *userHandler) UploadAvatar(c *gin.Context) {
	//	get input from user
	//	save images in folder named "images/"
	//	service will call repository
	//	will use JWT, but for now will set static ID = 1
	//	repo will get user ID =1
	//	update data user with file path

	//get file from input
	//"avatar" is the field name of input file
	file, err := c.FormFile("avatar")
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Failed to upload avatar image", http.StatusBadRequest, "error", data)

		c.JSON(http.StatusBadRequest, response)
		return
	}
	//	TODO: will get userId from JWT Token later, left it static now
	userId := 1
	//path => images/{userId}-{fileName}
	path := fmt.Sprintf("images/%d-%s", userId, file.Filename)

	//save upload file to some path + filename
	err = c.SaveUploadedFile(file, path)
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Failed to upload avatar image", http.StatusBadRequest, "error", data)

		c.JSON(http.StatusBadRequest, response)
		return
	}

	_, err = h.userService.SaveAvatar(userId, path)
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Failed to upload avatar image", http.StatusBadRequest, "error", data)

		c.JSON(http.StatusBadRequest, response)
		return
	}

	data := gin.H{"is_uploaded": true}
	response := helper.APIResponse("Avatar successfully uploaded", http.StatusBadRequest, "error", data)

	c.JSON(http.StatusBadRequest, response)
}
