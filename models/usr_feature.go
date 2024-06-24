package models

import (
	"gorm.io/gorm"
)

type USR_Feature struct {
	ID       uint       `json:"id" gorm:"primaryKey" `
	ModuleID uint       `json:"module_id" validate:"required"`
	Name     string     `json:"name" validate:"required"`
	Module   USR_Module `json:"module" gorm:"foreignKey:ModuleID"`
	gorm.Model
}

func (USR_Feature) TableName() string {
	return "usr_features"
}
