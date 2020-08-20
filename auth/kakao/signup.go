package kakao

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

// Signup validates kakao user information and saves it if valid
func Signup(ctx context.Context, request *events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// use headers since it is not modified later
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

	// unauthorized if already a member
	if isMember(user, db) {
		resp.StatusCode = http.StatusUnauthorized
		return resp, nil
	}

	// add account to db
	if err := db.Create(user).Error; err != nil {
		return resp, err
	}
	resp.StatusCode = http.StatusOK
	resp.Headers = headers

	return resp, nil
}
