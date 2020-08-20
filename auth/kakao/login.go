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
	resp := events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}

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

	// return token if account is a member, unauthorize  if not
	if isMember(user, db) {
		// copy headers due to adding cookie later
		resp.Headers = copyHeaders(headers)
		resp.StatusCode = http.StatusOK

		// marshal kakaoLoginResponse
		loginResp := &kakaoLoginResponse{AccessToken: kakaoToken.AccessToken, ExpiresIn: kakaoToken.ExpiresIn - 21290}
		data, err := json.Marshal(loginResp)
		if err != nil {
			return resp, err
		}
		resp.Body = string(data)

		// set httpOnly cookie
		cookie := &http.Cookie{Name: "refresh_token", Value: kakaoToken.RefreshToken, HttpOnly: true}
		setCookie(resp.Headers, cookie)

		// add refresh_token to the user
		db.Model(user).UpdateColumn("refresh_token", kakaoToken.RefreshToken)

		return resp, nil
	}
	resp.StatusCode = http.StatusUnauthorized

	return resp, nil
}

func setCookie(h map[string]string, c *http.Cookie) {
	cookieString := ""
	if c.HttpOnly {
		cookieString = fmt.Sprintf("%s=%s; HttpOnly", c.Name, c.Value)
	} else {
		cookieString = fmt.Sprintf("%s=%s;", c.Name, c.Value)
	}
	h["Set-Cookie"] = cookieString
}
