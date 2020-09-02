package kakao

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
)

const kakaoBaseURL = "https://kauth.kakao.com/oauth/logout"

// Logout logs out kakao user and removes refresh_token cookie
func Logout(ctx context.Context, request *events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	resp := events.APIGatewayProxyResponse{}

	// check if refresh_token is stored in cookie
	cookieString, ok := request.Headers["Cookie"]
	fmt.Println(cookieString)
	if !ok {
		resp.StatusCode = http.StatusBadRequest
		resp.Body = errEmptyCookie.Error()
		return resp, nil
	}

	refreshToken := parseCookie(cookieString)

	cookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  time.Now().Local().Add(-24 * time.Hour),
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
		HttpOnly: true,
	}

	resp.Headers = copyHeaders(headers)
	setCookie(resp.Headers, cookie)
	resp.StatusCode = http.StatusOK

	return resp, nil
}
