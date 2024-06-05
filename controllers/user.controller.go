package controllers

import (
	"fmt"
	"jxb-eprocurement/handlers"
	"jxb-eprocurement/handlers/dtos"
	"jxb-eprocurement/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserController struct {
	DB *gorm.DB
}

func (ctrl *UserController) CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		handlers.ResponseFormatter(c, http.StatusBadRequest, nil, "Invalid Input")
		return
	}

	if err := handlers.ValidateStruct(user); err != nil {
		handlers.ValidationErrorHandler(c, err)
		return
	}

	var checkEmailDuplicate models.User
	if err := ctrl.DB.Where("email = ?", user.Email).First(&checkEmailDuplicate).Error; err == nil {
		handlers.ResponseFormatter(c, http.StatusConflict, nil, fmt.Sprintf("User with email '%s' is already exist", user.Email))
		return
	}

	var role models.Role
	if err := ctrl.DB.First(&role, user.RoleID).Error; err != nil {
		handlers.ResponseFormatter(c, http.StatusBadRequest, nil, "Invalid role ID")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		handlers.ResponseFormatter(c, http.StatusInternalServerError, nil, "Failed to hash password")
		return
	}

	user.Password = string(hashedPassword)

	if err := ctrl.DB.Create(&user).Error; err != nil {
		handlers.ResponseFormatter(c, http.StatusInternalServerError, nil, "Internal Server Error")
		return
	}

	// Preload the Role for the created User
	if err := ctrl.DB.Preload("Role").First(&user, user.ID).Error; err != nil {
		handlers.ResponseFormatter(c, http.StatusInternalServerError, nil, "Internal Server Error")
		return
	}

	handlers.ResponseFormatter(c, http.StatusOK, nil, "User created successfully")
}

func (ctrl *UserController) GetAllUsers(c *gin.Context) {
	var users []models.User
	if err := ctrl.DB.Preload("Role").Find(&users).Error; err != nil {
		handlers.ResponseFormatter(c, http.StatusInternalServerError, nil, "Internal Server Error")
		return
	}

	var userDTOs []dtos.UserDTO
	for _, user := range users {
		userDTOs = append(userDTOs, dtos.ToUserDTO(user))
	}

	handlers.ResponseFormatter(c, http.StatusOK, userDTOs, "Users retrieved successfully")
}

func (ctrl *UserController) GetUser(c *gin.Context) {
	var user models.User
	id := c.Param("id")

	if err := ctrl.DB.Preload("Role").First(&user, id).Error; err != nil {
		handlers.ResponseFormatter(c, http.StatusNotFound, nil, "User not found")
		return
	}

	handlers.ResponseFormatter(c, http.StatusOK, dtos.ToUserDTO(user), "User retrieved successfully")
}

func (ctrl *UserController) UpdateUser(c *gin.Context) {
	var user models.User
	id := c.Param("id")

	// Fetch the existing user
	if err := ctrl.DB.First(&user, id).Error; err != nil {
		handlers.ResponseFormatter(c, http.StatusNotFound, nil, "User not found")
		return
	}

	var input models.User
	if err := c.ShouldBindJSON(&input); err != nil {
		handlers.ResponseFormatter(c, http.StatusBadRequest, nil, "Invalid Input")
		return
	}

	// Check if email is already used by another user
	var checkEmailDuplicate models.User
	if err := ctrl.DB.Where("email = ? AND id != ?", input.Email, id).First(&checkEmailDuplicate).Error; err == nil {
		handlers.ResponseFormatter(c, http.StatusConflict, nil, fmt.Sprintf("User with email '%s' is already exist", input.Email))
		return
	}

	// Validate the input
	if err := handlers.ValidateStruct(input); err != nil {
		handlers.ValidationErrorHandler(c, err)
		return
	}

	// Check if the role ID is valid
	var role models.Role
	if err := ctrl.DB.First(&role, input.RoleID).Error; err != nil {
		handlers.ResponseFormatter(c, http.StatusBadRequest, nil, "Invalid role ID")
		return
	}

	if input.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			handlers.ResponseFormatter(c, http.StatusInternalServerError, nil, "Failed to hash password")
			return
		}
		user.Password = string(hashedPassword)
	}

	// Update user fields
	user.Name = input.Name
	user.Email = input.Email
	user.RoleID = input.RoleID

	if err := ctrl.DB.Save(&user).Error; err != nil {
		handlers.ResponseFormatter(c, http.StatusInternalServerError, nil, "Internal Server Error")
		return
	}

	// Preload the Role for the created User
	if err := ctrl.DB.Preload("Role").First(&user, user.ID).Error; err != nil {
		handlers.ResponseFormatter(c, http.StatusInternalServerError, nil, "Internal Server Error")
		return
	}

	handlers.ResponseFormatter(c, http.StatusOK, nil, "User updated successfully")
}

func (ctrl *UserController) DeleteUser(c *gin.Context) {
	var user models.User
	id := c.Param("id")

	if err := ctrl.DB.First(&user, id).Error; err != nil {
		handlers.ResponseFormatter(c, http.StatusNotFound, nil, "User not found")
		return
	}

	// Preload the Role for the created User
	if err := ctrl.DB.Preload("Role").First(&user, user.ID).Error; err != nil {
		handlers.ResponseFormatter(c, http.StatusInternalServerError, nil, "Internal Server Error")
		return
	}

	if err := ctrl.DB.Delete(&user).Error; err != nil {
		handlers.ResponseFormatter(c, http.StatusInternalServerError, nil, "Internal Server Error")
		return
	}

	handlers.ResponseFormatter(c, http.StatusOK, nil, "User deleted successfully")
}
