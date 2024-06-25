package accesses

import (
	"jxb-eprocurement/controllers"
	"jxb-eprocurement/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitRoleRoutes(r *gin.RouterGroup, db *gorm.DB) {
	// Setup controller and route
	roleController := controllers.RoleControllerConstructor(service.RoleServiceConstructor(db))
	roleRoutes := r.Group("/roles")

	// Additional middleware to implement to the group routes
	// roleRoutes.Use(middlewares.AuthMiddleware()) // Uncomment this when the user module and feature module is finish

	// Collection of routes
	{
		roleRoutes.GET("", roleController.GetAllRoles)
		roleRoutes.GET("/:id", roleController.GetRole)
		roleRoutes.POST("", roleController.CreateRole)
		roleRoutes.PUT("/:id", roleController.UpdateRole)
		roleRoutes.DELETE("/:id", roleController.DeleteRole)
	}
}
