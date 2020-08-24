package auth

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/leegeobuk/Bible-User/auth/kakao"
)

var headers = map[string]string{
	"Access-Control-Allow-Headers": "Content-Type",
	"Access-Control-Allow-Origin":  "http://localhost:3000",
}

// Login authenticates user and decide whether to login or not
func Login(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if request.QueryStringParameters["type"] == "kakao" {
		return kakao.Login(ctx, &request)
	}

	resp := events.APIGatewayProxyResponse{Headers: headers}

	return resp, nil
}
