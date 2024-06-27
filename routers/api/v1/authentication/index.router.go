package authentication

import (
	"jxb-eprocurement/controllers"
	"jxb-eprocurement/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitAuthRoutes(r *gin.RouterGroup, db *gorm.DB) {
	authController := controllers.AuthControllerConstructor(service.AuthServiceConstructor(db))
	authRoutes := r.Group("/auth")

	{
		authRoutes.POST("/login", authController.LoginUser)
	}
}
