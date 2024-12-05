package entity

import "github.com/golang-jwt/jwt/v5"

type JWTCustomClaims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}