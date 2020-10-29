package auth

import (
	"context"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/leegeobuk/Bible-User/auth/kakao"
)

// Logout logs out users
func Logout(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if request.QueryStringParameters["type"] == "kakao" {
		resp := kakao.Logout(&request)
		addHeaders(resp.Headers, corsHeaders)
		return resp, nil
	}

	resp := events.APIGatewayProxyResponse{Headers: corsHeaders, MultiValueHeaders: map[string][]string{}}

	// check if refresh_token is stored in cookie
	cookieString, ok := request.Headers["Cookie"]
	if !ok {
		resp.Body = errEmptyCookie.Error()
		resp.StatusCode = http.StatusBadRequest
		return resp, nil
	}

	// expire the refresh_token stored in cookie
	refreshToken := parseCookie(cookieString)
	cookie := createRefreshCookie(refreshToken, -24*time.Hour)
	setCookie(resp.MultiValueHeaders, cookie)

	resp.StatusCode = http.StatusOK

	return resp, nil
}
