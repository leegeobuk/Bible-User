package kakao

import (
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
)

const kakaoBaseURL = "https://kauth.kakao.com/oauth/logout"

// Logout logs out kakao user and removes refresh_token cookie
func Logout(request *events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	resp := events.APIGatewayProxyResponse{Headers: map[string]string{}}

	// check if refresh_token is stored in cookie
	cookieString, ok := request.Headers["Cookie"]
	if !ok {
		resp.Body = errEmptyCookie.Error()
		resp.StatusCode = http.StatusBadRequest
		return resp
	}

	// expire the refresh_token stored in cookie
	refreshToken := parseCookie(cookieString)
	cookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  time.Now().Local().Add(-24 * time.Hour),
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
		HttpOnly: true,
	}
	setCookie(resp.Headers, cookie)

	resp.StatusCode = http.StatusOK

	return resp
}
