package accesses

import (
	"jxb-eprocurement/controllers"
	"jxb-eprocurement/middlewares"
	"jxb-eprocurement/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitRoleRoutes(r *gin.RouterGroup, db *gorm.DB) {
	// Setup controller and route
	roleController := controllers.RoleControllerConstructor(service.RoleServiceConstructor(db))
	roleRoutes := r.Group("/roles")

	// Additional middleware to implement to the group routes
	roleRoutes.Use(middlewares.Authentication()) // Uncomment this when the user module and feature module is finish

	// Collection of routes
	{
		// Get All
		roleRoutes.GET(
			"",
			middlewares.Authorization(
				[]string{
					"View Role",
					"Create Role",
					"Update Role",
					"Delete Role",
				},
				false,
			),
			roleController.GetAllRoles,
		)

		// Get Detail
		roleRoutes.GET(
			"/:id",
			middlewares.Authorization(
				[]string{
					"View Role",
					"Create Role",
					"Update Role",
					"Delete Role",
				},
				false,
			),
			roleController.GetRole,
		)

		// Create
		roleRoutes.POST(
			"",
			middlewares.Authorization(
				[]string{"Create Role"},
				false,
			),
			roleController.CreateRole,
		)

		// Update
		roleRoutes.PUT(
			"/:id",
			middlewares.Authorization(
				[]string{"Update Role"},
				false,
			),
			roleController.UpdateRole,
		)

		// Delete
		roleRoutes.DELETE(
			"/:id", middlewares.Authorization(
				[]string{"Create Role"},
				false,
			),
			roleController.DeleteRole,
		)
	}
}
