package controllers

import (
	"fmt"
	"jxb-eprocurement/handlers"
	"jxb-eprocurement/handlers/dtos"
	"jxb-eprocurement/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RoleController struct {
	DB *gorm.DB
}

func (ctrl *RoleController) CreateRole(c *gin.Context) {
	var role models.Role
	if err := c.ShouldBindJSON(&role); err != nil {
		handlers.ResponseFormatter(c, http.StatusBadRequest, nil, "Invalid Input")
		return
	}

	if err := handlers.ValidateStruct(role); err != nil {
		handlers.ValidationErrorHandler(c, err)
		return
	}

	var existingRole models.Role
	if err := ctrl.DB.Where("name = ?", role.Name).First(&existingRole).Error; err == nil {
		handlers.ResponseFormatter(c, http.StatusConflict, nil, fmt.Sprintf("Role with name %s is already exist", role.Name))
		return
	}

	if err := ctrl.DB.Create(&role).Error; err != nil {
		handlers.ResponseFormatter(c, http.StatusInternalServerError, nil, "Internal Server Error")
		return
	}

	handlers.ResponseFormatter(c, http.StatusCreated, nil, "Role created successfully")
}

func (ctrl *RoleController) GetAllRoles(c *gin.Context) {
	var roles []models.Role
	if err := ctrl.DB.Find(&roles).Error; err != nil {
		handlers.ResponseFormatter(c, http.StatusInternalServerError, nil, "Internal Server Error")
		return
	}

	var roleDTOs []dtos.RoleDTO
	for _, role := range roles {
		roleDTOs = append(roleDTOs, dtos.ToRoleDTO(role))
	}

	handlers.ResponseFormatter(c, http.StatusOK, roleDTOs, "Roles retrieved successfully")
}

func (ctrl *RoleController) GetRole(c *gin.Context) {
	var role models.Role
	id := c.Param("id")

	if err := ctrl.DB.First(&role, id).Error; err != nil {
		handlers.ResponseFormatter(c, http.StatusNotFound, nil, "Role not found")
		return
	}

	handlers.ResponseFormatter(c, http.StatusOK, dtos.ToRoleDTO(role), "Role retrieved successfully")
}

func (ctrl *RoleController) UpdateRole(c *gin.Context) {
	var role models.Role
	id := c.Param("id")

	if err := ctrl.DB.First(&role, id).Error; err != nil {
		handlers.ResponseFormatter(c, http.StatusNotFound, nil, "Role not found")
		return
	}

	if err := c.ShouldBindJSON(&role); err != nil {
		handlers.ResponseFormatter(c, http.StatusBadRequest, nil, "Invalid Input")
		return
	}

	var existingRole models.Role
	if err := ctrl.DB.Where("name = ? AND id != ?", role.Name, role.ID).First(&existingRole).Error; err == nil {
		handlers.ResponseFormatter(c, http.StatusConflict, nil, fmt.Sprintf("Role with name %s is already exist", role.Name))
		return
	}

	if err := handlers.ValidateStruct(role); err != nil {
		handlers.ValidationErrorHandler(c, err)
		return
	}

	if err := ctrl.DB.Save(&role).Error; err != nil {
		handlers.ResponseFormatter(c, http.StatusInternalServerError, nil, "Internal Server Error")
		return
	}

	handlers.ResponseFormatter(c, http.StatusOK, nil, "Role updated successfully")
}

func (ctrl *RoleController) DeleteRole(c *gin.Context) {
	var role models.Role
	id := c.Param("id")

	if err := ctrl.DB.First(&role, id).Error; err != nil {
		handlers.ResponseFormatter(c, http.StatusNotFound, nil, "Role not found")
		return
	}

	var userCount int64
	if err := ctrl.DB.Model(&models.User{}).Where("role_id = ?", role.ID).Count(&userCount).Error; err != nil {
		handlers.ResponseFormatter(c, http.StatusInternalServerError, nil, "Internal Server Error")
		return
	}

	if userCount > 0 {
		handlers.ResponseFormatter(c, http.StatusForbidden, nil, "Role cannot be deleted, because there are users with that role")
		return
	}

	if err := ctrl.DB.Delete(&role).Error; err != nil {
		handlers.ResponseFormatter(c, http.StatusInternalServerError, nil, "Internal Server Error")
		return
	}

	handlers.ResponseFormatter(c, http.StatusOK, nil, "Role deleted successfully")
}
