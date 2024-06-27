package models

import "gorm.io/gorm"

// var DB *gorm.DB // Variabel DB untuk menyimpan koneksi database
type User struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	RoleID   uint   `json:"role_id"`
	Name     string `json:"name" validate:"required,min=3,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	Role     Role   `json:"role" gorm:"foreignKey:RoleID" validate:"-"`
	gorm.Model
}

func (User) TableName() string {
	return "users"
}
