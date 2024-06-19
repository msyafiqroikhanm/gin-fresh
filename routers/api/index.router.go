package api

import (
	v1 "jxb-eprocurement/routers/api/v1"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitRoutes(r *gin.Engine, db *gorm.DB) {
	apiV1 := r.Group("/api/v1")
	v1.InitRoutes(apiV1, db)
}
