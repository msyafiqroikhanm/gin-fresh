package database

import (
	"jxb-eprocurement/database/seed"
	"jxb-eprocurement/models"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&models.Role{})
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Vehicle{})
	db.AutoMigrate(&models.VehicleType{})
	db.AutoMigrate(&models.Loan{})
	db.AutoMigrate(&models.USR_Module{})
	db.AutoMigrate(&models.USR_Feature{})

	// Seed initial data
	seed.Seed(db)
}
