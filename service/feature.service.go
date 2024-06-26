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

// FeatureService defines the methods for the feature service.
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
		var featureInstance models.USR_Module
		result = m.db.Limit(1).Where("id = ?", feature.ModuleID).Find(&featureInstance)
		if result.Error != nil || result.RowsAffected == 0 {
			errors["errors"]["feature_id"] = "Module Not Found"
			is_error = true
		}
	}

	return errors, is_error
}

// GetAllModules retrieves all features from the database and returns them in a ServiceResponse.
func (m *FeatureServiceImpl) GetAll(c *gin.Context) handlers.ServiceResponse {
	moduleIDStr := c.Query("module_id")

	query := m.db.Preload("Module")
	// Validate module_id and apply filter if present
	if moduleIDStr != "" {
		moduleID, err := strconv.ParseUint(moduleIDStr, 10, 64)
		if err != nil {
			return handlers.ServiceResponse{
				Status:  http.StatusBadRequest,
				Message: "Invalid module_id",
				Data:    nil,
				Err:     err,
			}
		}

		// Apply filter to the query
		query = query.Where("module_id = ?", moduleID)
	}

	var features []models.USR_Feature
	// Fetch all features from the database with the constructed query
	if err := query.Find(&features).Error; err != nil {
		return handlers.ServiceResponse{
			Status:  http.StatusInternalServerError,
			Message: "Error Getting Data",
			Data:    nil,
			Err:     err,
		}
	}

	// Convert features to DTOs
	featureDTOs := dtos.ToUSRFeatureMinimalWithModuleDTOs(features)

	// middlewares.LogSystem(identifier, "INFO", "controller/exampleHandler", "Request processed successfully", startTime, time.Now())

	return handlers.ServiceResponse{
		Status:  http.StatusOK,
		Message: "Success Getting All Feature Data",
		Data:    featureDTOs,
		Err:     nil,
	}
}

// GetModuleByID retrieves a feature by its ID and returns it in a ServiceResponse.
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

	var feature models.USR_Feature

	// Fetch the feature from the database by ID
	if err := m.db.Preload("Module").First(&feature, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return handlers.ServiceResponse{
				Status:  http.StatusNotFound,
				Message: "Feature Not Found",
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

	// Convert feature to DTO
	featureDTO := dtos.ToUSRFeatureMinimalWithModuleDTO(feature)

	return handlers.ServiceResponse{
		Status:  http.StatusOK,
		Message: "Success Getting Feature Data",
		Data:    featureDTO,
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

// UpdateModuleData updates an existing feature in the database.
func (m *FeatureServiceImpl) UpdateData(c *gin.Context) handlers.ServiceResponse {
	var featureDTO dtos.USRFeatureMinimalDTO
	if err := c.ShouldBind(&featureDTO); err != nil { // Binding body data to featureDTO
		return handlers.ServiceResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid ID",
			Data:    nil,
			Err:     nil,
		}
	}

	var feature models.USR_Feature
	input := dtos.ToUSRFeatureMinimalModel(featureDTO)

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

	// Check Feature Existence
	result := m.db.Limit(1).Where("id = ?", id).Find(&feature)
	if result.Error != nil || result.RowsAffected == 0 {
		return handlers.ServiceResponse{
			Status:  http.StatusNotFound,
			Message: "Feature not found",
			Data:    nil,
			Err:     nil,
		}
	}

	// Parsing id params to input dto
	input.ID = uint(id)

	// Validate input using golang validator
	if err := handlers.ValidateStruct(featureDTO); err != nil {
		errors := handlers.ValidationErrorHandlerV1(c, err, featureDTO)
		return handlers.ServiceResponse{
			Status:  http.StatusBadRequest,
			Message: "Error Invalid Data",
			Data:    nil,
			Err:     errors,
		}
	}

	// Check and validate input that cannot be validate by golang validator
	errors, errorHappen := m.inputValidator(input, "PUT")
	if errorHappen {
		return handlers.ServiceResponse{
			Status:  http.StatusBadRequest,
			Message: "Error Invalid Data",
			Data:    nil,
			Err:     errors,
		}
	}

	// Update the feature fields
	feature.Name = featureDTO.Name
	feature.ModuleID = featureDTO.ModuleID

	// Save the updated feature to the database
	if err := m.db.Save(&feature).Error; err != nil {
		return handlers.ServiceResponse{
			Status:  http.StatusInternalServerError,
			Message: "Error Updating Data",
			Data:    nil,
			Err:     err,
		}
	}

	return handlers.ServiceResponse{
		Status:  http.StatusOK,
		Message: "Feature Updated Successfully",
		// Data:    dtos.ToUSRModuleMinimalDTO(feature),
		// Err:     nil,
	}
}

// DeleteModule deletes a feature from the database.
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

	// Check Feature Existence
	var feature models.USR_Feature
	result := m.db.Limit(1).Where("id = ?", id).Find(&feature)
	if result.Error != nil || result.RowsAffected == 0 {
		return handlers.ServiceResponse{
			Status:  http.StatusNotFound,
			Message: "Feature not found",
			Data:    nil,
			Err:     nil,
		}
	}

	// Delete the feature from the database
	if err := m.db.Delete(&models.USR_Feature{}, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return handlers.ServiceResponse{
				Status:  http.StatusNotFound,
				Message: "Feature Not Found",
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
		Message: "Feature Deleted Successfully",
		Data:    nil,
		Err:     nil,
	}
}
