package api

import (
	v1 "jxb-eprocurement/routers/api/v1"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitRoutes(r *gin.Engine, db *gorm.DB) {
	apiRoutes := r.Group("/api")
	v1.InitRoutes(apiRoutes, db)
}
