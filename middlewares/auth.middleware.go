package middlewares

import (
	"jxb-eprocurement/handlers"
	"jxb-eprocurement/handlers/dtos"
	"jxb-eprocurement/models"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// Middleware to check user jwt is valid and correct
func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			handlers.ResponseFormatter(c, http.StatusUnauthorized, nil, "Authorization token not provided")
			c.Abort()
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		claims := &dtos.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) { return []byte(os.Getenv("JWT_SECRET")), nil })

		if err != nil || !token.Valid {
			if err == jwt.ErrSignatureInvalid {
				handlers.ResponseFormatter(c, http.StatusUnauthorized, nil, "Invalid token")
				c.Abort()
				return
			}

			validationError, ok := err.(*jwt.ValidationError)
			if ok && validationError.Errors == jwt.ValidationErrorExpired {
				handlers.ResponseFormatter(c, http.StatusUnauthorized, nil, "Token has expired")
				c.Abort()
				return
			}

			if ok && validationError.Errors == jwt.ValidationErrorNotValidYet {
				handlers.ResponseFormatter(c, http.StatusUnauthorized, nil, "Token is not valid yet")
				c.Abort()
				return
			}

			handlers.ResponseFormatter(c, http.StatusUnauthorized, nil, "Invalid token")
			c.Abort()
			return
		}

		// Parse claims data to context for further access authorization
		c.Set("user", &models.USR_User{ID: claims.UserID, Name: claims.User, RoleID: claims.RoleID})
		c.Set("role", &models.USR_Role{ID: claims.RoleID, Name: claims.Role, IsAdministrative: claims.IsAdministrative})
		c.Set("features", claims.Features)

		c.Next()
	}
}

// Middleware to check user have access (feature) to access the endpoint
// TODO: After Vendor Module Finish Development, Adding Check if Vendor Is Validated Or Not
func Authorization(allowedFeatures []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		features, exist := c.Get("features")
		if !exist {
			handlers.ResponseFormatter(c, http.StatusForbidden, nil, "Unauthorized to access this resource")
			c.Abort()
			return
		}

		// Assert that features is a slice of strings
		featureList, ok := features.([]string)
		if !ok {
			handlers.ResponseFormatter(c, http.StatusForbidden, nil, "Unauthorized to access this resource")
			c.Abort()
			return
		}

		// Convert allowed features to a map for efficient lookup
		allowedFeaturesMap := make(map[string]struct{})
		for _, feature := range allowedFeatures {
			allowedFeaturesMap[feature] = struct{}{}
		}

		// Check if user have allowed list
		for _, feature := range featureList {
			if _, exists := allowedFeaturesMap[feature]; exists {
				c.Next()
				return
			}
		}

		handlers.ResponseFormatter(c, http.StatusForbidden, nil, "Unauthorized to access this resource")
		c.Abort()
	}
}
