package service

import (
	"fmt"
	"jxb-eprocurement/handlers"
	"jxb-eprocurement/handlers/dtos"
	"jxb-eprocurement/helpers"
	"jxb-eprocurement/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AuthService defines the methods for the auth service.
type AuthService interface {
	Login(c *gin.Context) handlers.ServiceResponse
}

// AuthServiceImpl is the implementation of the AuthService interface.
type AuthServiceImpl struct {
	db *gorm.DB
}

// NewAuthService creates a new instance of AuthServiceImpl.
func AuthServiceConstructor(db *gorm.DB) AuthService {
	return &AuthServiceImpl{db: db}
}

// Validate user input that validator cannot check by validator v10 package
func (a *AuthServiceImpl) inputValidator(input dtos.InputLoginDTO) (models.USR_User, bool) {
	// Setup variable
	var user models.USR_User
	is_error := false

	// Check user with email  or username exist
	if result := a.db.Preload("Role").Limit(1).Where("email = ?", input.UsernameOrEmail).Or("username = ?", input.UsernameOrEmail).Find(&user); result.Error != nil || result.RowsAffected == 0 {
		is_error = true
	}

	// Check password form hashes form
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		is_error = true
	}

	return user, is_error
}

// AddAuthData adds a new user to the database.
func (a *AuthServiceImpl) Login(c *gin.Context) handlers.ServiceResponse {
	var input dtos.InputLoginDTO

	if err := c.ShouldBind(&input); err != nil {
		return handlers.ServiceResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid Input",
			Data:    nil,
			Err:     err.Error(),
		}
	}

	// Validate input using golang validator
	if err := handlers.ValidateStruct(input); err != nil {
		errors := handlers.ValidationErrorHandlerV1(c, err, input)
		return handlers.ServiceResponse{
			Status:  http.StatusBadRequest,
			Message: "Error Invalid Data",
			Data:    nil,
			Err:     errors,
		}
	}

	// Check and validate input that cannot be validate by golang validator
	user, err := a.inputValidator(input)
	if err {
		return handlers.ServiceResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid email or password",
			Data:    nil,
			Err:     "",
		}
	}

	token, error := helpers.GenerateJWT(user)
	if error != nil {
		return handlers.ServiceResponse{
			Status:  http.StatusBadRequest,
			Message: "Failed to generate token",
			Data:    nil,
			Err:     error.Error(),
		}
	}

	data := map[string]interface{}{
		"loginAt": time.Now(),
		"user":    user.Name,
		"role":    user.Role.Name,
		"token":   fmt.Sprintf("Bearer %s", token),
	}

	return handlers.ServiceResponse{
		Status:  http.StatusOK,
		Message: "User Login Successfully",
		Data:    data,
		Err:     nil,
	}
}
