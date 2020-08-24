package kakao

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

// Login authenticates kakao user and decide whether to login or not
func Login(ctx context.Context, request *events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	resp := events.APIGatewayProxyResponse{Headers: headers, StatusCode: http.StatusInternalServerError}

	// get token from Kakao Login API
	kakaoToken, err := getToken(request)
	if err != nil {
		return resp, err
	}

	// request member info from Kakao logic
	kakaoUserResp, err := getUserInfo(kakaoToken)
	if err != nil {
		return resp, err
	}

	// finding account in db logic
	db, err := connectDB()
	defer db.Close()
	if err != nil {
		return resp, err
	}

	user := kakaoUserResp.toKakaoUser()

	// unauthenticate if not a member
	if !isMember(user, db) {
		resp.StatusCode = http.StatusUnauthorized
		return resp, nil
	}

	// copy headers to set cookie later
	resp.Headers = copyHeaders(headers)
	resp.Headers["Access-Control-Expose-Headers"] = "Set-Cookie"
	fmt.Println("resp.headers", resp.Headers)

	// marshal kakaoLoginResponse
	loginResp := &kakaoTokenDTO{AccessToken: kakaoToken.AccessToken, ExpiresIn: kakaoToken.ExpiresIn - 21595}
	data, err := json.Marshal(loginResp)
	if err != nil {
		return resp, err
	}
	resp.Body = string(data)

	// set httpOnly cookie
	cookie := createRefreshCookie(kakaoToken.RefreshToken, kakaoToken.RefreshTokenExpiresIn)
	setCookie(resp.Headers, cookie)

	// add refresh_token to the user
	db.Model(user).UpdateColumn("refresh_token", kakaoToken.RefreshToken)

	resp.StatusCode = http.StatusOK

	return resp, nil
}
