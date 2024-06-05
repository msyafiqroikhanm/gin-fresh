package controllers

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

type VehicleController struct {
	DB *gorm.DB
}

func (ctrl *VehicleController) CreateVehicle(c *gin.Context) {
	var vehicle models.Vehicle
	if err := c.ShouldBindJSON(&vehicle); err != nil {
		handlers.ResponseFormatter(c, http.StatusBadRequest, nil, "Invalid Input")
		return
	}

	if err := handlers.ValidateStruct(vehicle); err != nil {
		handlers.ValidationErrorHandler(c, err)
		return
	}

	var checkDuplicatePoliceNumber models.Vehicle
	if err := ctrl.DB.Where("police_number = ?", vehicle.PoliceNumber).First(&checkDuplicatePoliceNumber).Error; err == nil {
		handlers.ResponseFormatter(c, http.StatusConflict, nil, fmt.Sprintf("Vehicle with police number '%s' is already exist", vehicle.PoliceNumber))
		return
	}

	var vehicleType models.VehicleType
	if err := ctrl.DB.First(&vehicleType, vehicle.TypeID).Error; err != nil {
		handlers.ResponseFormatter(c, http.StatusBadRequest, nil, "Invalid type ID")
		return
	}

	if err := ctrl.DB.Create(&vehicle).Error; err != nil {
		handlers.ResponseFormatter(c, http.StatusInternalServerError, nil, "Internal Server Error")
		return
	}

	handlers.ResponseFormatter(c, http.StatusCreated, nil, "Vehicle created successfully")
}

func (ctrl *VehicleController) GetAllVehicles(c *gin.Context) {
	var vehicles []models.Vehicle
	query := ctrl.DB.Preload("Type")

	// Check if type_id query parameter is provided
	typeIDStr := c.Query("type_id")
	if typeIDStr != "" {
		typeID, err := strconv.Atoi(typeIDStr)
		if err != nil {
			handlers.ResponseFormatter(c, http.StatusBadRequest, nil, "Invalid type_id")
			return
		}
		query = query.Where("type_id = ?", typeID)
	}

	// Check if is_available query parameter is provided
	isAvailableStr := c.Query("is_available")
	if isAvailableStr != "" {
		isAvailable, err := strconv.ParseBool(isAvailableStr)
		if err != nil {
			handlers.ResponseFormatter(c, http.StatusBadRequest, nil, "Invalid is_available")
			return
		}
		query = query.Where("is_available = ?", isAvailable)
	}
	if err := query.Find(&vehicles).Error; err != nil {
		handlers.ResponseFormatter(c, http.StatusInternalServerError, nil, "Internal Server Error")
		return
	}

	var vehicleDTOs []dtos.VehicleDTO
	for _, vehicle := range vehicles {
		vehicleDTOs = append(vehicleDTOs, dtos.ToVehicleDTO(vehicle))
	}

	handlers.ResponseFormatter(c, http.StatusOK, vehicleDTOs, "Vehicles retrieved successfully")
}

func (ctrl *VehicleController) GetVehicle(c *gin.Context) {
	var vehicle models.Vehicle
	id := c.Param("id")

	if err := ctrl.DB.Preload("Type").First(&vehicle, id).Error; err != nil {
		handlers.ResponseFormatter(c, http.StatusNotFound, nil, "Vehicle not found")
		return
	}

	handlers.ResponseFormatter(c, http.StatusOK, dtos.ToVehicleDTO(vehicle), "Vehicle retrieved successfully")
}

func (ctrl *VehicleController) UpdateVehicle(c *gin.Context) {
	var vehicle models.Vehicle
	id := c.Param("id")

	if err := ctrl.DB.First(&vehicle, id).Error; err != nil {
		handlers.ResponseFormatter(c, http.StatusNotFound, nil, "Vehicle not found")
		return
	}

	var input models.Vehicle
	if err := c.ShouldBindJSON(&input); err != nil {
		handlers.ResponseFormatter(c, http.StatusBadRequest, nil, "Invalid Inputs")
		return
	}

	// Check if police_number is already used by another vehicle
	var checkDuplicatePoliceNumber models.Vehicle
	if err := ctrl.DB.Where("police_number = ? AND id != ?", input.PoliceNumber, id).First(&checkDuplicatePoliceNumber).Error; err == nil {
		handlers.ResponseFormatter(c, http.StatusConflict, nil, fmt.Sprintf("Vehicle with police number '%s' is already exist", input.PoliceNumber))
		return
	}

	if err := handlers.ValidateStruct(input); err != nil {
		handlers.ValidationErrorHandler(c, err)
		return
	}

	// Check if the role ID is valid
	var vehicleType models.VehicleType
	if err := ctrl.DB.First(&vehicleType, input.TypeID).Error; err != nil {
		handlers.ResponseFormatter(c, http.StatusBadRequest, nil, "Invalid type ID")
		return
	}

	vehicle.Name = input.Name
	vehicle.PoliceNumber = input.PoliceNumber
	vehicle.TypeID = input.TypeID

	if err := ctrl.DB.Save(&vehicle).Error; err != nil {
		handlers.ResponseFormatter(c, http.StatusInternalServerError, nil, "Internal Server Error")
		return
	}

	handlers.ResponseFormatter(c, http.StatusOK, nil, "Vehicle updated successfully")
}

func (ctrl *VehicleController) DeleteVehicle(c *gin.Context) {
	var vehicle models.Vehicle
	id := c.Param("id")

	if err := ctrl.DB.First(&vehicle, id).Error; err != nil {
		handlers.ResponseFormatter(c, http.StatusNotFound, nil, "Vehicle not found")
		return
	}

	if err := ctrl.DB.Delete(&vehicle).Error; err != nil {
		handlers.ResponseFormatter(c, http.StatusInternalServerError, nil, "Internal Server Error")
		return
	}

	handlers.ResponseFormatter(c, http.StatusOK, nil, "Vehicle deleted successfully")
}
