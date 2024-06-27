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
	"gorm.io/gorm"
)

// RoleService defines the methods for the role service.
type RoleService interface {
	GetAll(c *gin.Context) handlers.ServiceResponse
	GetByID(c *gin.Context) handlers.ServiceResponse
	AddData(c *gin.Context) handlers.ServiceResponse
	UpdateData(c *gin.Context) handlers.ServiceResponse
	DeleteData(c *gin.Context) handlers.ServiceResponse
}

// RoleServiceImpl is the implementation of the RoleService interface.
type RoleServiceImpl struct {
	db *gorm.DB
}

// NewRoleService creates a new instance of RoleServiceImpl.
func RoleServiceConstructor(db *gorm.DB) RoleService {
	return &RoleServiceImpl{db: db}
}

// Validate user input that validator cannot check for POST and PUT / PATCH method
// method parameter option are: ["POST", "PUT", "PATCH"]
func (r *RoleServiceImpl) inputValidator(model models.USR_Role, method string) (map[string]map[string]string, bool) {
	// Setup variable
	errors := map[string]map[string]string{"errors": {}}
	is_error := false
	var result *gorm.DB

	// Check name duplication
	var duplicateName models.USR_Role
	if method == "POST" { // Check for POST method
		result = r.db.Limit(1).Where("name = ?", model.Name).Find(&duplicateName)
	} else { // Check for PUT and PATCH method
		result = r.db.Limit(1).Where("name = ?", model.Name).Not("id = ?", model.ID).Find(&duplicateName)
	}
	if result.Error != nil || result.RowsAffected >= 1 {
		errors["errors"]["name"] = fmt.Sprintf("Role name %s already exist", model.Name)
		is_error = true
	}

	return errors, is_error
}

// GetAllRoles retrieves all roles from the database and returns them in a ServiceResponse.
func (r *RoleServiceImpl) GetAll(c *gin.Context) handlers.ServiceResponse {
	var roles []models.USR_Role
	var data interface{}

	query := r.db

	// Check if using pagination and order in query
	// Apply pagination if the relevant query parameters are present
	if c.Query("page") != "" || c.Query("limit") != "" {
		query = query.Scopes(helpers.Paginate(c))
	}

	// Apply ordering if the relevant query parameters are present
	if c.Query("order_by") != "" || c.Query("order") != "" {
		// allowedOrderedFields is whitelist of option that user could use to order,
		// this is also used to prevent sql injection on order
		allowedOrderFields := []string{"id", "name", "is_administrative", "created_at", "updated_at"}
		query = query.Scopes(helpers.Order(c, allowedOrderFields))
	}

	if err := query.Find(&roles).Error; err != nil {
		return handlers.ServiceResponse{
			Status:  http.StatusInternalServerError,
			Message: "Error Getting Data",
			Data:    nil,
			Err:     err,
		}
	}

	// Convert role to DTOs
	roleDTOs := dtos.ToUSRRoleMinimalDTOs(roles)
	data = roleDTOs

	// Setup data for paginated result
	if c.Query("page") != "" || c.Query("limit") != "" {
		var totalRows int64
		r.db.Model(&models.USR_Role{}).Count(&totalRows)
		data = helpers.GeneratePaginatedQuery(c, totalRows, dtos.MinimalRoleDTOToInterfaceSlice(roleDTOs))
	}

	return handlers.ServiceResponse{
		Status:  http.StatusOK,
		Message: "Success Getting All Roles Data",
		Data:    data,
		Err:     nil,
	}
}

// GetRoleByID retrieves a role by its ID and returns it in a ServiceResponse.
func (r *RoleServiceImpl) GetByID(c *gin.Context) handlers.ServiceResponse {
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

	var role models.USR_Role

	// Fetch the role from the database by ID with preloaded features and modules
	result := r.db.Preload("Features").Preload("Features.Module").Limit(1).Where("id = ?", id).Find(&role)
	if result.Error != nil || result.RowsAffected == 0 {
		return handlers.ServiceResponse{
			Status:  http.StatusNotFound,
			Message: "Role not found",
			Data:    nil,
			Err:     nil,
		}
	}

	// Convert role to DTO
	roleDTO := dtos.ToUSRRoleDTO(role)

	return handlers.ServiceResponse{
		Status:  http.StatusOK,
		Message: "Success Getting Role Data",
		Data:    roleDTO,
		Err:     nil,
	}
}

