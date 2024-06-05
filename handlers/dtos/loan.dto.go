package dtos

import (
	"jxb-eprocurement/models"
	"time"
)

type LoanDTO struct {
	ID         uint       `gorm:"primaryKey" json:"id"`
	Purpose    string     `json:"purpose" validate:"required"`
	ReturnTime time.Time  `json:"return_time"`
	User       UserDTO    `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Vehicle    VehicleDTO `gorm:"foreignKey:VehicleID" json:"vehicle,omitempty"`
}

func ToLoanDTO(loan models.Loan) LoanDTO {
	return LoanDTO{
		ID:         loan.ID,
		Purpose:    loan.Purpose,
		ReturnTime: loan.ReturnTime,
		User:       ToUserDTO(loan.User),
		Vehicle:    ToVehicleDTO(loan.Vehicle),
	}
}
