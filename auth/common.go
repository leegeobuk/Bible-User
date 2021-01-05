package auth

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
)

var (
	origins = []string{
		"http://localhost:3000",
		"https://www.biblennium.com",
	}
	// ErrAccountExist returns error when account exists in db
	ErrAccountExist = errors.New("error account already exists")
	// ErrAccountNotExist returns error when account is not in db
	ErrAccountNotExist = errors.New("error account doesn't exist")
	// ErrEmptyCookie returns error when cookie is not in request headers
	ErrEmptyCookie = errors.New("error empty cookie from request")
)

// Response returns default response with headers and status code
func Response(request *events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	corsHeaders := map[string]string{
		"Access-Control-Allow-Credentials": "true",
		"Access-Control-Allow-Headers":     "Content-Type",
		"Access-Control-Allow-Origin":      "",
	}
	origin := request.Headers["origin"]
	setOrigin(corsHeaders, origin)
	return events.APIGatewayProxyResponse{
		Headers:           corsHeaders,
		MultiValueHeaders: map[string][]string{},
		StatusCode:        http.StatusInternalServerError,
	}
}

func setOrigin(headers map[string]string, origin string) {
	headers["Access-Control-Allow-Origin"] = getOrigin(origin)
}

func getOrigin(origin string) string {
	for _, o := range origins {
		if o == origin {
			return o
		}
	}
	return origins[0]
}

// CreateRefreshCookie returns cookie with given value expiring after given duration
func CreateRefreshCookie(value string, dur time.Duration) *http.Cookie {
	return &http.Cookie{
		Name:     "refresh_token",
		Value:    value,
		Expires:  time.Now().Local().Add(dur),
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
		HttpOnly: true,
	}
}

// SetCookie sets given cookie to given headers
func SetCookie(h map[string][]string, c *http.Cookie) {
	h["Set-Cookie"] = append(h["Set-Cookie"], c.String())
}

// ParseCookie returns the value of cookie from client
func ParseCookie(cookieString string) string {
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
