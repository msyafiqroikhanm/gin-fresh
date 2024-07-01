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
	Login(c *gin.Context) handlers.ServiceResponseWithLogging
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
func (a *AuthServiceImpl) inputValidator(input dtos.InputLoginDTO, c *gin.Context) (models.USR_User, bool) {
	// Setup variable
	var user models.USR_User
	isError := false

	// Create log
	log := helpers.CreateLog(c, a)

	// Check user with email  or username exist
	if result := a.db.Preload("Role").Preload("Role.Features").Preload("Role.Features.Module").Limit(1).Where("email = ?", input.UsernameOrEmail).Or("username = ?", input.UsernameOrEmail).Find(&user); result.Error != nil || result.RowsAffected == 0 {
		isError = true
	}

	// Check password form hashes form
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		isError = true
	}

	if isError {
		handlers.WriteLog(c, http.StatusBadRequest, "Validation errors encountered", "Invalid email or username or password", log)
	} else {
		handlers.WriteLog(c, http.StatusProcessing, "Validation passed, continuing", nil, log)
	}

	return user, isError
}

// AddAuthData adds a new user to the database.
func (a *AuthServiceImpl) Login(c *gin.Context) handlers.ServiceResponseWithLogging {
	log := helpers.CreateLog(c, a)
	var input dtos.InputLoginDTO

	if err := c.ShouldBind(&input); err != nil {
		return handlers.ServiceResponseWithLogging{
			Status:  http.StatusBadRequest,
			Message: "Invalid Input",
			Data:    nil,
			Err:     err.Error(),
			Log:     log,
		}
	}

	// Validate input using golang validator
	if err := handlers.ValidateStruct(input); err != nil {
		errors := handlers.ValidationErrorHandlerV1(c, err, input)
		return handlers.ServiceResponseWithLogging{
			Status:  http.StatusBadRequest,
			Message: "Error Invalid Data",
			Data:    nil,
			Err:     errors,
			Log:     log,
		}
	}

	// Check and validate input that cannot be validate by golang validator
	user, err := a.inputValidator(input, c)
	if err {
		return handlers.ServiceResponseWithLogging{
			Status:  http.StatusBadRequest,
			Message: "Invalid email or password",
			Data:    nil,
			Err:     "Error in input validator",
			Log:     log,
		}
	}

	token, error := helpers.GenerateJWT(user)
	if error != nil {
		return handlers.ServiceResponseWithLogging{
			Status:  http.StatusBadRequest,
			Message: "Failed to generate token",
			Data:    nil,
			Err:     error.Error(),
			Log:     log,
		}
	}

	data := map[string]interface{}{
		"loginAt": time.Now(),
		"user":    user.Name,
		"role":    user.Role.Name,
		"token":   fmt.Sprintf("Bearer %s", token),
	}

	return handlers.ServiceResponseWithLogging{
		Status:  http.StatusOK,
		Message: "User Login Successfully",
		Data:    data,
		Err:     nil,
		Log:     log,
	}
}
