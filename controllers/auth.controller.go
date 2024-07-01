package controllers

import (
	"jxb-eprocurement/handlers"
	"jxb-eprocurement/service"

	"github.com/gin-gonic/gin"
)

type AuthController interface {
	LoginUser(c *gin.Context)
}

// AuthControllerImpl is the implementation of the AuthController interface.
type AuthControllerImpl struct {
	service service.AuthService
}

// AuthControllerConstructor creates a new instance of AuthControllerImpl.
func AuthControllerConstructor(service service.AuthService) AuthController {
	return &AuthControllerImpl{service: service}
}

// LoginUser handles the request for user to login and get token for authentication.
func (ac *AuthControllerImpl) LoginUser(c *gin.Context) {
	response := ac.service.Login(c)
	if response.Err != nil {
		handlers.ResponseFormatter(c, response.Status, response.Err, response.Message)
	} else {
		handlers.ResponseFormatter(c, response.Status, response.Data, response.Message)
	}
}
