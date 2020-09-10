package kakao

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/leegeobuk/Bible-User/dbutil"
)

// Login authenticates kakao user and decide whether to login or not
func Login(request *events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	resp := events.APIGatewayProxyResponse{Headers: map[string]string{}, StatusCode: http.StatusInternalServerError}

	// get token from Kakao Login API
	kakaoToken, err := getToken(request)
	if err != nil {
		resp.Body = err.Error()
		return resp
	}

	// request member info from Kakao logic
	kakaoUserResp, err := getUserInfo(kakaoToken)
	if err != nil {
		resp.Body = err.Error()
		return resp
	}

	// finding account in db logic
	db, err := dbutil.ConnectDB()
	defer db.Close()
	if err != nil {
		resp.Body = err.Error()
		return resp
	}

	user := kakaoUserResp.toKakaoUser()

	// unauthorize if not a member
	if err := dbutil.FindMember(db, user); err != nil {
		resp.Body = err.Error()
		resp.StatusCode = http.StatusUnauthorized
		return resp
	}

	// marshal kakaoLoginResponse
	loginResp := &kakaoTokenDTO{AccessToken: kakaoToken.AccessToken, ExpiresIn: kakaoToken.ExpiresIn}
	data, err := json.Marshal(loginResp)
	if err != nil {
		resp.Body = err.Error()
		return resp
	}

	// set httpOnly cookie
	cookie := createRefreshCookie(kakaoToken.RefreshToken, kakaoToken.RefreshTokenExpiresIn)
	setCookie(resp.Headers, cookie)

	resp.Body = string(data)
	resp.StatusCode = http.StatusOK

	return resp
}
