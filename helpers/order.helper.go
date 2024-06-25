package helpers

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Order(c *gin.Context, allowedOrderFields []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		// Ordering logic
		orderBy := c.DefaultQuery("order_by", "")
		order := c.DefaultQuery("order", "asc")

		// Validate the order_by field
		isValidOrderField := false
		for _, field := range allowedOrderFields {
			if field == orderBy {
				isValidOrderField = true
				break
			}
		}

		if isValidOrderField {
			if order != "asc" && order != "desc" {
				order = "asc"
			}
			db = db.Order(fmt.Sprintf("%s %s", orderBy, order))
		}

		return db
	}
}
