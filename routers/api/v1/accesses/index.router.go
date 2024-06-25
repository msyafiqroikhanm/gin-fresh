package accesses

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitAccessRoutes(r *gin.RouterGroup, db *gorm.DB) {
	accessRoutes := r.Group("/accesses")

	InitModuleRoutes(accessRoutes, db)
	InitFeatureRoutes(accessRoutes, db)
	InitRoleRoutes(accessRoutes, db)
}
