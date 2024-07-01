package v1

import (
	"jxb-eprocurement/routers/api/v1/accesses"
	"jxb-eprocurement/routers/api/v1/authentication"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitRoutes(r *gin.RouterGroup, db *gorm.DB) {
	v1Routes := r.Group("/v1")

	accesses.InitAccessRoutes(v1Routes, db)
	authentication.InitAuthRoutes(v1Routes, db)
}
