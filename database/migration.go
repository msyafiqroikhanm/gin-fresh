package database

import (
	"jxb-eprocurement/database/seed"
	"jxb-eprocurement/models"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&models.USR_Module{})
	db.AutoMigrate(&models.USR_Feature{})
	db.AutoMigrate(&models.USR_Role{})
	db.AutoMigrate(&models.USR_User{})

	// Seed initial data
	seed.Seed(db)
}
