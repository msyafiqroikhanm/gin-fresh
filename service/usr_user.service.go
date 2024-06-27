package service

import (
	"fmt"
	"jxb-eprocurement/handlers"
	"jxb-eprocurement/handlers/dtos"
	"jxb-eprocurement/helpers"
	"jxb-eprocurement/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserService defines the methods for the user service.
type UserService interface {
	GetAll(c *gin.Context) handlers.ServiceResponse
	GetByID(c *gin.Context) handlers.ServiceResponse
	AddData(c *gin.Context) handlers.ServiceResponse
	UpdateData(c *gin.Context) handlers.ServiceResponse
	DeleteData(c *gin.Context) handlers.ServiceResponse
	ResetPass(c *gin.Context) handlers.ServiceResponse
	ChangePass(c *gin.Context) handlers.ServiceResponse
}

// UserServiceImpl is the implementation of the UserService interface.
type UserServiceImpl struct {
	db *gorm.DB
}

// NewUserService creates a new instance of UserServiceImpl.
func UserServiceConstructor(db *gorm.DB) UserService {
	return &UserServiceImpl{db: db}
}

// Validate user input that validator cannot check for POST and PUT / PATCH method
// method parameter option are: ["POST", "PUT", "PATCH"]
func (u *UserServiceImpl) inputValidator(model models.USR_User, method string) (map[string]map[string]string, bool) {
	// Setup variable
	errors := map[string]map[string]string{"errors": {}}
	is_error := false
	var result *gorm.DB

	// Check email duplication
	var duplicateEmail models.USR_User
	if method == "POST" { // Check for POST method
		result = u.db.Limit(1).Where("email = ?", model.Email).Find(&duplicateEmail)
	} else { // Check for PUT and PATCH method
		result = u.db.Limit(1).Where("email = ?", model.Email).Not("id = ?", model.ID).Find(&duplicateEmail)
	}
	if result.Error != nil || result.RowsAffected >= 1 {
		errors["errors"]["email"] = fmt.Sprintf("User email %s already exist", model.Email)
		is_error = true
	}

	// Check if role exist
	var role models.USR_Role
	result = u.db.Table("usr_roles").Where("id = ?", model.RoleID).Limit(1).Find(&role)
	if result.Error != nil || result.RowsAffected == 0 {
		errors["errors"]["role_id"] = fmt.Sprintf("Role with id %d not found", model.RoleID)
		is_error = true
	}

	return errors, is_error
}

// GetAllUsers retrieves all users from the database and returns them in a ServiceResponse.
func (u *UserServiceImpl) GetAll(c *gin.Context) handlers.ServiceResponse {
	var users []models.USR_User
	var data interface{}

	query := u.db.Preload("Role", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "name").Unscoped()
	})

	// Check if using pagination and order in query
	// Apply pagination if the relevant query parameters are present
	if c.Query("page") != "" || c.Query("limit") != "" {
		query = query.Scopes(helpers.Paginate(c))
	}

	// Apply ordering if the relevant query parameters are present
	if c.Query("order_by") != "" || c.Query("order") != "" {
		// allowedOrderedFields is whitelist of option that user could use to order,
		// this is also used to prevent sql injection on order
		allowedOrderFields := []string{"id", "name", "email", "role_id", "created_at", "updated_at"}
		query = query.Scopes(helpers.Order(c, allowedOrderFields))
	}

	if err := query.Find(&users).Error; err != nil {
		return handlers.ServiceResponse{
			Status:  http.StatusInternalServerError,
			Message: "Error Getting Data",
			Data:    nil,
			Err:     err,
		}
	}

	// Convert user to DTOs
	userDTOs := dtos.ToUSRUserMinimalDTOs(users)
	data = userDTOs

	// Setup data for paginated result
	if c.Query("page") != "" || c.Query("limit") != "" {
		var totalRows int64
		u.db.Model(&models.USR_User{}).Count(&totalRows)
		data = helpers.GeneratePaginatedQuery(c, totalRows, dtos.MinimalUserDTOToInterfaceSlice(userDTOs))
	}

	return handlers.ServiceResponse{
		Status:  http.StatusOK,
		Message: "Success Getting All Users Data",
		Data:    data,
		Err:     nil,
	}
}

