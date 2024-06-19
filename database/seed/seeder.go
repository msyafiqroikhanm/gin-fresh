package seed

import (
	"gorm.io/gorm"
)

func Seed(db *gorm.DB) {
	// // Seed Roles
	// var roleCount int64
	// db.Model(&models.Role{}).Count(&roleCount)
	// if roleCount == 0 {
	// 	roles := []models.Role{
	// 		{Name: "admin"},
	// 		{Name: "user"},
	// 	}
	// 	// db.Create(&roles)
	// 	if err := db.Create(&roles).Error; err != nil {
	// 		log.Fatalf("Error seeding User: %v", err)
	// 	} else {
	// 		log.Println("User seeded successfully.")
	// 	}
	// }

	// // Seed Users
	// var userCount int64
	// db.Model(&models.User{}).Count(&userCount)
	// if userCount == 0 {
	// 	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	// 	if err != nil {
	// 		log.Fatalf("Failed to hash password: %v", err)
	// 	}
	// 	users := []models.User{
	// 		{
	// 			Name:     "Admin User",
	// 			Email:    "admin@example.com",
	// 			Password: string(hashedPassword),
	// 			RoleID:   1, // Assuming 1 is the ID for admin role
	// 		},
	// 		{
	// 			Name:     "Regular User",
	// 			Email:    "user@example.com",
	// 			Password: string(hashedPassword),
	// 			RoleID:   2, // Assuming 2 is the ID for regular user role
	// 		},
	// 	}
	// 	// db.Create(&users)
	// 	if err := db.Create(&users).Error; err != nil {
	// 		log.Fatalf("Error seeding Role: %v", err)
	// 	} else {
	// 		log.Println("Role seeded successfully.")
	// 	}
	// }

	// // Seed Vehicle Types
	// var vehicleTypeCount int64
	// db.Model(&models.VehicleType{}).Count(&vehicleTypeCount)
	// if vehicleTypeCount == 0 {
	// 	vehicleTypes := []models.VehicleType{
	// 		{Name: "Motorcycle"},
	// 		{Name: "Car"},
	// 		{Name: "Minibus"},
	// 	}
	// 	if err := db.Create(&vehicleTypes).Error; err != nil {
	// 		log.Fatalf("Error seeding vehicle types: %v", err)
	// 	} else {
	// 		log.Println("Vehicle types seeded successfully.")
	// 	}
	// }

	// // Get vehicle types
	// var vehicleTypes []models.VehicleType
	// if err := db.Find(&vehicleTypes).Error; err != nil {
	// 	log.Fatalf("Error retrieving vehicle types: %v", err)
	// }

	// // Seed Vehicles
	// for _, vehicleType := range vehicleTypes {
	// 	var vehicleCount int64
	// 	db.Model(&models.Vehicle{}).Where("type_id = ?", vehicleType.ID).Count(&vehicleCount)
	// 	if vehicleCount == 0 {
	// 		var vehicles []models.Vehicle
	// 		switch vehicleType.Name {
	// 		case "Motorcycle":
	// 			vehicles = []models.Vehicle{
	// 				{Name: "Yamaha NMAX", TypeID: vehicleType.ID, PoliceNumber: "B1234XYZ", IsAvailable: true},
	// 				{Name: "Honda Beat", TypeID: vehicleType.ID, PoliceNumber: "B5678ABC", IsAvailable: true},
	// 			}
	// 		case "Car":
	// 			vehicles = []models.Vehicle{
	// 				{Name: "Toyota Avanza", TypeID: vehicleType.ID, PoliceNumber: "B9101DEF", IsAvailable: true},
	// 				{Name: "Honda Civic", TypeID: vehicleType.ID, PoliceNumber: "B1121GHI", IsAvailable: true},
	// 			}
	// 		case "Minibus":
	// 			vehicles = []models.Vehicle{
	// 				{Name: "Isuzu Elf", TypeID: vehicleType.ID, PoliceNumber: "B3141JKL", IsAvailable: true},
	// 				{Name: "Toyota Hiace", TypeID: vehicleType.ID, PoliceNumber: "B5161MNO", IsAvailable: true},
	// 			}
	// 		}

	// 		if err := db.Create(&vehicles).Error; err != nil {
	// 			log.Fatalf("Error seeding vehicles for type %s: %v", vehicleType.Name, err)
	// 		} else {
	// 			log.Printf("Vehicles for type %s seeded successfully.", vehicleType.Name)
	// 		}
	// 	}
	// }
}
