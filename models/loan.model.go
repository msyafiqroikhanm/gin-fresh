package models

import (
	"time"

	"gorm.io/gorm"
)

type Loan struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `json:"user_id" `
	VehicleID  uint      `json:"vehicle_id" validate:"required"`
	Purpose    string    `json:"purpose" validate:"required"`
	ReturnTime time.Time `json:"return_time"`
	User       User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Vehicle    Vehicle   `gorm:"foreignKey:VehicleID" json:"vehicle,omitempty"`
	gorm.Model
}

func (Loan) TableName() string {
	return "loans"
}

type LoanValidation struct {
	UserID    uint   `json:"user_id" validate:"required"`
	VehicleID uint   `json:"vehicle_id" validate:"required"`
	Purpose   string `json:"purpose" validate:"required"`
}
