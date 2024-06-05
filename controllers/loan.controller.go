package controllers

import (
	"jxb-eprocurement/handlers"
	"jxb-eprocurement/handlers/dtos"
	"jxb-eprocurement/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type LoanController struct {
	DB *gorm.DB
}

func (ctrl *LoanController) CreateLoan(c *gin.Context) {
	var loanValidation models.LoanValidation

	// Bind JSON input to loanValidation struct
	if err := c.ShouldBindJSON(&loanValidation); err != nil {
		handlers.ResponseFormatter(c, http.StatusBadRequest, err.Error(), "Invalid Input")
		return
	}

	// Validate the struct
	if err := handlers.ValidateStruct(loanValidation); err != nil {
		handlers.ValidationErrorHandler(c, err)
		return
	}

	// If valid, create the actual Loan model
	loan := models.Loan{
		UserID:    loanValidation.UserID,
		VehicleID: loanValidation.VehicleID,
		Purpose:   loanValidation.Purpose,
		// ReturnTime: nil,
		// Set other fields if necessary
	}

	var vehicle models.Vehicle
	if err := ctrl.DB.Where("is_available=?", true).First(&vehicle, loan.VehicleID).Error; err != nil {
		handlers.ResponseFormatter(c, http.StatusNotFound, nil, "Vehicle Not Available")
		return
	}

	vehicle.IsAvailable = false
	if err := ctrl.DB.Save(&vehicle).Error; err != nil {
		handlers.ResponseFormatter(c, http.StatusInternalServerError, nil, "Internal Server Error")
		return
	}

	if err := ctrl.DB.Create(&loan).Error; err != nil {
		handlers.ResponseFormatter(c, http.StatusInternalServerError, nil, err.Error())
		return
	}

	handlers.ResponseFormatter(c, http.StatusCreated, nil, "Loan created successfully")
}

func (ctrl *LoanController) GetAllLoans(c *gin.Context) {
	var loans []models.Loan
	if err := ctrl.DB.Preload("Vehicle").Preload("User").Find(&loans).Error; err != nil {
		handlers.ResponseFormatter(c, http.StatusInternalServerError, nil, "Internal Server Error")
		return
	}

	loanDTOs := dtos.ToLoanDTOs(loans)
	handlers.ResponseFormatter(c, http.StatusOK, loanDTOs, "Loans retrieved successfully")
}

func (ctrl *LoanController) GetLoan(c *gin.Context) {
	var loan models.Loan
	id := c.Param("id")

	if err := ctrl.DB.Preload("Vehicle").Preload("User").First(&loan, id).Error; err != nil {
		handlers.ResponseFormatter(c, http.StatusNotFound, nil, "Loan not found")
		return
	}

	handlers.ResponseFormatter(c, http.StatusOK, dtos.ToLoanDTO(loan), "Loan retrieved successfully")
}

