package accesses

import (
	"jxb-eprocurement/controllers"
	"jxb-eprocurement/middlewares"
	"jxb-eprocurement/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitUserRoutes(r *gin.RouterGroup, db *gorm.DB) {
	// Setup controller and route
	userController := controllers.UserControllerConstructor(service.UserServiceConstructor(db))
	userRoutes := r.Group("/users")

	// Additional middleware to implement to the group routes
	userRoutes.Use(middlewares.Authentication())

	// Collection of routes
	{
		// Get All
		userRoutes.GET(
			"",
			middlewares.Authorization(
				[]string{
					"View User",
					"Create User",
					"Update User",
					"Delete User",
					"Reset User Password",
				},
				false,
			),
			userController.GetAllUsers,
		)

		// Get Detail
		userRoutes.GET(
			"/:id",
			middlewares.Authorization(
				[]string{
					"View User",
					"Create User",
					"Update User",
					"Delete User",
				},
				false,
			),
			userController.GetUser,
		)

		// Create
		userRoutes.POST(
			"",
			middlewares.Authorization([]string{"Create User"}, false),
			userController.CreateUser,
		)

		// Edit
		userRoutes.PUT(
			"/:id",
			middlewares.Authorization([]string{"Update User"}, false),
			userController.UpdateUser,
		)

		// Delete
		userRoutes.DELETE(
			"/:id",
			middlewares.Authorization([]string{"Delete User"}, false),
			userController.DeleteUser,
		)

		// Reset Password
		userRoutes.PATCH(
			"/reset-pass/:id",
			middlewares.Authorization([]string{"Reset User Password"}, true),
			userController.ResetPassUser,
		)

		// Change Password
		userRoutes.PATCH(
			"/change-pass/:id",
			middlewares.Authorization([]string{"Reset User Password"}, false),
			userController.ChangePassUser,
		)
	}
}
