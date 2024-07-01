package models

import (
	"gorm.io/gorm"
)

type USR_Feature struct {
	ID       uint        `json:"id" gorm:"primaryKey" `
	ModuleID uint        `json:"module_id" validate:"required"`
	Name     string      `json:"name" validate:"required"`
	Module   USR_Module  `json:"module" gorm:"foreignKey:ModuleID"`
	Roles    []*USR_Role `gorm:"many2many:usr_rolefeatures;" json:"roles"`
	gorm.Model
}

func (USR_Feature) TableName() string {
	return "usr_features"
}
