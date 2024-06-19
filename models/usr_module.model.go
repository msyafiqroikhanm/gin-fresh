package models

import "gorm.io/gorm"

type USR_Module struct {
	ID       uint         `gorm:"primaryKey" json:"id"`
	Name     string       `json:"name"`
	ParentID *uint        `json:"parent_id"`
	Child    []USR_Module `gorm:"foreignkey:ParentID"`
	gorm.Model
}

func (USR_Module) TableName() string {
	return "usr_modules"
}
