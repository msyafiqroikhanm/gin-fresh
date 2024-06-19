package v1

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitRoutes(r *gin.RouterGroup, db *gorm.DB) {
	InitModuleRoutes(r, db)
}
