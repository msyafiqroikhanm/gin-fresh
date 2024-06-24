package service

import (
	"fmt"
	"jxb-eprocurement/handlers"
	"jxb-eprocurement/handlers/dtos"
	"jxb-eprocurement/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// FeatureService defines the methods for the module service.
type FeatureService interface {
	GetAll(c *gin.Context) handlers.ServiceResponse
	GetByID(c *gin.Context) handlers.ServiceResponse
	AddData(c *gin.Context) handlers.ServiceResponse
	UpdateData(c *gin.Context) handlers.ServiceResponse
	DeleteData(c *gin.Context) handlers.ServiceResponse
}

// FeatureServiceImpl is the implementation of the FeatureService interface.
type FeatureServiceImpl struct {
	db *gorm.DB
}

// NewFeatureService creates a new instance of FeatureServiceImpl.
func FeatureServiceConstructor(db *gorm.DB) FeatureService {
	return &FeatureServiceImpl{db: db}
}

// Validate user input that validator cannot check for POST and PUT / PATCH method
// method parameter option are: ["POST", "PUT", "PATCH"]
func (m *FeatureServiceImpl) inputValidator(feature models.USR_Feature, method string) (map[string]map[string]string, bool) {
	// Setup variable
	errors := map[string]map[string]string{"errors": {}}
	is_error := false
	var result *gorm.DB

	// Check name duplication
	var duplicateName models.USR_Feature
	if method == "POST" { // Check for POST method
		result = m.db.Limit(1).Where("name = ?", feature.Name).Find(&duplicateName)
	} else { // Check for PUT and PATCH method
		result = m.db.Limit(1).Where("name = ?", feature.Name).Not("id = ?", feature.ID).Find(&duplicateName)
	}
	if result.Error != nil || result.RowsAffected >= 1 {
		errors["errors"]["name"] = fmt.Sprintf("Module name %s already exist", feature.Name)
		is_error = true
	}

	// Check parent_id input validity
	if feature.ModuleID != 0 {
		var moduleInstance models.USR_Feature
		result = m.db.Limit(1).Where("id = ?", feature.ModuleID).Find(&moduleInstance)
		if result.Error != nil || result.RowsAffected == 0 {
			errors["errors"]["module_id"] = "Module Not Found"
			is_error = true
		}
	}

	return errors, is_error
}

// GetAllModules retrieves all modules from the database and returns them in a ServiceResponse.
func (m *FeatureServiceImpl) GetAll(c *gin.Context) handlers.ServiceResponse {
	var modules []models.USR_Module

	// Fetch all modules from the database
	if err := m.db.Preload("Child").Find(&modules).Error; err != nil {
		return handlers.ServiceResponse{
			Status:  http.StatusInternalServerError,
			Message: "Error Getting Data",
			Data:    nil,
			Err:     err,
		}
	}

	// Convert modules to DTOs
	moduleDTOs := dtos.ToUSRModuleMinimalDTOs(modules)

	return handlers.ServiceResponse{
		Status:  http.StatusOK,
		Message: "Success Getting All Modules Data",
		Data:    moduleDTOs,
		Err:     nil,
	}
}

// GetModuleByID retrieves a module by its ID and returns it in a ServiceResponse.
func (m *FeatureServiceImpl) GetByID(c *gin.Context) handlers.ServiceResponse {
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

	var module models.USR_Module

	// Fetch the module from the database by ID
	if err := m.db.Preload("Child").First(&module, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return handlers.ServiceResponse{
				Status:  http.StatusNotFound,
				Message: "Module Not Found",
				Data:    nil,
				Err:     err,
			}
		}
		return handlers.ServiceResponse{
			Status:  http.StatusInternalServerError,
			Message: "Error Getting Data",
			Data:    nil,
			Err:     err,
		}
	}

	// Convert module to DTO
	moduleDTO := dtos.ToUSRModuleDTO(module)

	return handlers.ServiceResponse{
		Status:  http.StatusOK,
		Message: "Success Getting Module Data",
		Data:    moduleDTO,
		Err:     nil,
	}
}

