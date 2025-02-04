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
	GetAll(c *gin.Context) handlers.ServiceResponseWithLogging
	GetByID(c *gin.Context) handlers.ServiceResponseWithLogging
	AddData(c *gin.Context) handlers.ServiceResponseWithLogging
	UpdateData(c *gin.Context) handlers.ServiceResponseWithLogging
	DeleteData(c *gin.Context) handlers.ServiceResponseWithLogging
	ResetPass(c *gin.Context) handlers.ServiceResponseWithLogging
	ChangePass(c *gin.Context) handlers.ServiceResponseWithLogging
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
func (u *UserServiceImpl) inputValidator(model models.USR_User, method string, c *gin.Context) (map[string]map[string]string, bool) {
	// Setup variable
	errors := map[string]map[string]string{"errors": {}}
	is_error := false
	var result *gorm.DB

	// Create log
	log := helpers.CreateLog(c, u)

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

	if is_error {
		handlers.WriteLog(c, http.StatusBadRequest, "Validation errors encountered", errors, log)
	} else {
		handlers.WriteLog(c, http.StatusProcessing, "Validation passed, continuing", nil, log)
	}

	return errors, is_error
}

// GetAllUsers retrieves all users from the database and returns them in a ServiceResponseWithLogging.
func (u *UserServiceImpl) GetAll(c *gin.Context) handlers.ServiceResponseWithLogging {
	log := helpers.CreateLog(c, u)

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
		return handlers.ServiceResponseWithLogging{
			Status:  http.StatusInternalServerError,
			Message: "Error Getting Data",
			Data:    nil,
			Err:     err.Error(),
			Log:     log,
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

	return handlers.ServiceResponseWithLogging{
		Status:  http.StatusOK,
		Message: "Success Getting All Users Data",
		Data:    data,
		Err:     nil,
		Log:     log,
	}
}

// GetUserByID retrieves a user by its ID and returns it in a ServiceResponseWithLogging.
func (u *UserServiceImpl) GetByID(c *gin.Context) handlers.ServiceResponseWithLogging {
	log := helpers.CreateLog(c, u)

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return handlers.ServiceResponseWithLogging{
			Status:  http.StatusBadRequest,
			Message: "Invalid ID",
			Data:    nil,
			Err:     err.Error(),
			Log:     log,
		}
	}

	var user models.USR_User

	// Fetch the user from the database by ID
	result := u.db.Preload("Role").Preload("Role.Features").Preload("Role.Features.Module").Limit(1).Where("id = ?", id).Omit("password").Find(&user)
	if result.Error != nil || result.RowsAffected == 0 {
		return handlers.ServiceResponseWithLogging{
			Status:  http.StatusNotFound,
			Message: "User not found",
			Data:    nil,
			Err:     nil,
			Log:     log,
		}
	}

	// Convert user to DTO
	userDTO := dtos.ToUSRUserDTO(user)

	return handlers.ServiceResponseWithLogging{
		Status:  http.StatusOK,
		Message: "Success Getting User Data",
		Data:    userDTO,
		Err:     nil,
		Log:     log,
	}
}

// AddUserData adds a new user to the database.
func (u *UserServiceImpl) AddData(c *gin.Context) handlers.ServiceResponseWithLogging {
	log := helpers.CreateLog(c, u)

	var input dtos.CreateUSRUserInputDTO

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

	userModel := dtos.InputCreateToUSRUserModel(input)
	// Check and validate input that cannot be validate by golang validator
	errors, errorHappen := u.inputValidator(userModel, "POST", c)
	if errorHappen {
		return handlers.ServiceResponseWithLogging{
			Status:  http.StatusBadRequest,
			Message: "Error Invalid Data",
			Data:    errors,
			Err:     errors,
			Log:     log,
		}
	}

	// Add the user to the database
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userModel.Password), bcrypt.DefaultCost)
	if err != nil {
		return handlers.ServiceResponseWithLogging{
			Status:  http.StatusBadRequest,
			Message: "Failed to hash password",
			Data:    nil,
			Err:     err.Error(),
			Log:     log,
		}
	}

	userModel.Password = string(hashedPassword)

	if err := u.db.Create(&userModel).Error; err != nil {
		return handlers.ServiceResponseWithLogging{
			Status:  http.StatusInternalServerError,
			Message: "Error Creating Data",
			Data:    nil,
			Err:     err.Error(),
			Log:     log,
		}
	}

	input.ID = userModel.ID

	return handlers.ServiceResponseWithLogging{
		Status:  http.StatusCreated,
		Message: "User Created Successfully",
		Data:    input,
		Err:     nil,
		Log:     log,
	}
}

