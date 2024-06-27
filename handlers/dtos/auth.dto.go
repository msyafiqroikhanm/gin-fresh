package dtos

import "github.com/dgrijalva/jwt-go"

type (
	Claims struct {
		UserID           uint     `json:"user_id"`
		RoleID           uint     `json:"role_id"`
		User             string   `json:"user"`
		Role             string   `json:"role"`
		IsAdministrative bool     `json:"is_administrative"`
		Features         []string `json:"features"`
		jwt.StandardClaims
	}

	InputLoginDTO struct {
		UsernameOrEmail string `json:"username_or_email" form:"username_or_email" validate:"required,no_space"`
		Password        string `json:"password" form:"password" validate:"required"`
	}
)