// AddfeatureData adds a new feature to the database.
func (m *FeatureServiceImpl) AddData(c *gin.Context) handlers.ServiceResponse {
	var input dtos.USRFeatureMinimalDTO

	if err := c.ShouldBind(&input); err != nil {
		return handlers.ServiceResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid Input",
			Data:    nil,
			Err:     err,
		}
	}

	// Validate input using golang validator with custom validations
	if err := handlers.ValidateStruct(input); err != nil {
		errors := handlers.ValidationErrorHandlerV1(c, err, input)
		return handlers.ServiceResponse{
			Status:  http.StatusBadRequest,
			Message: "Error Invalid Data",
			Data:    nil,
			Err:     errors,
		}
	}

	feature := dtos.ToUSRFeatureMinimalModel(input)
	// Check and validate input that cannot be validate by golang validator
	errors, errorHappen := m.inputValidator(feature, "POST")
	if errorHappen {
		return handlers.ServiceResponse{
			Status:  http.StatusBadRequest,
			Message: "Error Invalid Data",
			Data:    nil,
			Err:     errors,
		}
	}

	// Add the feature to the database
	if err := m.db.Create(&feature).Error; err != nil {
		return handlers.ServiceResponse{
			Status:  http.StatusInternalServerError,
			Message: "Error Creating Data",
			Data:    nil,
			Err:     err,
		}
	}

	return handlers.ServiceResponse{
		Status:  http.StatusCreated,
		Message: "Module Created Successfully",
		Data:    dtos.ToUSRFeatureMinimalDTO(feature),
		Err:     nil,
	}
}

// UpdateModuleData updates an existing module in the database.
func (m *FeatureServiceImpl) UpdateData(c *gin.Context) handlers.ServiceResponse {
	// var moduleDTO dtos.USRModuleMinimalDTO
	// if err := c.ShouldBind(&moduleDTO); err != nil { // Binding body data to moduleDTO
	// 	return handlers.ServiceResponse{
	// 		Status:  http.StatusBadRequest,
	// 		Message: "Invalid ID",
	// 		Data:    nil,
	// 		Err:     nil,
	// 	}
	// }

	// var module models.USR_Module
	// input := dtos.ToUSRModuleMinimalModel(moduleDTO)

	// // Check Params Validity
	// idStr := c.Param("id")
	// id, err := strconv.Atoi(idStr)
	// if err != nil {
	// 	return handlers.ServiceResponse{
	// 		Status:  http.StatusBadRequest,
	// 		Message: "Invalid ID",
	// 		Data:    nil,
	// 		Err:     nil,
	// 	}
	// }

	// // Check Module Existence
	// result := m.db.Limit(1).Where("id = ?", id).Find(&module)
	// if result.Error != nil || result.RowsAffected == 0 {
	// 	return handlers.ServiceResponse{
	// 		Status:  http.StatusNotFound,
	// 		Message: "Module not found",
	// 		Data:    nil,
	// 		Err:     nil,
	// 	}
	// }

	// // Parsing id params to input dto
	// input.ID = uint(id)

	// // Validate input using golang validator
	// if err := handlers.ValidateStruct(moduleDTO); err != nil {
	// 	errors := handlers.ValidationErrorHandlerV1(c, err, moduleDTO)
	// 	return handlers.ServiceResponse{
	// 		Status:  http.StatusBadRequest,
	// 		Message: "Error Invalid Data",
	// 		Data:    nil,
	// 		Err:     errors,
	// 	}
	// }

	// // Check and validate input that cannot be validate by golang validator
	// errors, errorHappen := m.inputValidator(input, "PUT")
	// if errorHappen {
	// 	return handlers.ServiceResponse{
	// 		Status:  http.StatusBadRequest,
	// 		Message: "Error Invalid Data",
	// 		Data:    nil,
	// 		Err:     errors,
	// 	}
	// }

	// // Update the module fields
	// module.Name = moduleDTO.Name
	// module.ParentID = moduleDTO.ParentID

	// // Save the updated module to the database
	// if err := m.db.Save(&module).Error; err != nil {
	// 	return handlers.ServiceResponse{
	// 		Status:  http.StatusInternalServerError,
	// 		Message: "Error Updating Data",
	// 		Data:    nil,
	// 		Err:     err,
	// 	}
	// }

	return handlers.ServiceResponse{
		Status:  http.StatusOK,
		Message: "Module Updated Successfully",
		// Data:    dtos.ToUSRModuleMinimalDTO(module),
		// Err:     nil,
	}
}

// DeleteModule deletes a module from the database.
func (m *FeatureServiceImpl) DeleteData(c *gin.Context) handlers.ServiceResponse {
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

	// Check Module Existence
	var module models.USR_Module
	result := m.db.Limit(1).Where("id = ?", id).Find(&module)
	if result.Error != nil || result.RowsAffected == 0 {
		return handlers.ServiceResponse{
			Status:  http.StatusNotFound,
			Message: "Module not found",
			Data:    nil,
			Err:     nil,
		}
	}

	// Delete the module from the database
	if err := m.db.Delete(&models.USR_Module{}, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return handlers.ServiceResponse{
				Status:  http.StatusNotFound,
				Message: "Module Not Found",
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
		Message: "Module Deleted Successfully",
		Data:    nil,
		Err:     nil,
	}
}