// GetUserByID retrieves a user by its ID and returns it in a ServiceResponse.
func (u *UserServiceImpl) GetByID(c *gin.Context) handlers.ServiceResponse {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return handlers.ServiceResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid ID",
			Data:    nil,
			Err:     nil,
		}
	}

	var user models.USR_User

	// Fetch the user from the database by ID
	result := u.db.Preload("Role").Preload("Role.Features").Preload("Role.Features.Module").Limit(1).Where("id = ?", id).Omit("password").Find(&user)
	if result.Error != nil || result.RowsAffected == 0 {
		return handlers.ServiceResponse{
			Status:  http.StatusNotFound,
			Message: "Role not found",
			Data:    nil,
			Err:     nil,
		}
	}

	// Convert user to DTO
	userDTO := dtos.ToUSRUserDTO(user)

	return handlers.ServiceResponse{
		Status:  http.StatusOK,
		Message: "Success Getting User Data",
		Data:    userDTO,
		Err:     nil,
	}
}

// AddUserData adds a new user to the database.
func (u *UserServiceImpl) AddData(c *gin.Context) handlers.ServiceResponse {
	var input dtos.CreateUSRUserInputDTO

	if err := c.ShouldBind(&input); err != nil {
		return handlers.ServiceResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid Input",
			Data:    nil,
			Err:     err,
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

	userModel := dtos.InputCreateToUSRUserModel(input)
	// Check and validate input that cannot be validate by golang validator
	errors, errorHappen := u.inputValidator(userModel, "POST")
	if errorHappen {
		return handlers.ServiceResponse{
			Status:  http.StatusBadRequest,
			Message: "Error Invalid Data",
			Data:    nil,
			Err:     errors,
		}
	}

	// Add the user to the database
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userModel.Password), bcrypt.DefaultCost)
	if err != nil {
		return handlers.ServiceResponse{
			Status:  http.StatusBadRequest,
			Message: "Failed to hash password",
			Data:    nil,
			Err:     err,
		}
	}

	userModel.Password = string(hashedPassword)

	if err := u.db.Create(&userModel).Error; err != nil {
		return handlers.ServiceResponse{
			Status:  http.StatusInternalServerError,
			Message: "Error Creating Data",
			Data:    nil,
			Err:     err,
		}
	}

	input.ID = userModel.ID

	return handlers.ServiceResponse{
		Status:  http.StatusCreated,
		Message: "User Created Successfully",
		Data:    input,
		Err:     nil,
	}
}

// UpdateUserData updates an existing user in the database.
func (u *UserServiceImpl) UpdateData(c *gin.Context) handlers.ServiceResponse {
	var input dtos.UpdateUSRUserInputDTO

	if err := c.ShouldBind(&input); err != nil {
		return handlers.ServiceResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid Input",
			Data:    nil,
			Err:     err,
		}
	}

	// Check Params Validity
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return handlers.ServiceResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid ID",
			Data:    nil,
			Err:     nil,
		}
	}

	var user models.USR_User
	data := dtos.InputUpdateToUSRUserModel(input)
	input.ID = uint(id)
	data.ID = uint(id)

	// Check User Existence
	result := u.db.Limit(1).Where("id = ?", id).Find(&user)
	if result.Error != nil || result.RowsAffected == 0 {
		return handlers.ServiceResponse{
			Status:  http.StatusNotFound,
			Message: "User not found",
			Data:    nil,
			Err:     nil,
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
	errors, errorHappen := u.inputValidator(data, "PUT")
	if errorHappen {
		return handlers.ServiceResponse{
			Status:  http.StatusBadRequest,
			Message: "Error Invalid Data",
			Data:    nil,
			Err:     errors,
		}
	}

	// Update the user fields
	user.Name = data.Name
	user.Email = data.Email
	user.RoleID = data.RoleID

	// Save the updated user to the database
	if err := u.db.Save(&user).Error; err != nil {
		return handlers.ServiceResponse{
			Status:  http.StatusInternalServerError,
			Message: "Error Updating Data",
			Data:    nil,
			Err:     err,
		}
	}

	return handlers.ServiceResponse{
		Status:  http.StatusOK,
		Message: "User Updated Successfully",
		Data:    input,
		Err:     nil,
	}
}

// DeleteUser deletes a user from the database.
func (u *UserServiceImpl) DeleteData(c *gin.Context) handlers.ServiceResponse {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil { // Check Params Validity
		return handlers.ServiceResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid ID",
			Data:    nil,
			Err:     nil,
		}
	}

	// Check User Existence
	var user models.USR_User
	result := u.db.Limit(1).Where("id = ?", id).Find(&user)
	if result.Error != nil || result.RowsAffected == 0 {
		return handlers.ServiceResponse{
			Status:  http.StatusNotFound,
			Message: "User not found",
			Data:    nil,
			Err:     nil,
		}
	}

	// Delete the user from the database
	if err := u.db.Delete(&models.USR_User{}, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return handlers.ServiceResponse{
				Status:  http.StatusNotFound,
				Message: "User Not Found",
				Data:    nil,
				Err:     err,
			}
		}
		return handlers.ServiceResponse{
			Status:  http.StatusInternalServerError,
			Message: "Error Deleting Data",
			Data:    nil,
			Err:     err,
		}
	}

	return handlers.ServiceResponse{
		Status:  http.StatusOK,
		Message: "User Deleted Successfully",
		Data:    nil,
		Err:     nil,
	}
}

