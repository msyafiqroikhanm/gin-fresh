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
	handlers.ResponseFormatterWithLogging(c, response)
}

// GetUser handles the request to get a user by ID.
func (uc *UserControllerImpl) GetUser(c *gin.Context) {
	response := uc.service.GetByID(c)
	handlers.ResponseFormatterWithLogging(c, response)
}

// CreateUser handles the request to add a new user.
func (uc *UserControllerImpl) CreateUser(c *gin.Context) {
	response := uc.service.AddData(c)
	handlers.ResponseFormatterWithLogging(c, response)
}

// UpdateUser handles the request to update a user.
func (uc *UserControllerImpl) UpdateUser(c *gin.Context) {
	response := uc.service.UpdateData(c)
	handlers.ResponseFormatterWithLogging(c, response)
}

// DeleteUser handles the request to delete a user.
func (uc *UserControllerImpl) DeleteUser(c *gin.Context) {
	response := uc.service.DeleteData(c)
	handlers.ResponseFormatterWithLogging(c, response)
}

// ResetPassUser handle the request to reset user's password by admin
func (uc *UserControllerImpl) ResetPassUser(c *gin.Context) {
	response := uc.service.ResetPass(c)
	handlers.ResponseFormatterWithLogging(c, response)
}

// ChangePassUser handle the request to change user's password by user it self
func (uc *UserControllerImpl) ChangePassUser(c *gin.Context) {
	response := uc.service.ChangePass(c)
	handlers.ResponseFormatterWithLogging(c, response)
}
