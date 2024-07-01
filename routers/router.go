package routers

import (
	"jxb-eprocurement/handlers"
	"jxb-eprocurement/middlewares"
	apis "jxb-eprocurement/routers/api"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	router := gin.Default()

	// Apply middlewares here
	router.Use(middlewares.RequestIDMiddleware())
	router.Use(handlers.APILogger())

	// Initialize route groups for versioning
	apis.InitRoutes(router, db)

	return router
}
