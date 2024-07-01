package middlewares

import (
	"jxb-eprocurement/helpers"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CORS middleware, setup routes to use CORS Policy
func CORS() gin.HandlerFunc {
	// Setup allowed origins slice
	var origins []string
	originENV := helpers.GetENVWithDefault("FRONTEND_URLS", "http://localhost:3000")
	origins = strings.Split(originENV, ",")

	// Setup cors maxAge
	maxAgeENV := helpers.GetENVWithDefault("CORS_MAX_AGE", "6")
	maxAge, err := strconv.Atoi(maxAgeENV)
	if err != nil {
		maxAge = 6
	}

	// Setup config for cors
	config := cors.Config{
		AllowOrigins: origins,
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Accept",
			"Authorization",
			"Accept-Encoding",
			"Cache-Control",
			"Connection",
			"Content-Length",
			"User-Agent",
			"Host",
		},
		AllowCredentials: true,
		MaxAge:           time.Duration(maxAge) * time.Hour,
	}
	return cors.New(config)
}
