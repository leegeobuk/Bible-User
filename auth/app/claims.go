package app

import "github.com/dgrijalva/jwt-go"

type claims struct {
	userID string
	jwt.StandardClaims
}
