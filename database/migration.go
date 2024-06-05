package database

import (
	"jxb-eprocurement/models"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&models.Role{})
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Vehicle{})
	db.AutoMigrate(&models.VehicleType{})
	db.AutoMigrate(&models.Loan{})

	// Seed initial data
	Seed(db)
}
