package accesses

import (
	"jxb-eprocurement/controllers"
	"jxb-eprocurement/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitFeatureRoutes(r *gin.RouterGroup, db *gorm.DB) {
	// Setup controller and route
	featureController := controllers.FeatureControllerConstructor(service.FeatureServiceConstructor(db))
	moduleRoutes := r.Group("/features")

	// Additional middleware to implement to the group routes
	// moduleRoutes.Use(middlewares.AuthMiddleware()) // Uncomment this when the user module and feature module is finish

	// Collection of routes
	{
		moduleRoutes.GET("", featureController.GetAllFeatures)
		moduleRoutes.GET("/:id", featureController.GetFeature)
		moduleRoutes.POST("", featureController.CreateFeature)
		moduleRoutes.PUT("/:id", featureController.UpdateFeature)
		moduleRoutes.DELETE("/:id", featureController.DeleteFeature)
	}
}
