package kakao

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/leegeobuk/Bible-User/dbutil"
)

// Signup validates kakao user information and saves it if valid
func Signup(request *events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
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

	// unauthorized if already a member
	if dbutil.IsMember(db, user) {
		resp.StatusCode = http.StatusUnauthorized
		resp.Body = errAccountExist.Error()
		return resp
	}

	// add account to db
	if err := db.Create(user).Error; err != nil {
		resp.Body = err.Error()
		return resp
	}

	resp.StatusCode = http.StatusOK

	return resp
}
