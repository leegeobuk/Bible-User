package auth

import (
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	origin = "http://localhost:3000"
)

var (
	corsHeaders = map[string]string{
		"Access-Control-Allow-Credentials": "true",
		"Access-Control-Allow-Headers":     "Content-Type",
		"Access-Control-Allow-Origin":      origin,
	}
	errAccountExist    = errors.New("error account already exists")
	errAccountNotExist = errors.New("error account doesn't exist")
	errEmptyCookie     = errors.New("error empty cookie from request")
	accessSignKey      = os.Getenv("ACCESS_SIGN_KEY")
	refreshSignKey     = os.Getenv("REFRESH_SIGN_KEY")
)

func addHeaders(m map[string]string, h map[string]string) {
	for k, v := range h {
		m[k] = v
	}
}

func createRefreshCookie(value string, dur time.Duration) *http.Cookie {
	return &http.Cookie{
		Name:     "refresh_token",
		Value:    value,
		Expires:  time.Now().Local().Add(dur),
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
		HttpOnly: true,
	}
}

func setCookie(h map[string][]string, c *http.Cookie) {
	h["Set-Cookie"] = append(h["Set-Cookie"], c.String())
}

func parseCookie(cookieString string) string {
	cookies := strings.Split(cookieString, "; ")
	var c string
	for _, v := range cookies {
		if strings.HasPrefix(v, "refresh_token") {
			i := strings.Index(v, "=")
			c = v[i+1:]
		}
	}
	return c
}

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
