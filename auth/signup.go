package auth

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/leegeobuk/Bible-User/auth/kakao"
)

// Signup validates user information and saves it if valid
func Signup(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if request.QueryStringParameters["type"] == "kakao" {
		return kakao.Signup(ctx, &request)
	}

	resp := events.APIGatewayProxyResponse{Headers: headers}

	return resp, nil
}
