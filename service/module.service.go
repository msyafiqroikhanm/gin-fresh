package service

import (
	"jxb-eprocurement/handlers"
	"jxb-eprocurement/handlers/dtos"
	"jxb-eprocurement/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ModuleService defines the methods for the module service.
type ModuleService interface {
	GetAll(c *gin.Context) handlers.ServiceResponse
	GetByID(c *gin.Context, id uint) handlers.ServiceResponse
	AddData(c *gin.Context, module dtos.USRModuleMinimalDTO) handlers.ServiceResponse
	UpdateData(c *gin.Context, id uint, module dtos.USRModuleMinimalDTO) handlers.ServiceResponse
	DeleteData(c *gin.Context, id uint) handlers.ServiceResponse
}

// ModuleServiceImpl is the implementation of the ModuleService interface.
type ModuleServiceImpl struct {
	db *gorm.DB
}

// NewModuleService creates a new instance of ModuleServiceImpl.
func NewModuleService(db *gorm.DB) ModuleService {
	return &ModuleServiceImpl{db: db}
}

// GetAllModules retrieves all modules from the database and returns them in a ServiceResponse.
func (m *ModuleServiceImpl) GetAll(c *gin.Context) handlers.ServiceResponse {
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
func (m *ModuleServiceImpl) GetByID(c *gin.Context, id uint) handlers.ServiceResponse {
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

// AddModuleData adds a new module to the database.
func (m *ModuleServiceImpl) AddData(c *gin.Context, moduleDTO dtos.USRModuleMinimalDTO) handlers.ServiceResponse {
	module := dtos.ToUSRModuleMinimalModel(moduleDTO)

	// Add the module to the database
	if err := m.db.Create(&module).Error; err != nil {
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
		Data:    dtos.ToUSRModuleMinimalDTO(module),
		Err:     nil,
	}
}

// UpdateModuleData updates an existing module in the database.
func (m *ModuleServiceImpl) UpdateData(c *gin.Context, id uint, moduleDTO dtos.USRModuleMinimalDTO) handlers.ServiceResponse {
	var module models.USR_Module

	// Find the existing module by ID
	if err := m.db.First(&module, id).Error; err != nil {
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

	// Update the module fields
	module.Name = moduleDTO.Name
	module.ParentID = moduleDTO.ParentID

	// Save the updated module to the database
	if err := m.db.Save(&module).Error; err != nil {
		return handlers.ServiceResponse{
			Status:  http.StatusInternalServerError,
			Message: "Error Updating Data",
			Data:    nil,
			Err:     err,
		}
	}

	return handlers.ServiceResponse{
		Status:  http.StatusOK,
		Message: "Module Updated Successfully",
		Data:    dtos.ToUSRModuleMinimalDTO(module),
		Err:     nil,
	}
}

// DeleteModule deletes a module from the database.
func (m *ModuleServiceImpl) DeleteData(c *gin.Context, id uint) handlers.ServiceResponse {
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