// UpdateUserData updates an existing user in the database.
func (u *UserServiceImpl) UpdateData(c *gin.Context) handlers.ServiceResponseWithLogging {
	log := helpers.CreateLog(c, u)

	// Check Params Validity
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return handlers.ServiceResponseWithLogging{
			Status:  http.StatusBadRequest,
			Message: "Invalid ID",
			Data:    nil,
			Err:     err.Error(),
			Log:     log,
		}
	}

	// Check if user can change password (it self or admin)
	var userPayload models.USR_User
	var rolePayload models.USR_Role
	helpers.GetUserPayload(c, &userPayload)
	helpers.GetRolePayload(c, &rolePayload)
	if int(userPayload.ID) != id && !rolePayload.IsAdministrative {
		return handlers.ServiceResponseWithLogging{
			Status:  http.StatusForbidden,
			Message: "Forbidden, unable to alter another user's data",
			Data:    nil,
			Err:     "Non admin user trying to update other user's data",
			Log:     log,
		}
	}

	var input dtos.UpdateUSRUserInputDTO

	if err := c.ShouldBind(&input); err != nil {
		return handlers.ServiceResponseWithLogging{
			Status:  http.StatusBadRequest,
			Message: "Invalid Input",
			Data:    nil,
			Err:     err.Error(),
			Log:     log,
		}
	}

	var user models.USR_User
	data := dtos.InputUpdateToUSRUserModel(input)
	input.ID = uint(id)
	data.ID = uint(id)

	// Check User Existence
	result := u.db.Limit(1).Where("id = ?", id).Find(&user)
	if result.Error != nil || result.RowsAffected == 0 {
		return handlers.ServiceResponseWithLogging{
			Status:  http.StatusNotFound,
			Message: "User not found",
			Data:    nil,
			Err:     nil,
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
	errors, errorHappen := u.inputValidator(data, "PUT", c)
	if errorHappen {
		return handlers.ServiceResponseWithLogging{
			Status:  http.StatusBadRequest,
			Message: "Error Invalid Data",
			Data:    nil,
			Err:     errors,
			Log:     log,
		}
	}

	// Update the user fields
	user.Username = data.Username
	user.Name = data.Name
	user.Email = data.Email
	user.RoleID = data.RoleID

	// Save the updated user to the database
	if err := u.db.Save(&user).Error; err != nil {
		return handlers.ServiceResponseWithLogging{
			Status:  http.StatusInternalServerError,
			Message: "Error Updating Data",
			Data:    nil,
			Err:     err.Error(),
			Log:     log,
		}
	}

	return handlers.ServiceResponseWithLogging{
		Status:  http.StatusOK,
		Message: "User Updated Successfully",
		Data:    input,
		Err:     nil,
		Log:     log,
	}
}

// DeleteUser deletes a user from the database.
func (u *UserServiceImpl) DeleteData(c *gin.Context) handlers.ServiceResponseWithLogging {
	log := helpers.CreateLog(c, u)

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil { // Check Params Validity
		return handlers.ServiceResponseWithLogging{
			Status:  http.StatusBadRequest,
			Message: "Invalid ID",
			Data:    nil,
			Err:     err.Error(),
			Log:     log,
		}
	}

	// Check if user can change password (it self or admin)
	var userPayload models.USR_User
	var rolePayload models.USR_Role
	helpers.GetUserPayload(c, &userPayload)
	helpers.GetRolePayload(c, &rolePayload)
	if int(userPayload.ID) != id && !rolePayload.IsAdministrative {
		return handlers.ServiceResponseWithLogging{
			Status:  http.StatusForbidden,
			Message: "Forbidden, unable to delete another user",
			Data:    nil,
			Err:     "Non admin user trying to delete other user",
			Log:     log,
		}
	}

	// Check User Existence
	var user models.USR_User
	result := u.db.Limit(1).Where("id = ?", id).Find(&user)
	if result.Error != nil || result.RowsAffected == 0 {
		return handlers.ServiceResponseWithLogging{
			Status:  http.StatusNotFound,
			Message: "User not found",
			Data:    nil,
			Err:     nil,
			Log:     log,
		}
	}

	// Delete the user from the database
	if err := u.db.Delete(&models.USR_User{}, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return handlers.ServiceResponseWithLogging{
				Status:  http.StatusNotFound,
				Message: "User Not Found",
				Data:    nil,
				Err:     err.Error(),
				Log:     log,
			}
		}
		return handlers.ServiceResponseWithLogging{
			Status:  http.StatusInternalServerError,
			Message: "Error Deleting Data",
			Data:    nil,
			Err:     err.Error(),
			Log:     log,
		}
	}

	return handlers.ServiceResponseWithLogging{
		Status:  http.StatusOK,
		Message: "User Deleted Successfully",
		Data:    nil,
		Err:     nil,
		Log:     log,
	}
}

