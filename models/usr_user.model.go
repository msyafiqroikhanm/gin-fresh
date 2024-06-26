package models

import "gorm.io/gorm"

type USR_User struct {
	ID       uint     `gorm:"primaryKey" json:"id"`
	RoleID   uint     `json:"role_id"`
	Name     string   `json:"name" validate:"required,min=3,max=100"`
	Email    string   `json:"email" validate:"required,email"`
	Password string   `json:"password" validate:"required"`
	Role     USR_Role `json:"role" gorm:"foreignKey:RoleID"`
	gorm.Model
}

func (USR_User) TableName() string {
	return "usr_users"
}
