package v1

import (
	"jxb-eprocurement/controllers"
	"jxb-eprocurement/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitModuleRoutes(r *gin.RouterGroup, db *gorm.DB) {
	moduleController := controllers.NewModuleController(service.NewModuleService(db))
	r.GET("/module", moduleController.GetAllModules)
	r.GET("/module/:id", moduleController.GetModule)
	r.POST("/module", moduleController.CreateModule)
	r.PUT("/module/:id", moduleController.UpdateModule)
	r.DELETE("/module/:id", moduleController.DeleteModule)
}