// AddRoleData adds a new role to the database.
func (r *RoleServiceImpl) AddData(c *gin.Context) handlers.ServiceResponse {
	var roleDTO dtos.InputUSRRoleDTO
	if err := c.ShouldBindJSON(&roleDTO); err != nil {
		return handlers.ServiceResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid Input",
			Data:    nil,
			Err:     err,
		}
	}

	role := dtos.InputToUSRRoleModel(roleDTO)

	// Validate input using golang validator
	if err := handlers.ValidateStruct(roleDTO); err != nil {
		errors := handlers.ValidationErrorHandlerV1(c, err, roleDTO)
		return handlers.ServiceResponse{
			Status:  http.StatusBadRequest,
			Message: "Error Invalid Data",
			Data:    nil,
			Err:     errors,
		}
	}

	// Check and validate input that cannot be validate by golang validator
	errors, errorHappen := r.inputValidator(role, "POST")
	if errorHappen {
		return handlers.ServiceResponse{
			Status:  http.StatusBadRequest,
			Message: "Error Invalid Data",
			Data:    nil,
			Err:     errors,
		}
	}

	// Fetch features from the database
	var features []*models.USR_Feature
	if err := r.db.Where("id IN ?", roleDTO.Features).Find(&features).Error; err != nil {
		return handlers.ServiceResponse{
			Status:  http.StatusBadRequest,
			Message: "Error Invalid Features Data",
			Data:    nil,
			Err:     err,
		}
	}

	role.Features = features

	// Add the role to the database
	if err := r.db.Create(&role).Error; err != nil {
		return handlers.ServiceResponse{
			Status:  http.StatusInternalServerError,
			Message: "Error Creating Data",
			Data:    nil,
			Err:     err,
		}
	}

	return handlers.ServiceResponse{
		Status:  http.StatusCreated,
		Message: "Role Created Successfully",
		Data:    dtos.ToUSRRoleDTO(role),
		Err:     nil,
	}
}

// UpdateRoleData updates an existing role in the database.
func (r *RoleServiceImpl) UpdateData(c *gin.Context) handlers.ServiceResponse {
	var roleDTO dtos.InputUSRRoleDTO
	if err := c.ShouldBindJSON(&roleDTO); err != nil {
		return handlers.ServiceResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid Input",
			Data:    nil,
			Err:     err,
		}
	}

	var role models.USR_Role
	input := dtos.InputToUSRRoleModel(roleDTO)

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

	// Check Role Existence
	result := r.db.Preload("Features").Limit(1).Where("id = ?", id).Find(&role)
	if result.Error != nil || result.RowsAffected == 0 {
		return handlers.ServiceResponse{
			Status:  http.StatusNotFound,
			Message: "Role not found",
			Data:    nil,
			Err:     nil,
		}
	}

	// Parsing id params to input dto
	input.ID = uint(id)

	// Validate input using golang validator
	if err := handlers.ValidateStruct(roleDTO); err != nil {
		errors := handlers.ValidationErrorHandlerV1(c, err, roleDTO)
		return handlers.ServiceResponse{
			Status:  http.StatusBadRequest,
			Message: "Error Invalid Data",
			Data:    nil,
			Err:     errors,
		}
	}

	// Check and validate input that cannot be validate by golang validator
	errors, errorHappen := r.inputValidator(input, "PUT")
	if errorHappen {
		return handlers.ServiceResponse{
			Status:  http.StatusBadRequest,
			Message: "Error Invalid Data",
			Data:    nil,
			Err:     errors,
		}
	}

	// Fetch features from the database
	var features []*models.USR_Feature
	if err := r.db.Where("id IN ?", roleDTO.Features).Find(&features).Error; err != nil {
		return handlers.ServiceResponse{
			Status:  http.StatusBadRequest,
			Message: "Error Invalid Features Data",
			Data:    nil,
			Err:     err,
		}
	}

	// Update the role fields
	role.Name = roleDTO.Name
	role.IsAdministrative = roleDTO.IsAdministrative

	// Update role features
	// Set new features directly
	if err := r.db.Model(&role).Association("Features").Replace(features); err != nil {
		return handlers.ServiceResponse{
			Status:  http.StatusBadRequest,
			Message: "Error Updating Role Features Data",
			Data:    nil,
			Err:     err,
		}
	}

	// Save the updated role to the database
	if err := r.db.Save(&role).Error; err != nil {
		return handlers.ServiceResponse{
			Status:  http.StatusInternalServerError,
			Message: "Error Updating Data",
			Data:    nil,
			Err:     err,
		}
	}

	return handlers.ServiceResponse{
		Status:  http.StatusOK,
		Message: "Role Updated Successfully",
		Data:    dtos.ToUSRRoleDTO(role),
		Err:     nil,
	}
}

// DeleteRole deletes a role from the database.
func (r *RoleServiceImpl) DeleteData(c *gin.Context) handlers.ServiceResponse {
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

	// Check Role Existence
	var role models.USR_Role
	result := r.db.Limit(1).Where("id = ?", id).Find(&role)
	if result.Error != nil || result.RowsAffected == 0 {
		return handlers.ServiceResponse{
			Status:  http.StatusNotFound,
			Message: "Role not found",
			Data:    nil,
			Err:     nil,
		}
	}

	// Delete the role from the database
	if err := r.db.Delete(&models.USR_Role{}, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return handlers.ServiceResponse{
				Status:  http.StatusNotFound,
				Message: "Role Not Found",
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
		Message: "Role Deleted Successfully",
		Data:    nil,
		Err:     nil,
	}
}
