package controllers

import (
	"jxb-eprocurement/handlers"
	"jxb-eprocurement/service"

	"github.com/gin-gonic/gin"
)

type USRUserController interface {
	GetAllUsers(c *gin.Context)
	GetUser(c *gin.Context)
	CreateUser(c *gin.Context)
	UpdateUser(c *gin.Context)
	DeleteUser(c *gin.Context)
	ChangePassUser(c *gin.Context)
	ResetPassUser(c *gin.Context)
}

// UserControllerImpl is the implementation of the UserController interface.
type UserControllerImpl struct {
	service service.UserService
}

// NewUserController creates a new instance of UserControllerImpl.
func UserControllerConstructor(service service.UserService) USRUserController {
	return &UserControllerImpl{service: service}
}

// GetAllUsers handles the request to get all users.
func (uc *UserControllerImpl) GetAllUsers(c *gin.Context) {
	response := uc.service.GetAll(c)
	if response.Err != nil {
		handlers.ResponseFormatter(c, response.Status, response.Err, response.Message)
	} else {
		handlers.ResponseFormatter(c, response.Status, response.Data, response.Message)
	}
}

// GetUser handles the request to get a user by ID.
func (uc *UserControllerImpl) GetUser(c *gin.Context) {
	response := uc.service.GetByID(c)
	if response.Err != nil {
		handlers.ResponseFormatter(c, response.Status, response.Err, response.Message)
	} else {
		handlers.ResponseFormatter(c, response.Status, response.Data, response.Message)
	}
}

// CreateUser handles the request to add a new user.
func (uc *UserControllerImpl) CreateUser(c *gin.Context) {
	response := uc.service.AddData(c)
	if response.Err != nil {
		handlers.ResponseFormatter(c, response.Status, response.Err, response.Message)
	} else {
		handlers.ResponseFormatter(c, response.Status, response.Data, response.Message)
	}
}

// UpdateUser handles the request to update a user.
func (uc *UserControllerImpl) UpdateUser(c *gin.Context) {
	response := uc.service.UpdateData(c)
	if response.Err != nil {
		handlers.ResponseFormatter(c, response.Status, response.Err, response.Message)
	} else {
		handlers.ResponseFormatter(c, response.Status, response.Data, response.Message)
	}
}

// DeleteUser handles the request to delete a user.
func (uc *UserControllerImpl) DeleteUser(c *gin.Context) {
	response := uc.service.DeleteData(c)
	if response.Err != nil {
		handlers.ResponseFormatter(c, response.Status, response.Err, response.Message)
	} else {
		handlers.ResponseFormatter(c, response.Status, response.Data, response.Message)
	}
}

// ResetPassUser handle the request to reset user's password by admin
func (uc *UserControllerImpl) ResetPassUser(c *gin.Context) {
	response := uc.service.ResetPass(c)
	if response.Err != nil {
		handlers.ResponseFormatter(c, response.Status, response.Err, response.Message)
	} else {
		handlers.ResponseFormatter(c, response.Status, response.Data, response.Message)
	}
}

// ChangePassUser handle the request to change user's password by user it self
func (uc *UserControllerImpl) ChangePassUser(c *gin.Context) {
	response := uc.service.ChangePass(c)
	if response.Err != nil {
		handlers.ResponseFormatter(c, response.Status, response.Err, response.Message)
	} else {
		handlers.ResponseFormatter(c, response.Status, response.Data, response.Message)
	}
}
