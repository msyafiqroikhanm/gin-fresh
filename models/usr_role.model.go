package models

import "gorm.io/gorm"

type USR_Role struct {
	ID               uint   `gorm:"primaryKey" json:"id"`
	Name             string `json:"name" validate:"required"`
	IsAdministrative bool   `json:"is_administrative"`
	gorm.Model
}

func (USR_Role) TableName() string {
	return "usr_roles"
}
