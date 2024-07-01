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

// ModuleService defines the methods for the module service.
type ModuleService interface {
	GetAll(c *gin.Context) handlers.ServiceResponseWithLogging
	GetByID(c *gin.Context) handlers.ServiceResponseWithLogging
	AddData(c *gin.Context) handlers.ServiceResponseWithLogging
	UpdateData(c *gin.Context) handlers.ServiceResponseWithLogging
	DeleteData(c *gin.Context) handlers.ServiceResponseWithLogging
}

// ModuleServiceImpl is the implementation of the ModuleService interface.
type ModuleServiceImpl struct {
	db *gorm.DB
}

// NewModuleService creates a new instance of ModuleServiceImpl.
func ModuleServiceConstructor(db *gorm.DB) ModuleService {
	return &ModuleServiceImpl{db: db}
}

// Validate user input that validator cannot check for POST and PUT / PATCH method
// method parameter option are: ["POST", "PUT", "PATCH"]
func (m *ModuleServiceImpl) inputValidator(model models.USR_Module, method string, c *gin.Context) (map[string]map[string]string, bool) {
	// Setup variable
	errors := map[string]map[string]string{"errors": {}}
	is_error := false
	var result *gorm.DB

	// Create log
	log := helpers.CreateLog(c, m)

	// Check name duplication
	var duplicateName models.USR_Module
	if method == "POST" { // Check for POST method
		result = m.db.Limit(1).Where("name = ?", model.Name).Find(&duplicateName)
	} else { // Check for PUT and PATCH method
		result = m.db.Limit(1).Where("name = ?", model.Name).Not("id = ?", model.ID).Find(&duplicateName)
	}
	if result.Error != nil || result.RowsAffected >= 1 {
		errors["errors"]["name"] = fmt.Sprintf("Module name %s already exist", model.Name)
		is_error = true
	}

	// Check parent_id input validity
	if model.ParentID != nil {
		var parentData models.USR_Module
		result = m.db.Limit(1).Where("id = ?", model.ParentID).Find(&parentData)
		if result.Error != nil || result.RowsAffected == 0 {
			errors["errors"]["parent_id"] = "Parent Module Not Found"
			is_error = true
		}
	}

	if is_error {
		handlers.WriteLog(c, http.StatusBadRequest, "Validation errors encountered", errors, log)
	} else {
		handlers.WriteLog(c, http.StatusProcessing, "Validation passed, continuing", nil, log)
	}

	return errors, is_error
}

// GetAllModules retrieves all modules from the database and returns them in a ServiceResponseWithLogging.
func (m *ModuleServiceImpl) GetAll(c *gin.Context) handlers.ServiceResponseWithLogging {
	log := helpers.CreateLog(c, m)

	var modules []models.USR_Module
	var data interface{}

	query := m.db.Preload("Child")

	// Check if using pagination and order in query
	// Apply pagination if the relevant query parameters are present
	if c.Query("page") != "" || c.Query("limit") != "" {
		query = query.Scopes(helpers.Paginate(c))
	}

	// Apply ordering if the relevant query parameters are present
	if c.Query("order_by") != "" || c.Query("order") != "" {
		// allowedOrderedFields is whitelist of option that user could use to order,
		// this is also used to prevent sql injection on order
		allowedOrderFields := []string{"id", "name", "created_at", "updated_at"}
		query = query.Scopes(helpers.Order(c, allowedOrderFields))
	}

	// Fetch all modules from the database
	if err := query.Find(&modules).Error; err != nil {
		return handlers.ServiceResponseWithLogging{
			Status:  http.StatusInternalServerError,
			Message: "Error Getting Data",
			Data:    nil,
			Err:     err.Error(),
			Log:     log,
		}
	}

	// Convert modules to DTOs
	moduleDTOs := dtos.ToUSRModuleMinimalDTOs(modules)
	data = moduleDTOs

	// Setup data for paginated result
	if c.Query("page") != "" || c.Query("limit") != "" {
		var totalRows int64
		m.db.Model(&models.USR_Module{}).Count(&totalRows)
		data = helpers.GeneratePaginatedQuery(c, totalRows, dtos.MinimalUSRModuleDTOToInterfaceSlice(moduleDTOs))
	}

	return handlers.ServiceResponseWithLogging{
		Status:  http.StatusOK,
		Message: "Success Getting All Modules Data",
		Data:    data,
		Err:     nil,
		Log:     log,
	}
}

// GetModuleByID retrieves a module by its ID and returns it in a ServiceResponseWithLogging.
func (m *ModuleServiceImpl) GetByID(c *gin.Context) handlers.ServiceResponseWithLogging {
	log := helpers.CreateLog(c, m)

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

	var module models.USR_Module

	// Fetch the module from the database by ID
	result := m.db.Preload("Child").Preload("Features").Limit(1).Where("id = ?", id).Find(&module)
	if result.Error != nil || result.RowsAffected == 0 {
		return handlers.ServiceResponseWithLogging{
			Status:  http.StatusNotFound,
			Message: "Module not found",
			Data:    nil,
			Err:     nil,
			Log:     log,
		}
	}

	// Convert module to DTO
	moduleDTO := dtos.ToUSRModuleWithFeaturesDTO(module)

	return handlers.ServiceResponseWithLogging{
		Status:  http.StatusOK,
		Message: "Success Getting Module Data",
		Data:    moduleDTO,
		Err:     nil,
		Log:     log,
	}
}