// ResetPass reset user password data.
func (u *UserServiceImpl) ResetPass(c *gin.Context) handlers.ServiceResponse {
	var input dtos.ResetPassUSRUserInputDTO

	if err := c.ShouldBind(&input); err != nil {
		return handlers.ServiceResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid Input",
			Data:    nil,
			Err:     err,
		}
	}

	// Check Params Validity
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return handlers.ServiceResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid ID",
			Data:    nil,
			Err:     nil,
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

	// Check User Existence
	var user models.USR_User
	result := u.db.Limit(1).Where("id = ?", id).Find(&user)
	if result.Error != nil || result.RowsAffected == 0 {
		return handlers.ServiceResponse{
			Status:  http.StatusNotFound,
			Message: "User not found",
			Data:    nil,
			Err:     nil,
		}
	}

	// Check change password input value
	errors := map[string]map[string]string{"errors": {}}

	// Check if password and re-password is identical
	if input.Password != input.RePassword {
		errors["errors"]["password"] = "Re-Password and Password are different"
		errors["errors"]["re_password"] = "Re-Password and Password are different"
	}

	if len(errors["errors"]) != 0 {
		return handlers.ServiceResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid Data",
			Data:    nil,
			Err:     errors,
		}
	}

	// Add the user to the database
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return handlers.ServiceResponse{
			Status:  http.StatusBadRequest,
			Message: "Failed to hash password",
			Data:    nil,
			Err:     err,
		}
	}

	user.Password = string(hashedPassword)

	// Save the new user password to the database
	if err := u.db.Save(&user).Error; err != nil {
		return handlers.ServiceResponse{
			Status:  http.StatusInternalServerError,
			Message: "Error Updating Data",
			Data:    nil,
			Err:     err,
		}
	}

	return handlers.ServiceResponse{
		Status:  http.StatusOK,
		Message: "User Password Reset Successfully",
		Data:    nil,
		Err:     nil,
	}
}

// ChangePass change user password data.
func (u *UserServiceImpl) ChangePass(c *gin.Context) handlers.ServiceResponse {
	var input dtos.ChangePassUSRUserInputDTO

	if err := c.ShouldBind(&input); err != nil {
		return handlers.ServiceResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid Input",
			Data:    nil,
			Err:     err,
		}
	}

	// Check Params Validity
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return handlers.ServiceResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid ID",
			Data:    nil,
			Err:     nil,
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

	// Check User Existence
	var user models.USR_User
	result := u.db.Limit(1).Where("id = ?", id).Find(&user)
	if result.Error != nil || result.RowsAffected == 0 {
		return handlers.ServiceResponse{
			Status:  http.StatusNotFound,
			Message: "User not found",
			Data:    nil,
			Err:     nil,
		}
	}

	// Check change password input value
	errors := map[string]map[string]string{"errors": {}}

	// Check if password and re-password is identical
	if input.Password != input.RePassword {
		errors["errors"]["password"] = "Re-Password and Password are different"
		errors["errors"]["re_password"] = "Re-Password and Password are different"
	}

	// Check if old password is correct
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.OldPassword)); err != nil {
		errors["errors"]["old_password"] = "The old password is incorrect"
	}

	if len(errors["errors"]) != 0 {
		return handlers.ServiceResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid Data",
			Data:    nil,
			Err:     errors,
		}
	}

	// Add the user to the database
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return handlers.ServiceResponse{
			Status:  http.StatusBadRequest,
			Message: "Failed to hash password",
			Data:    nil,
			Err:     err,
		}
	}

	user.Password = string(hashedPassword)

	// Save the new user password to the database
	if err := u.db.Save(&user).Error; err != nil {
		return handlers.ServiceResponse{
			Status:  http.StatusInternalServerError,
			Message: "Error Updating Data",
			Data:    nil,
			Err:     err,
		}
	}

	return handlers.ServiceResponse{
		Status:  http.StatusOK,
		Message: "User Password Reset Successfully",
		Data:    nil,
		Err:     nil,
	}
}
