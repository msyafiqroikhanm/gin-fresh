package controllers

import (
	"jxb-eprocurement/handlers"
	"jxb-eprocurement/handlers/dtos"
	"jxb-eprocurement/service"
	"net/http"
	"strconv"

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
	handlers.ResponseFormatter(c, response.Status, response.Data, response.Message)
}

// GetModule handles the request to get a module by ID.
func (mc *ModuleControllerImpl) GetModule(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		handlers.ResponseFormatter(c, http.StatusBadRequest, nil, "Invalid ID")
		return
	}
	response := mc.service.GetByID(c, uint(id))
	handlers.ResponseFormatter(c, response.Status, response.Data, response.Message)
}

// CreateModule handles the request to add a new module.
func (mc *ModuleControllerImpl) CreateModule(c *gin.Context) {
	var moduleDTO dtos.USRModuleMinimalDTO
	if err := c.ShouldBind(&moduleDTO); err != nil {
		handlers.ResponseFormatter(c, http.StatusBadRequest, nil, "Invalid input")
		return
	}

	response := mc.service.AddData(c, moduleDTO)
	if response.Err != nil {
		handlers.ResponseFormatter(c, response.Status, response.Err, response.Message)
	} else {
		handlers.ResponseFormatter(c, response.Status, response.Data, response.Message)
	}
}

// UpdateModule handles the request to update a module.
func (mc *ModuleControllerImpl) UpdateModule(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		handlers.ResponseFormatter(c, http.StatusBadRequest, nil, "Invalid ID")
		return
	}
	var moduleDTO dtos.USRModuleMinimalDTO
	if err := c.ShouldBindJSON(&moduleDTO); err != nil {
		handlers.ResponseFormatter(c, http.StatusBadRequest, nil, "Invalid input")
		return
	}
	response := mc.service.UpdateData(c, uint(id), moduleDTO)
	handlers.ResponseFormatter(c, response.Status, response.Data, response.Message)
}

// DeleteModule handles the request to delete a module.
func (mc *ModuleControllerImpl) DeleteModule(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		handlers.ResponseFormatter(c, http.StatusBadRequest, nil, "Invalid ID")
		return
	}
	response := mc.service.DeleteData(c, uint(id))
	handlers.ResponseFormatter(c, response.Status, response.Data, response.Message)
}
