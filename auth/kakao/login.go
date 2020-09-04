package kakao

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/leegeobuk/Bible-User/db"
	"github.com/leegeobuk/Bible-User/util"
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
	database, err := db.ConnectDB()
	defer database.Close()
	if err != nil {
		return resp, err
	}

	user := kakaoUserResp.toKakaoUser()

	// unauthenticate if not a member
	if !db.IsMember(database, user) {
		resp.StatusCode = http.StatusUnauthorized
		return resp, nil
	}

	// marshal kakaoLoginResponse
	loginResp := &kakaoTokenDTO{AccessToken: kakaoToken.AccessToken, ExpiresIn: kakaoToken.ExpiresIn}
	data, err := json.Marshal(loginResp)
	if err != nil {
		return resp, err
	}

	// set httpOnly cookie
	resp.Headers = util.CopyMap(headers)
	cookie := createRefreshCookie(kakaoToken.RefreshToken, kakaoToken.RefreshTokenExpiresIn)
	setCookie(resp.Headers, cookie)

	resp.Body = string(data)
	resp.StatusCode = http.StatusOK

	return resp, nil
}
