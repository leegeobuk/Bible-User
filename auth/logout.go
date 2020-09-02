package auth

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/leegeobuk/Bible-User/auth/kakao"
)

// Logout logs out users
func Logout(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if request.QueryStringParameters["type"] == "kakao" {
		return kakao.Logout(ctx, &request)
	}

	resp := events.APIGatewayProxyResponse{Headers: headers}

	return resp, nil
}
