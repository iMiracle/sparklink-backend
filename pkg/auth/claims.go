package auth

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	UserID uint   `json:"user_id"`
	Phone  string `json:"phone"`
	jwt.RegisteredClaims
}
