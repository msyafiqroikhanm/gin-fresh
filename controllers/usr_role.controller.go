package controllers

import (
	"jxb-eprocurement/handlers"
	"jxb-eprocurement/service"

	"github.com/gin-gonic/gin"
)

type USRRoleController interface {
	GetAllRoles(c *gin.Context)
	GetRole(c *gin.Context)
	CreateRole(c *gin.Context)
	UpdateRole(c *gin.Context)
	DeleteRole(c *gin.Context)
}

// RoleControllerImpl is the implementation of the RoleController interface.
type RoleControllerImpl struct {
	service service.RoleService
}

// NewRoleController creates a new instance of RoleControllerImpl.
func RoleControllerConstructor(service service.RoleService) USRRoleController {
	return &RoleControllerImpl{service: service}
}

// GetAllRoles handles the request to get all roles.
func (mc *RoleControllerImpl) GetAllRoles(c *gin.Context) {
	response := mc.service.GetAll(c)
	if response.Err != nil {
		handlers.ResponseFormatter(c, response.Status, response.Err, response.Message)
	} else {
		handlers.ResponseFormatter(c, response.Status, response.Data, response.Message)
	}
}

// GetRole handles the request to get a role by ID.
func (mc *RoleControllerImpl) GetRole(c *gin.Context) {
	response := mc.service.GetByID(c)
	if response.Err != nil {
		handlers.ResponseFormatter(c, response.Status, response.Err, response.Message)
	} else {
		handlers.ResponseFormatter(c, response.Status, response.Data, response.Message)
	}
}

// CreateRole handles the request to add a new role.
func (mc *RoleControllerImpl) CreateRole(c *gin.Context) {
	response := mc.service.AddData(c)
	if response.Err != nil {
		handlers.ResponseFormatter(c, response.Status, response.Err, response.Message)
	} else {
		handlers.ResponseFormatter(c, response.Status, response.Data, response.Message)
	}
}

// UpdateRole handles the request to update a role.
func (mc *RoleControllerImpl) UpdateRole(c *gin.Context) {
	response := mc.service.UpdateData(c)
	if response.Err != nil {
		handlers.ResponseFormatter(c, response.Status, response.Err, response.Message)
	} else {
		handlers.ResponseFormatter(c, response.Status, response.Data, response.Message)
	}
}

// DeleteRole handles the request to delete a role.
func (mc *RoleControllerImpl) DeleteRole(c *gin.Context) {
	response := mc.service.DeleteData(c)
	if response.Err != nil {
		handlers.ResponseFormatter(c, response.Status, response.Err, response.Message)
	} else {
		handlers.ResponseFormatter(c, response.Status, response.Data, response.Message)
	}
}
