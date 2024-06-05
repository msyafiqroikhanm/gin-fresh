package models

import "gorm.io/gorm"

type Role struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `gorm:"not null" json:"name" validate:"required"`
	gorm.Model
}

func (Role) TableName() string {
	return "roles"
}
