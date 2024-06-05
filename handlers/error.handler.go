package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorHandler(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			// Cetak log kesalahan
			log.Printf("Unhandled error: %v", r)
			ResponseFormatter(c, http.StatusInternalServerError, nil, "Internal Server Error")
		}
	}()
	c.Next()
}

func NotFoundHandler(c *gin.Context) {
	ResponseFormatter(c, http.StatusNotFound, nil, "Resource not found")
	// c.JSON(http.StatusNotFound, gin.H{
	// 	"status":  http.StatusNotFound,
	// 	"message": "Resource not found",
	// })
}
