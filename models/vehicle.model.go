package models

import "gorm.io/gorm"

type Vehicle struct {
	ID           uint        `gorm:"primaryKey" json:"id"`
	TypeID       uint        `json:"type_id" validate:"required"`
	Name         string      `json:"name" validate:"required,min=3,max=100"`
	PoliceNumber string      `json:"police_number" validate:"required"`
	IsAvailable  bool        `json:"is_available"`
	Type         VehicleType `json:"type" gorm:"foreignKey:TypeID" validate:"-"`
	gorm.Model
}

func (Vehicle) TableName() string {
	return "vehicles"
}
