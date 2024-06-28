package accesses

import (
	"jxb-eprocurement/controllers"
	"jxb-eprocurement/middlewares"
	"jxb-eprocurement/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitModuleRoutes(r *gin.RouterGroup, db *gorm.DB) {
	// Setup controller and route
	moduleController := controllers.ModuleControllerConstructor(service.ModuleServiceConstructor(db))
	moduleRoutes := r.Group("/modules")

	// Additional middleware to implement to the group routes
	moduleRoutes.Use(middlewares.Authentication())

	// Collection of routes
	{
		// Get All
		moduleRoutes.GET(
			"",
			middlewares.Authorization(
				[]string{
					"View Module",
					"Create Module",
					"Update Module",
					"Delete Module",
				},
				false,
			),
			moduleController.GetAllModules,
		)

		// Get Detail
		moduleRoutes.GET(
			"/:id",
			middlewares.Authorization(
				[]string{
					"View Module",
					"Create Module",
					"Update Module",
					"Delete Module",
				},
				false,
			),
			moduleController.GetModule,
		)

		// Create
		moduleRoutes.POST(
			"",
			middlewares.Authorization(
				[]string{"Create Module"},
				false,
			),
			moduleController.CreateModule,
		)

		// Update
		moduleRoutes.PUT(
			"/:id",
			middlewares.Authorization(
				[]string{"Update Module"},
				false,
			),
			moduleController.UpdateModule,
		)

		// Delete
		moduleRoutes.DELETE(
			"/:id",
			middlewares.Authorization(
				[]string{"Delete Module"},
				false,
			),
			moduleController.DeleteModule,
		)
	}
}
