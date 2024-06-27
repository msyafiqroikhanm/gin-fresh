package accesses

import (
	"jxb-eprocurement/controllers"
	"jxb-eprocurement/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitUserRoutes(r *gin.RouterGroup, db *gorm.DB) {
	// Setup controller and route
	userController := controllers.UserControllerConstructor(service.UserServiceConstructor(db))
	userRoutes := r.Group("/users")

	// Additional middleware to implement to the group routes
	// userRoutes.Use(middlewares.AuthMiddleware()) // Uncomment this when the user module and feature module is finish

	// Collection of routes
	{
		userRoutes.GET("", userController.GetAllUsers)
		userRoutes.GET("/:id", userController.GetUser)
		userRoutes.POST("", userController.CreateUser)
		userRoutes.PUT("/:id", userController.UpdateUser)
		userRoutes.DELETE("/:id", userController.DeleteUser)
		userRoutes.PATCH("/reset-pass/:id", userController.ResetPassUser)
		userRoutes.PATCH("/change-pass/:id", userController.ChangePassUser)
	}
}