// AddModuleData adds a new module to the database.
func (m *ModuleServiceImpl) AddData(c *gin.Context) handlers.ServiceResponseWithLogging {
	log := helpers.CreateLog(c, m)

	var moduleDTO dtos.USRModuleMinimalDTO
	if err := c.ShouldBind(&moduleDTO); err != nil {
		return handlers.ServiceResponseWithLogging{
			Status:  http.StatusBadRequest,
			Message: "Invalid Input",
			Data:    nil,
			Err:     err.Error(),
			Log:     log,
		}
	}

	module := dtos.ToUSRModuleMinimalModel(moduleDTO)

	// Validate input using golang validator
	if err := handlers.ValidateStruct(moduleDTO); err != nil {
		errors := handlers.ValidationErrorHandlerV1(c, err, moduleDTO)
		return handlers.ServiceResponseWithLogging{
			Status:  http.StatusBadRequest,
			Message: "Error Invalid Data",
			Data:    nil,
			Err:     errors,
			Log:     log,
		}
	}

	// Check and validate input that cannot be validate by golang validator
	errors, errorHappen := m.inputValidator(module, "POST", c)
	if errorHappen {
		return handlers.ServiceResponseWithLogging{
			Status:  http.StatusBadRequest,
			Message: "Error Invalid Data",
			Data:    nil,
			Err:     errors,
			Log:     log,
		}
	}

	// Add the module to the database
	if err := m.db.Create(&module).Error; err != nil {
		return handlers.ServiceResponseWithLogging{
			Status:  http.StatusInternalServerError,
			Message: "Error Creating Data",
			Data:    nil,
			Err:     err,
			Log:     log,
		}
	}

	return handlers.ServiceResponseWithLogging{
		Status:  http.StatusCreated,
		Message: "Module Created Successfully",
		Data:    dtos.ToUSRModuleMinimalDTO(module),
		Err:     nil,
		Log:     log,
	}
}

// UpdateModuleData updates an existing module in the database.
func (m *ModuleServiceImpl) UpdateData(c *gin.Context) handlers.ServiceResponseWithLogging {
	log := helpers.CreateLog(c, m)

	var moduleDTO dtos.USRModuleMinimalDTO
	if err := c.ShouldBind(&moduleDTO); err != nil { // Binding body data to moduleDTO
		return handlers.ServiceResponseWithLogging{
			Status:  http.StatusBadRequest,
			Message: "Invalid ID",
			Data:    nil,
			Err:     err.Error(),
			Log:     log,
		}
	}

	var module models.USR_Module
	input := dtos.ToUSRModuleMinimalModel(moduleDTO)

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

	// Check Module Existence
	result := m.db.Limit(1).Where("id = ?", id).Find(&module)
	if result.Error != nil || result.RowsAffected == 0 {
		return handlers.ServiceResponseWithLogging{
			Status:  http.StatusNotFound,
			Message: "Module not found",
			Data:    nil,
			Err:     nil,
			Log:     log,
		}
	}

	// Parsing id params to input dto
	input.ID = uint(id)

	// Validate input using golang validator
	if err := handlers.ValidateStruct(moduleDTO); err != nil {
		errors := handlers.ValidationErrorHandlerV1(c, err, moduleDTO)
		return handlers.ServiceResponseWithLogging{
			Status:  http.StatusBadRequest,
			Message: "Error Invalid Data",
			Data:    nil,
			Err:     errors,
			Log:     log,
		}
	}

	// Check and validate input that cannot be validate by golang validator
	errors, errorHappen := m.inputValidator(input, "PUT", c)
	if errorHappen {
		return handlers.ServiceResponseWithLogging{
			Status:  http.StatusBadRequest,
			Message: "Error Invalid Data",
			Data:    nil,
			Err:     errors,
			Log:     log,
		}
	}

	// Update the module fields
	module.Name = moduleDTO.Name
	module.ParentID = moduleDTO.ParentID

	// Save the updated module to the database
	if err := m.db.Save(&module).Error; err != nil {
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
		Message: "Module Updated Successfully",
		Data:    dtos.ToUSRModuleMinimalDTO(module),
		Err:     nil,
		Log:     log,
	}
}

// DeleteModule deletes a module from the database.
func (m *ModuleServiceImpl) DeleteData(c *gin.Context) handlers.ServiceResponseWithLogging {
	log := helpers.CreateLog(c, m)

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

	// Check Module Existence
	var module models.USR_Module
	result := m.db.Limit(1).Where("id = ?", id).Find(&module)
	if result.Error != nil || result.RowsAffected == 0 {
		return handlers.ServiceResponseWithLogging{
			Status:  http.StatusNotFound,
			Message: "Module not found",
			Data:    nil,
			Err:     nil,
			Log:     log,
		}
	}

	// Delete the module from the database
	if err := m.db.Delete(&models.USR_Module{}, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return handlers.ServiceResponseWithLogging{
				Status:  http.StatusNotFound,
				Message: "Module Not Found",
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
		Message: "Module Deleted Successfully",
		Data:    nil,
		Err:     nil,
		Log:     log,
	}
}
