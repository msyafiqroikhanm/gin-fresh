package handlers

import (
	"jxb-eprocurement/models"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var JwtKey = []byte(os.Getenv("JWT_SECRET"))

type Claims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.StandardClaims
}

func GenerateJWT(user models.User) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID: user.ID,
		Role:   user.Role.Name,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JwtKey)
}

func Authenticate(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		ResponseFormatter(c, http.StatusBadRequest, nil, "Invalid input")
		return
	}

	var user models.User
	if err := models.DB.Where("email = ?", input.Email).Preload("Role").First(&user).Error; err != nil {
		ResponseFormatter(c, http.StatusUnauthorized, nil, "Invalid email or password")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		ResponseFormatter(c, http.StatusUnauthorized, nil, "Invalid email or password")
		return
	}

	token, err := GenerateJWT(user)
	if err != nil {
		ResponseFormatter(c, http.StatusInternalServerError, nil, "Failed to generate token")
		return
	}

	ResponseFormatter(c, http.StatusOK, gin.H{"token": token}, "Login successful")
}
