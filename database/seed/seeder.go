package seed

import (
	"jxb-eprocurement/helpers"
	"jxb-eprocurement/models"
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Seed(db *gorm.DB) {
	// Seed models.USR_Module
	var modules []models.USR_Module
	{
		var moduleCount int64
		db.Model(&models.USR_Module{}).Count(&moduleCount)
		if moduleCount > 0 {
			log.Println("models.USR_Modules already seeded.")
			return
		}

		parentModules := []models.USR_Module{
			{Name: "Access Management"},
			{Name: "Vendor Management"},
			{Name: "Plan Management"},
			{Name: "Procurement Management"},
			{Name: "Contract Management"},
		}

		// Creating data using transaction
		err := db.Transaction(
			func(tx *gorm.DB) error {
				if err := tx.Create(&parentModules).Error; err != nil {
					return err
				}

				var subModules []models.USR_Module

				for _, parent := range parentModules {
					parentID := parent.ID // Ensure parent.ID is captured correctly in each iteration
					switch parent.Name {
					case "Access Management":
						subModules = append(subModules,
							models.USR_Module{Name: "Module", ParentID: &parentID},
							models.USR_Module{Name: "Feature", ParentID: &parentID},
							models.USR_Module{Name: "Role", ParentID: &parentID},
							models.USR_Module{Name: "User", ParentID: &parentID},
						)
					case "Vendor Management":
						subModules = append(subModules,
							models.USR_Module{Name: "Vendor", ParentID: &parentID},
							models.USR_Module{Name: "Vendor Profile", ParentID: &parentID},
						)
					case "Plan Management":
						subModules = append(subModules,
							models.USR_Module{Name: "RKA", ParentID: &parentID},
							models.USR_Module{Name: "DRUPP", ParentID: &parentID},
						)
					case "Procurement Management":
						subModules = append(subModules,
							models.USR_Module{Name: "Procurement", ParentID: &parentID},
							models.USR_Module{Name: "Interest", ParentID: &parentID},
							models.USR_Module{Name: "Offer", ParentID: &parentID},
							models.USR_Module{Name: "Invoice", ParentID: &parentID},
							models.USR_Module{Name: "Bill", ParentID: &parentID},
							models.USR_Module{Name: "Review", ParentID: &parentID},
						)
					case "Contract Management":
						subModules = append(subModules,
							models.USR_Module{Name: "Contract", ParentID: &parentID},
							models.USR_Module{Name: "Contract Monitoring", ParentID: &parentID},
							models.USR_Module{Name: "Addendum", ParentID: &parentID},
							models.USR_Module{Name: "SPMK", ParentID: &parentID},
						)
					}
				}

				if err := tx.Create(&subModules).Error; err != nil {
					return err
				}

				modules = append(modules, parentModules...)
				modules = append(modules, subModules...)

				return nil
			},
		)

		if err != nil {
			log.Fatalf("Error seeding models.USR_Modules: %v", err)
		} else {
			log.Println("models.USR_Modules seeded successfully.")
		}
	}

	// Seed Features
	var features []models.USR_Feature
	{
		var featureCount int64
		db.Model(&models.USR_Feature{}).Count(&featureCount)

		if featureCount == 0 {

			for _, module := range modules {
				switch module.Name {
				case "Module":
					features = append(
						features,
						models.USR_Feature{Name: "View Module", ModuleID: module.ID},
						models.USR_Feature{Name: "Create Module", ModuleID: module.ID},
						models.USR_Feature{Name: "Update Module", ModuleID: module.ID},
						models.USR_Feature{Name: "Delete Module", ModuleID: module.ID},
					)
				case "Feature":
					features = append(
						features,
						models.USR_Feature{Name: "View Feature", ModuleID: module.ID},
						models.USR_Feature{Name: "Create Feature", ModuleID: module.ID},
						models.USR_Feature{Name: "Update Feature", ModuleID: module.ID},
						models.USR_Feature{Name: "Delete Feature", ModuleID: module.ID},
					)
				case "Role":
					features = append(
						features,
						models.USR_Feature{Name: "View Role", ModuleID: module.ID},
						models.USR_Feature{Name: "Create Role", ModuleID: module.ID},
						models.USR_Feature{Name: "Update Role", ModuleID: module.ID},
						models.USR_Feature{Name: "Delete Role", ModuleID: module.ID},
					)
				case "User":
					features = append(
						features,
						models.USR_Feature{Name: "View User", ModuleID: module.ID},
						models.USR_Feature{Name: "Create User", ModuleID: module.ID},
						models.USR_Feature{Name: "Update User", ModuleID: module.ID},
						models.USR_Feature{Name: "Delete User", ModuleID: module.ID},
					)
				case "Vendor":
					features = append(
						features,
						models.USR_Feature{Name: "View Vendor", ModuleID: module.ID},
						models.USR_Feature{Name: "Update Vendor", ModuleID: module.ID},
						models.USR_Feature{Name: "Delete Vendor", ModuleID: module.ID},
						models.USR_Feature{Name: "Validate Vendor", ModuleID: module.ID},
						models.USR_Feature{Name: "Blacklist Vendor", ModuleID: module.ID},
						models.USR_Feature{Name: "View Blacklisted Vendor", ModuleID: module.ID},
					)
				case "Vendor Profile":
					features = append(
						features,
						models.USR_Feature{Name: "View Vendor Profile", ModuleID: module.ID},
						models.USR_Feature{Name: "Update Vendor Profile", ModuleID: module.ID},
					)
				case "RKA":
					features = append(
						features,
						models.USR_Feature{Name: "View RKA", ModuleID: module.ID},
						models.USR_Feature{Name: "Create RKA", ModuleID: module.ID},
						models.USR_Feature{Name: "Update RKA", ModuleID: module.ID},
						models.USR_Feature{Name: "Delete RKA", ModuleID: module.ID},
					)
				case "DRUPP":
					features = append(
						features,
						models.USR_Feature{Name: "View DRUPP", ModuleID: module.ID},
						models.USR_Feature{Name: "Create DRUPP", ModuleID: module.ID},
						models.USR_Feature{Name: "Update DRUPP", ModuleID: module.ID},
						models.USR_Feature{Name: "Delete DRUPP", ModuleID: module.ID},
					)
				case "Procurement":
					features = append(
						features,
						models.USR_Feature{Name: "View Procurement", ModuleID: module.ID},
						models.USR_Feature{Name: "Create Procurement", ModuleID: module.ID},
						models.USR_Feature{Name: "Update Procurement", ModuleID: module.ID},
						models.USR_Feature{Name: "Delete Procurement", ModuleID: module.ID},
						models.USR_Feature{Name: "Announce Procurement", ModuleID: module.ID},
						models.USR_Feature{Name: "Choose Winner Procurement", ModuleID: module.ID},
					)
				case "Interest":
					features = append(
						features,
						models.USR_Feature{Name: "View Interest", ModuleID: module.ID},
						models.USR_Feature{Name: "Create Interest", ModuleID: module.ID},
						models.USR_Feature{Name: "Update Interest", ModuleID: module.ID},
						models.USR_Feature{Name: "Delete Interest", ModuleID: module.ID},
					)
				case "Offer":
					features = append(
						features,
						models.USR_Feature{Name: "View Offer", ModuleID: module.ID},
						models.USR_Feature{Name: "Create Offer", ModuleID: module.ID},
						models.USR_Feature{Name: "Update Offer", ModuleID: module.ID},
						models.USR_Feature{Name: "Delete Offer", ModuleID: module.ID},
						models.USR_Feature{Name: "Validate Offer", ModuleID: module.ID},
					)
				case "Invoice":
					features = append(
						features,
						models.USR_Feature{Name: "View Invoice", ModuleID: module.ID},
						models.USR_Feature{Name: "Create Invoice", ModuleID: module.ID},
						models.USR_Feature{Name: "Update Invoice", ModuleID: module.ID},
						models.USR_Feature{Name: "Delete Invoice", ModuleID: module.ID},
					)
				case "Bill":
					features = append(
						features,
						models.USR_Feature{Name: "View Bill", ModuleID: module.ID},
						models.USR_Feature{Name: "Create Bill", ModuleID: module.ID},
						models.USR_Feature{Name: "Update Bill", ModuleID: module.ID},
						models.USR_Feature{Name: "Delete Bill", ModuleID: module.ID},
						models.USR_Feature{Name: "Pay Bill", ModuleID: module.ID},
					)
				case "Review":
					features = append(
						features,
						models.USR_Feature{Name: "View Review", ModuleID: module.ID},
						models.USR_Feature{Name: "Create Review", ModuleID: module.ID},
						models.USR_Feature{Name: "Update Review", ModuleID: module.ID},
						models.USR_Feature{Name: "Delete Review", ModuleID: module.ID},
					)
				case "Contract":
					features = append(
						features,
						models.USR_Feature{Name: "View Contract", ModuleID: module.ID},
						models.USR_Feature{Name: "Create Contract", ModuleID: module.ID},
						models.USR_Feature{Name: "Update Contract", ModuleID: module.ID},
						models.USR_Feature{Name: "Delete Contract", ModuleID: module.ID},
					)
				case "Contract Monitoring":
					features = append(
						features,
						models.USR_Feature{Name: "View Contract Monitoring", ModuleID: module.ID},
						models.USR_Feature{Name: "Create Contract Monitoring", ModuleID: module.ID},
						models.USR_Feature{Name: "Update Contract Monitoring", ModuleID: module.ID},
						models.USR_Feature{Name: "Delete Contract Monitoring", ModuleID: module.ID},
					)
				case "Addendum":
					features = append(
						features,
						models.USR_Feature{Name: "View Addendum", ModuleID: module.ID},
						models.USR_Feature{Name: "Create Addendum", ModuleID: module.ID},
						models.USR_Feature{Name: "Update Addendum", ModuleID: module.ID},
						models.USR_Feature{Name: "Delete Addendum", ModuleID: module.ID},
					)
				case "SPMK":
					features = append(
						features,
						models.USR_Feature{Name: "View SPMK", ModuleID: module.ID},
						models.USR_Feature{Name: "Create SPMK", ModuleID: module.ID},
						models.USR_Feature{Name: "Update SPMK", ModuleID: module.ID},
						models.USR_Feature{Name: "Delete SPMK", ModuleID: module.ID},
					)
				}
			}

			if err := db.Create(&features).Error; err != nil {
				log.Fatalf("Error seeding USR_Feature: %v", err)
			} else {
				log.Println("USR_Feature seeded successfully.")
			}
		}
	}

	// Seed Roles
	var roles []models.USR_Role
	{
		var roleCount int64
		db.Model(&models.USR_Role{}).Count(&roleCount)

		// Convert features to slice of pointers
		var featurePtrs []*models.USR_Feature
		for i := range features {
			featurePtrs = append(featurePtrs, &features[i])
		}

		if roleCount == 0 {
			// Inritialize the map of feature for vendor
			vendorFeature := map[string]struct{}{
				"View Vendor Profile":   {},
				"Update Vendor Profile": {},
				"View Procurement":      {},
				"View Interest":         {},
				"Create Interest":       {},
				"Update Interest":       {},
				"Delete Interest":       {},
				"View Offer":            {},
				"Create Offer":          {},
				"Update Offer":          {},
				"Delete Offer":          {},
				"View Contract":         {},
				"View Addendum":         {},
				"View SPMK":             {},
				"View Invoice":          {},
				"Create Invoice":        {},
				"Update Invoice":        {},
				"Delete Invoice":        {},
			}

			// Role for vendor
			vendorRole := models.USR_Role{
				Name:             "Vendor",
				IsAdministrative: false,
				Features:         []*models.USR_Feature{},
			}
			for _, feature := range featurePtrs {
				if _, exist := vendorFeature[feature.Name]; exist {
					vendorRole.Features = append(vendorRole.Features, feature)
				}
			}

			roles = append(roles, models.USR_Role{Name: "Admin", IsAdministrative: true, Features: featurePtrs})
			roles = append(roles, vendorRole)

			if err := db.Create(&roles).Error; err != nil {
				log.Fatalf("Error seeding USR_Role: %v", err)
			} else {
				log.Println("USR_Roles seeded successfully.")
			}
		}
	}

	// Seed Users
	var userCount int64
	db.Model(&models.USR_User{}).Count(&userCount)
	if userCount == 0 {
		user := models.USR_User{Name: "Admin", Username: "admin", Email: "admin@jxboard.id"}

		// Search for admin role
		for _, role := range roles {
			if role.Name == "Admin" {
				user.RoleID = role.ID
			}
		}

		// hashed env password
		hashedPassword, hashErr := bcrypt.GenerateFromPassword([]byte(helpers.GetENVWithDefault("ADMIN_PASS", "|(-=-)|")), bcrypt.DefaultCost)
		if hashErr == nil {
			user.Password = string(hashedPassword)

			// Create User
			if err := db.Create(&user).Error; err != nil {
				log.Fatalf("Error seeding USR_User: %v", err)
			} else {
				log.Println("USR_User seeded successfully.")
			}
		} else {
			log.Fatalln("Error seeding USR_User hash failed")
		}
	}
}
