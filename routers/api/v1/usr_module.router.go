package v1

import (
	"jxb-eprocurement/controllers"
	"jxb-eprocurement/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitModuleRoutes(r *gin.RouterGroup, db *gorm.DB) {
	// Setup controller and route
	moduleController := controllers.ModuleControllerConstructor(service.ModuleServiceConstructor(db))
	moduleRoutes := r.Group("/module")

	// Additional middleware to implement to the group routes
	// moduleRoutes.Use(middlewares.AuthMiddleware()) // Uncomment this when the user module and feature module is finish

	// Collection of routes
	{
		moduleRoutes.GET("", moduleController.GetAllModules)
		moduleRoutes.GET("/:id", moduleController.GetModule)
		moduleRoutes.POST("", moduleController.CreateModule)
		moduleRoutes.PUT("/:id", moduleController.UpdateModule)
		moduleRoutes.DELETE("/:id", moduleController.DeleteModule)
	}
}
