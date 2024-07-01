package models

import "gorm.io/gorm"

type USR_Module struct {
	ID       uint          `gorm:"primaryKey" json:"id"`
	Name     string        `json:"name" validate:"required"`
	ParentID *uint         `json:"parent_id"`
	Child    []USR_Module  `gorm:"foreignkey:ParentID"`
	Features []USR_Feature `gorm:"foreignKey:ModuleID"`
	gorm.Model
}

func (USR_Module) TableName() string {
	return "usr_modules"
}
