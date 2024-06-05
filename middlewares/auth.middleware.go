package middlewares

import (
	"jxb-eprocurement/handlers"
	"jxb-eprocurement/models"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			handlers.ResponseFormatter(c, http.StatusUnauthorized, nil, "Authorization token not provided")
			c.Abort()
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		claims := &handlers.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return handlers.JwtKey, nil
		})

		if err != nil || !token.Valid {
			handlers.ResponseFormatter(c, http.StatusUnauthorized, nil, "Invalid authorization token")
			c.Abort()
			return
		}

		var user models.User
		if err := models.DB.Where("id = ?", claims.UserID).Preload("Role").First(&user).Error; err != nil {
			handlers.ResponseFormatter(c, http.StatusUnauthorized, nil, "Invalid user")
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Set("role", user.Role.Name)
		c.Next()
	}
}

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role != "admin" {
			handlers.ResponseFormatter(c, http.StatusForbidden, nil, "Admin access required")
			c.Abort()
			return
		}
		c.Next()
	}
}
