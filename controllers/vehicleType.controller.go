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

type VehicleTypeController struct {
	DB *gorm.DB
}

func (ctrl *VehicleTypeController) CreateVehicleType(c *gin.Context) {
	var vehicleType models.VehicleType
	if err := c.ShouldBindJSON(&vehicleType); err != nil {
		handlers.ResponseFormatter(c, http.StatusBadRequest, nil, "Invalid Input")
		return
	}

	if err := handlers.ValidateStruct(vehicleType); err != nil {
		handlers.ValidationErrorHandler(c, err)
		return
	}

	var existingVehicleType models.VehicleType
	if err := ctrl.DB.Where("name = ?", vehicleType.Name).First(&existingVehicleType).Error; err == nil {
		handlers.ResponseFormatter(c, http.StatusConflict, nil, fmt.Sprintf("Vehicle Type with name %s is already exist", vehicleType.Name))
		return
	}

	if err := ctrl.DB.Create(&vehicleType).Error; err != nil {
		handlers.ResponseFormatter(c, http.StatusInternalServerError, nil, "Internal Server Error")
		return
	}

	handlers.ResponseFormatter(c, http.StatusCreated, nil, "Vehicle Type created successfully")
}

func (ctrl *VehicleTypeController) GetAllVehicleTypes(c *gin.Context) {
	var vehicleTypes []models.VehicleType
	if err := ctrl.DB.Find(&vehicleTypes).Error; err != nil {
		handlers.ResponseFormatter(c, http.StatusInternalServerError, nil, "Internal Server Error")
		return
	}

	var vehicleTypeDTOs []dtos.VehicleTypeDTO
	for _, vehicleType := range vehicleTypes {
		vehicleTypeDTOs = append(vehicleTypeDTOs, dtos.ToVehicleTypeDTO(vehicleType))
	}

	handlers.ResponseFormatter(c, http.StatusOK, vehicleTypeDTOs, "Vehicle Types retrieved successfully")
}

func (ctrl *VehicleTypeController) GetVehicleType(c *gin.Context) {
	var vehicleType models.VehicleType
	id := c.Param("id")

	if err := ctrl.DB.First(&vehicleType, id).Error; err != nil {
		handlers.ResponseFormatter(c, http.StatusNotFound, nil, "Vehicle Type not found")
		return
	}

	handlers.ResponseFormatter(c, http.StatusOK, dtos.ToVehicleTypeDTO(vehicleType), "Vehicle Type retrieved successfully")
}

func (ctrl *VehicleTypeController) UpdateVehicleType(c *gin.Context) {
	var vehicleType models.VehicleType
	id := c.Param("id")

	if err := ctrl.DB.First(&vehicleType, id).Error; err != nil {
		handlers.ResponseFormatter(c, http.StatusNotFound, nil, "Vehicle Type not found")
		return
	}

	if err := c.ShouldBindJSON(&vehicleType); err != nil {
		handlers.ResponseFormatter(c, http.StatusBadRequest, nil, "Invalid Input")
		return
	}

	var existingVehicleType models.VehicleType
	if err := ctrl.DB.Where("name = ? AND id != ?", vehicleType.Name, vehicleType.ID).First(&existingVehicleType).Error; err == nil {
		handlers.ResponseFormatter(c, http.StatusConflict, nil, fmt.Sprintf("Vehicle Type with name %s is already exist", vehicleType.Name))
		return
	}

	if err := handlers.ValidateStruct(vehicleType); err != nil {
		handlers.ValidationErrorHandler(c, err)
		return
	}

	if err := ctrl.DB.Save(&vehicleType).Error; err != nil {
		handlers.ResponseFormatter(c, http.StatusInternalServerError, nil, "Internal Server Error")
		return
	}

	handlers.ResponseFormatter(c, http.StatusOK, nil, "Vehicle Type updated successfully")
}

func (ctrl *VehicleTypeController) DeleteVehicleType(c *gin.Context) {
	var vehicleType models.VehicleType
	id := c.Param("id")

	if err := ctrl.DB.First(&vehicleType, id).Error; err != nil {
		handlers.ResponseFormatter(c, http.StatusNotFound, nil, "Vehicle Type not found")
		return
	}

	var userCount int64
	if err := ctrl.DB.Model(&models.User{}).Where("role_id = ?", vehicleType.ID).Count(&userCount).Error; err != nil {
		handlers.ResponseFormatter(c, http.StatusInternalServerError, nil, "Internal Server Error")
		return
	}

	if userCount > 0 {
		handlers.ResponseFormatter(c, http.StatusForbidden, nil, "Vehicle Type cannot be deleted, because there are users with that vehicleType")
		return
	}

	if err := ctrl.DB.Delete(&vehicleType).Error; err != nil {
		handlers.ResponseFormatter(c, http.StatusInternalServerError, nil, "Internal Server Error")
		return
	}

	handlers.ResponseFormatter(c, http.StatusOK, nil, "Vehicle Type deleted successfully")
}
