package models

import "gorm.io/gorm"

type VehicleType struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `json:"name"`
	gorm.Model
}