// ResetPass reset user password data.
func (u *UserServiceImpl) ResetPass(c *gin.Context) handlers.ServiceResponseWithLogging {
	log := helpers.CreateLog(c, u)
	var input dtos.ResetPassUSRUserInputDTO

	if err := c.ShouldBind(&input); err != nil {
		return handlers.ServiceResponseWithLogging{
			Status:  http.StatusBadRequest,
			Message: "Invalid Input",
			Data:    nil,
			Err:     err.Error(),
			Log:     log,
		}
	}

	// Check Params Validity
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return handlers.ServiceResponseWithLogging{
			Status:  http.StatusBadRequest,
			Message: "Invalid ID",
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

	// Check User Existence
	var user models.USR_User
	result := u.db.Limit(1).Where("id = ?", id).Find(&user)
	if result.Error != nil || result.RowsAffected == 0 {
		return handlers.ServiceResponseWithLogging{
			Status:  http.StatusNotFound,
			Message: "User not found",
			Data:    nil,
			Err:     nil,
			Log:     log,
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
		return handlers.ServiceResponseWithLogging{
			Status:  http.StatusBadRequest,
			Message: "Invalid Data",
			Data:    nil,
			Err:     errors,
			Log:     log,
		}
	}

	// Add the user to the database
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return handlers.ServiceResponseWithLogging{
			Status:  http.StatusBadRequest,
			Message: "Failed to hash password",
			Data:    nil,
			Err:     err.Error(),
			Log:     log,
		}
	}

	user.Password = string(hashedPassword)

	// Save the new user password to the database
	if err := u.db.Save(&user).Error; err != nil {
		return handlers.ServiceResponseWithLogging{
			Status:  http.StatusInternalServerError,
			Message: "Error Updating Data",
			Data:    nil,
			Err:     err.Error(),
			Log:     log,
		}
	}

	return handlers.ServiceResponseWithLogging{
		Status:  http.StatusOK,
		Message: "User Password Reset Successfully",
		Data:    nil,
		Err:     nil,
		Log:     log,
	}
}

// ChangePass change user password data.
func (u *UserServiceImpl) ChangePass(c *gin.Context) handlers.ServiceResponseWithLogging {
	log := helpers.CreateLog(c, u)

	// Check Params Validity
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return handlers.ServiceResponseWithLogging{
			Status:  http.StatusBadRequest,
			Message: "Invalid ID",
			Data:    nil,
			Err:     err.Error(),
			Log:     log,
		}
	}

	// Check if user can change password (it self or admin)
	var userPayload models.USR_User
	var rolePayload models.USR_Role
	helpers.GetUserPayload(c, &userPayload)
	helpers.GetRolePayload(c, &rolePayload)
	if int(userPayload.ID) != id && !rolePayload.IsAdministrative {
		return handlers.ServiceResponseWithLogging{
			Status:  http.StatusForbidden,
			Message: "Forbidden, unable to alter another user's password",
			Data:    nil,
			Err:     "Non admin user trying to change other user's password",
			Log:     log,
		}
	}

	var input dtos.ChangePassUSRUserInputDTO

	if err := c.ShouldBind(&input); err != nil {
		return handlers.ServiceResponseWithLogging{
			Status:  http.StatusBadRequest,
			Message: "Invalid Input",
			Data:    nil,
			Err:     err,
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

	// Check User Existence
	var user models.USR_User
	result := u.db.Limit(1).Where("id = ?", id).Find(&user)
	if result.Error != nil || result.RowsAffected == 0 {
		return handlers.ServiceResponseWithLogging{
			Status:  http.StatusNotFound,
			Message: "User not found",
			Data:    nil,
			Err:     nil,
			Log:     log,
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
		return handlers.ServiceResponseWithLogging{
			Status:  http.StatusBadRequest,
			Message: "Invalid Data",
			Data:    nil,
			Err:     errors,
			Log:     log,
		}
	}

	// Add the user to the database
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return handlers.ServiceResponseWithLogging{
			Status:  http.StatusBadRequest,
			Message: "Failed to hash password",
			Data:    nil,
			Err:     err.Error(),
			Log:     log,
		}
	}

	user.Password = string(hashedPassword)

	// Save the new user password to the database
	if err := u.db.Save(&user).Error; err != nil {
		return handlers.ServiceResponseWithLogging{
			Status:  http.StatusInternalServerError,
			Message: "Error Updating Data",
			Data:    nil,
			Err:     err.Error(),
			Log:     log,
		}
	}

	return handlers.ServiceResponseWithLogging{
		Status:  http.StatusOK,
		Message: "User Password Reset Successfully",
		Data:    nil,
		Err:     nil,
		Log:     log,
	}
}
