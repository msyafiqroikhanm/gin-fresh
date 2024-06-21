package controllers

import (
	"jxb-eprocurement/handlers"
	"jxb-eprocurement/service"

	"github.com/gin-gonic/gin"
)

type ModuleController interface {
	GetAllModules(c *gin.Context)
	GetModule(c *gin.Context)
	CreateModule(c *gin.Context)
	UpdateModule(c *gin.Context)
	DeleteModule(c *gin.Context)
}

// ModuleControllerImpl is the implementation of the ModuleController interface.
type ModuleControllerImpl struct {
	service service.ModuleService
}

// NewModuleController creates a new instance of ModuleControllerImpl.
func NewModuleController(service service.ModuleService) ModuleController {
	return &ModuleControllerImpl{service: service}
}

// GetAllModules handles the request to get all modules.
func (mc *ModuleControllerImpl) GetAllModules(c *gin.Context) {
	response := mc.service.GetAll(c)
	if response.Err != nil {
		handlers.ResponseFormatter(c, response.Status, response.Err, response.Message)
	} else {
		handlers.ResponseFormatter(c, response.Status, response.Data, response.Message)
	}
}

// GetModule handles the request to get a module by ID.
func (mc *ModuleControllerImpl) GetModule(c *gin.Context) {
	response := mc.service.GetByID(c)
	if response.Err != nil {
		handlers.ResponseFormatter(c, response.Status, response.Err, response.Message)
	} else {
		handlers.ResponseFormatter(c, response.Status, response.Data, response.Message)
	}
}

// CreateModule handles the request to add a new module.
func (mc *ModuleControllerImpl) CreateModule(c *gin.Context) {
	response := mc.service.AddData(c)
	if response.Err != nil {
		handlers.ResponseFormatter(c, response.Status, response.Err, response.Message)
	} else {
		handlers.ResponseFormatter(c, response.Status, response.Data, response.Message)
	}
}

// UpdateModule handles the request to update a module.
func (mc *ModuleControllerImpl) UpdateModule(c *gin.Context) {
	response := mc.service.UpdateData(c)
	if response.Err != nil {
		handlers.ResponseFormatter(c, response.Status, response.Err, response.Message)
	} else {
		handlers.ResponseFormatter(c, response.Status, response.Data, response.Message)
	}
}

// DeleteModule handles the request to delete a module.
func (mc *ModuleControllerImpl) DeleteModule(c *gin.Context) {
	response := mc.service.DeleteData(c)
	if response.Err != nil {
		handlers.ResponseFormatter(c, response.Status, response.Err, response.Message)
	} else {
		handlers.ResponseFormatter(c, response.Status, response.Data, response.Message)
	}
}
