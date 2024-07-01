package dtos

import (
	"jxb-eprocurement/models"
	"strconv"
)

type (
	// USRUserDTO represents a Data Transfer Object for the USR_User model in detail format.
	// It includes only the fields necessary for data transfer and serialization.
	USRUserDTO struct {
		ID       uint       `json:"id"`               // Unique identifier of the user
		Username string     `json:"username"`         // Name of the user
		Name     string     `json:"name"`             // Name of the user
		Email    string     `json:"email"`            // Name of the user
		RoleID   uint       `json:"role_id" gorm:"-"` // Foreign Key To Role Table
		Role     USRRoleDTO `json:"role"`             // Role Data
	}

	// USRUserDTO represents a Data Transfer Object for the USR_User model in minimal format.
	// It includes only the fields necessary for data transfer and serialization.
	USRUserMinimalDTO struct {
		ID       uint   `json:"id" form:"id"`                                 // Unique identifier of the user
		Username string `json:"username" form:"username" validate:"required"` // Username of the user
		Name     string `json:"name" form:"name" validate:"required"`         // Name of the user
		Email    string `json:"email" form:"email" validate:"required"`       // Email of the user
		RoleName string `json:"role_name"`
		RoleID   uint   `json:"role_id" form:"role_id" validate:"required"`
	}

	CreateUSRUserInputDTO struct {
		ID       uint   `json:"id"`
		Username string `json:"username" form:"username" validate:"required,no_space,min=3,max=100"`
		Name     string `json:"name" form:"name" validate:"required,min=3,max=100"`
		Email    string `json:"email" form:"email" validate:"required,email"`
		Password string `json:"password" form:"password" validate:"required,min=6"`
		RoleID   string `json:"role_id" form:"role_id" validate:"required,numeric"`
	}

	UpdateUSRUserInputDTO struct {
		ID       uint   `json:"id"`
		Username string `json:"username" form:"username" validate:"required,no_space,min=3,max=100"`
		Name     string `json:"name" form:"name" validate:"required,min=3,max=100"`
		Email    string `json:"email" form:"email" validate:"required,email"`
		RoleID   string `json:"role_id" form:"role_id" validate:"required,numeric"`
	}

	ResetPassUSRUserInputDTO struct {
		Password   string `json:"password" form:"password" validate:"required,min=6"`
		RePassword string `json:"re_password" form:"re_password" validate:"required,min=6"`
	}

	ChangePassUSRUserInputDTO struct {
		Password    string `json:"password" form:"password" validate:"required,min=6"`
		RePassword  string `json:"re_password" form:"re_password" validate:"required,min=6"`
		OldPassword string `json:"old_password" form:"old_password" validate:"required,min=6"`
	}

	LogUserInfo struct {
		ID       string `json:"id"`
		Username string `json:"username"`
	}
)

// ToUSRUserDTO converts a USR_User model to a USRUserDTO in minimal format.
// Use this function to where detail information of role not needed.
func ToUSRUserMinimalDTO(user models.USR_User) USRUserMinimalDTO {
	// Return the DTO with converted fields
	return USRUserMinimalDTO{
		ID:       user.ID,
		Username: user.Username,
		Name:     user.Name,
		Email:    user.Email,
		RoleID:   user.RoleID,
		RoleName: user.Role.Name,
	}
}

// ToUSRUserDTO converts a USR_User model to a USRUserDTO in minimal format.
// Use this function to where detail information of role not needed.
func ToUSRUserMinimalDTOs(users []models.USR_User) []USRUserMinimalDTO {
	var userDTOs []USRUserMinimalDTO

	for _, user := range users {
		// Assigning role name to user
		userDTO := ToUSRUserMinimalDTO(user)
		userDTO.RoleName = user.Role.Name

		userDTOs = append(userDTOs, userDTO)
	}

	// Return the DTO with converted fields
	return userDTOs
}

func ToUSRUserDTO(user models.USR_User) USRUserDTO {
	// Return the DTO with converted fields
	return USRUserDTO{
		ID:       user.ID,
		Username: user.Username,
		Name:     user.Name,
		Email:    user.Email,
		RoleID:   user.RoleID,
		Role:     ToUSRRoleDTO(user.Role),
	}
}

// ToUSRUserModel converts a USRUserDTO to a USR_User model in minimal format.
// Use this function to where detail information of role not needed.
func ToUSRUserMinimalModel(dto USRUserMinimalDTO) models.USR_User {
	// Return the model with converted fields
	return models.USR_User{
		ID:     dto.ID,
		Name:   dto.Name,
		RoleID: dto.RoleID,
	}
}

// InputToUSRUserModel converts a serialization from InputUSRUserDTO to a USR_User model in detail format.
// Use this function to where the feature that role have is needed.
func InputCreateToUSRUserModel(dto CreateUSRUserInputDTO) models.USR_User {
	// convert string into uint
	roleId, _ := strconv.ParseUint(dto.RoleID, 10, 32)

	// Return the model with converted fields
	return models.USR_User{
		Username: dto.Username,
		Name:     dto.Name,
		Email:    dto.Email,
		Password: dto.Password,
		RoleID:   uint(roleId),
	}
}

// InputToUSRUserModel converts a serialization from InputUSRUserDTO to a USR_User model in detail format.
// Use this function to where the feature that role have is needed.
func InputUpdateToUSRUserModel(dto UpdateUSRUserInputDTO) models.USR_User {
	// convert string into uint
	roleId, _ := strconv.ParseUint(dto.RoleID, 10, 32)

	// Return the model with converted fields
	return models.USR_User{
		Username: dto.Username,
		Name:     dto.Name,
		Email:    dto.Email,
		RoleID:   uint(roleId),
	}
}

// Function to convert slice of USRUserMinimalDTO into slice of interface
// Used for generating pagination data
func MinimalUserDTOToInterfaceSlice(slice []USRUserMinimalDTO) []interface{} {
	interfaceSlice := make([]interface{}, len(slice))
	for i, v := range slice {
		interfaceSlice[i] = v
	}
	return interfaceSlice
}
