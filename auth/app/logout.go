package app

import (
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/leegeobuk/Bible-User/auth"
)

// Logout logs out users
func Logout(request *events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	resp := auth.Response(request)

	// check if refresh_token is stored in cookie
	cookieString, ok := request.Headers["Cookie"]
	if !ok {
		resp.Body = auth.ErrEmptyCookie.Error()
		resp.StatusCode = http.StatusBadRequest
		return resp
	}

	// expire the refresh_token stored in cookie
	refreshToken := auth.ParseCookie(cookieString)
	cookie := auth.CreateRefreshCookie(refreshToken, -24*time.Hour)
	auth.SetCookie(resp.MultiValueHeaders, cookie)

	resp.StatusCode = http.StatusOK

	return resp
}
