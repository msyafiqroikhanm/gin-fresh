package routers

import (
	apis "jxb-eprocurement/routers/api"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	router := gin.Default()

	// Initialize route groups for versioning
	apis.InitRoutes(router, db)

	return router
}
