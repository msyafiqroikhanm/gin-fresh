package accesses

import (
	"jxb-eprocurement/controllers"
	"jxb-eprocurement/middlewares"
	"jxb-eprocurement/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitFeatureRoutes(r *gin.RouterGroup, db *gorm.DB) {
	// Setup controller and route
	featureController := controllers.FeatureControllerConstructor(service.FeatureServiceConstructor(db))
	moduleRoutes := r.Group("/features")

	// Additional middleware to implement to the group routes
	moduleRoutes.Use(middlewares.Authentication()) // Uncomment this when the user module and feature module is finish

	// Collection of routes
	{
		// Get All
		moduleRoutes.GET(
			"",
			middlewares.Authorization(
				[]string{
					"View Feature",
					"Create Feature",
					"Update Feature",
					"Delete Feature",
				},
				false,
			),
			featureController.GetAllFeatures)

		// Get Detail
		moduleRoutes.GET(
			"/:id",
			middlewares.Authorization(
				[]string{
					"View Feature",
					"Create Feature",
					"Update Feature",
					"Delete Feature",
				},
				false,
			),
			featureController.GetFeature,
		)

		// Create
		moduleRoutes.POST(
			"",
			middlewares.Authorization(
				[]string{"Create Feature"},
				false,
			),
			featureController.CreateFeature,
		)

		// Update
		moduleRoutes.PUT(
			"/:id",
			middlewares.Authorization(
				[]string{"Update Feature"},
				false,
			),
			featureController.UpdateFeature,
		)

		// Delete
		moduleRoutes.DELETE(
			"/:id",
			middlewares.Authorization(
				[]string{"Delete Feature"},
				false,
			),
			featureController.DeleteFeature,
		)
	}
}
