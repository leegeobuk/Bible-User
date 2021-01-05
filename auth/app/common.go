package app

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	accessSignKey  = os.Getenv("ACCESS_SIGN_KEY")
	refreshSignKey = os.Getenv("REFRESH_SIGN_KEY")
)

// generateAccessToken generates access token expiring after given duration hours
func generateAccessToken(uid string, dur time.Duration) (string, error) {
	// jwt.StandardClaims.ExpiresAt takes unix time
	c := &claims{uid, jwt.StandardClaims{ExpiresAt: time.Now().Local().Add(dur).Unix()}}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return token.SignedString([]byte(accessSignKey))
}

// generateRefreshToken generates refresh token expiring after given duration days
func generateRefreshToken(uid string, dur time.Duration) (string, error) {
	c := &claims{uid, jwt.StandardClaims{ExpiresAt: time.Now().Local().Add(dur).Unix()}}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return token.SignedString([]byte(refreshSignKey))
}
