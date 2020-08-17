package kakao

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

// Login authenticates kakao user and decide whether to login or not
func Login(ctx context.Context, request *events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	resp := events.APIGatewayProxyResponse{Headers: header}

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

	// ok if account is a member, unauthorized  if not
	if isMember(user, db) {
		resp.StatusCode = http.StatusOK
	} else {
		resp.StatusCode = http.StatusUnauthorized
	}

	return resp, nil
}