func (ctrl *LoanController) UpdateLoan(c *gin.Context) {
	var loan models.Loan
	id := c.Param("id")

	if err := ctrl.DB.First(&loan, id).Error; err != nil {
		handlers.ResponseFormatter(c, http.StatusNotFound, nil, "Loan not found")
		return
	}

	var loanValidation models.LoanValidation

	// Bind JSON input to loanValidation struct
	if err := c.ShouldBindJSON(&loanValidation); err != nil {
		handlers.ResponseFormatter(c, http.StatusBadRequest, nil, "Invalid Input")
		return
	}

	// Validate the struct
	if err := handlers.ValidateStruct(loanValidation); err != nil {
		handlers.ValidationErrorHandler(c, err)
		return
	}

	// If valid, create the actual Loan model
	input := models.Loan{
		UserID:    loanValidation.UserID,
		VehicleID: loanValidation.VehicleID,
		Purpose:   loanValidation.Purpose,
		// Set other fields if necessary
	}

	var vehicle models.Vehicle
	// if err := ctrl.DB.Where("is_available = ? OR id != ?", true, id).First(&vehicle).Error; err == nil {
	// 	handlers.ResponseFormatter(c, http.StatusNotFound, nil, "Vehicle Not Available")
	// 	return
	// }

	// if err := ctrl.DB.Preload("Type").First(&vehicle, id).Error; err != nil {
	// 	handlers.ResponseFormatter(c, http.StatusNotFound, nil, "Vehicle not found")
	// 	return
	// }
	if input.VehicleID != loan.VehicleID {
		if err := ctrl.DB.Where("is_available = ?", true).First(&vehicle, input.VehicleID).Error; err != nil {
			handlers.ResponseFormatter(c, http.StatusNotFound, nil, "Vehicle Not Available")
			return
		}

		var vehicleOld models.Vehicle
		if err := ctrl.DB.First(&vehicleOld, loan.VehicleID).Error; err != nil {
			handlers.ResponseFormatter(c, http.StatusNotFound, nil, "Vehicle Not Available")
			return
		}

		vehicleOld.IsAvailable = true
		if err := ctrl.DB.Save(&vehicleOld).Error; err != nil {
			handlers.ResponseFormatter(c, http.StatusInternalServerError, nil, "Internal Server Error")
			return
		}

		var vehicleNew models.Vehicle
		if err := ctrl.DB.First(&vehicleNew, input.VehicleID).Error; err != nil {
			handlers.ResponseFormatter(c, http.StatusNotFound, nil, "Vehicle Not Available")
			return
		}

		vehicleNew.IsAvailable = false
		if err := ctrl.DB.Save(&vehicleNew).Error; err != nil {
			handlers.ResponseFormatter(c, http.StatusInternalServerError, nil, "Internal Server Error")
			return
		}

		loan.VehicleID = input.VehicleID
	}

	// handlers.ResponseFormatter(c, http.StatusOK, vehicle, "Loan updated successfully")
	// return

	loan.UserID = input.UserID
	loan.Purpose = input.Purpose

	if err := ctrl.DB.Save(&loan).Error; err != nil {
		handlers.ResponseFormatter(c, http.StatusInternalServerError, nil, "Internal Server Error")
		return
	}

	handlers.ResponseFormatter(c, http.StatusOK, nil, "Loan updated successfully")
}

func (ctrl *LoanController) DeleteLoan(c *gin.Context) {
	var loan models.Loan
	id := c.Param("id")

	if err := ctrl.DB.First(&loan, id).Error; err != nil {
		handlers.ResponseFormatter(c, http.StatusNotFound, nil, "Loan not found")
		return
	}

	if err := ctrl.DB.Delete(&loan).Error; err != nil {
		handlers.ResponseFormatter(c, http.StatusInternalServerError, nil, "Internal Server Error")
		return
	}

	var vehicle models.Vehicle
	if err := ctrl.DB.First(&vehicle, loan.VehicleID).Error; err != nil {
		handlers.ResponseFormatter(c, http.StatusNotFound, nil, "Vehicle Not Available")
		return
	}

	vehicle.IsAvailable = true
	if err := ctrl.DB.Save(&vehicle).Error; err != nil {
		handlers.ResponseFormatter(c, http.StatusInternalServerError, nil, "Internal Server Error")
		return
	}

	handlers.ResponseFormatter(c, http.StatusOK, nil, "Loan deleted successfully")
}

func (ctrl *LoanController) ReturnLoan(c *gin.Context) {
	var loan models.Loan
	id := c.Param("id")

	if err := ctrl.DB.First(&loan, id).Error; err != nil {
		handlers.ResponseFormatter(c, http.StatusNotFound, nil, "Loan not found")
		return
	}

	loan.ReturnTime = time.Now()

	if err := ctrl.DB.Save(&loan).Error; err != nil {
		handlers.ResponseFormatter(c, http.StatusInternalServerError, nil, "Internal Server Error")
		return
	}

	var vehicle models.Vehicle
	if err := ctrl.DB.First(&vehicle, loan.VehicleID).Error; err != nil {
		handlers.ResponseFormatter(c, http.StatusNotFound, nil, "Vehicle Not Available")
		return
	}

	vehicle.IsAvailable = true
	if err := ctrl.DB.Save(&vehicle).Error; err != nil {
		handlers.ResponseFormatter(c, http.StatusInternalServerError, nil, "Internal Server Error")
		return
	}

	handlers.ResponseFormatter(c, http.StatusOK, nil, "Loan returned successfully")
}
