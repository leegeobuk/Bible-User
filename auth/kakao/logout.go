package kakao

import (
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
)

const kakaoBaseURL = "https://kauth.kakao.com/oauth/logout"

// Logout logs out kakao user and removes refresh_token cookie
func Logout(request *events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	resp := events.APIGatewayProxyResponse{Headers: map[string]string{}, MultiValueHeaders: map[string][]string{}}

	// check if refresh_token is stored in cookie
	cookieString, ok := request.Headers["Cookie"]
	if !ok {
		resp.Body = errEmptyCookie.Error()
		resp.StatusCode = http.StatusBadRequest
		return resp
	}

	// expire the refresh_token stored in cookie
	refreshToken := parseCookie(cookieString)
	cookie := createRefreshCookie(refreshToken, int(-24*time.Hour.Seconds()))
	setCookie(resp.MultiValueHeaders, cookie)

	resp.StatusCode = http.StatusOK

	return resp
}
