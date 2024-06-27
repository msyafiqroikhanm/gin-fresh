package helpers

import (
	"jxb-eprocurement/handlers/dtos"
	"jxb-eprocurement/models"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func GenerateJWT(user models.USR_User) (string, error) {
	// Getting jwt related data from env
	jwtKey := []byte(os.Getenv("JWT_SECRET"))
	jwtLimitStr := GetENVWithDefault("JWT_TIME", "24")

	jwtLimitInt, err := strconv.Atoi(jwtLimitStr)
	if err != nil {
		jwtLimitInt = 24
	}

	// Set expiration time of jwt token
	expirationTime := time.Now().Add(time.Duration(jwtLimitInt) * time.Hour)

	// Create jwt claim to generate jwt token
	claims := &dtos.Claims{
		UserID: user.ID,
		User:   user.Name,
		RoleID: user.Role.ID,
		Role:   user.Role.Name,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
			NotBefore: time.Now().Unix(),
		},
	}

	// Generate jwt token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
