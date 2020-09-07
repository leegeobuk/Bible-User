package auth

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/leegeobuk/Bible-User/auth/kakao"
)

// RefreshToken returns new access_token and expiring time in seconds
func RefreshToken(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if request.QueryStringParameters["type"] == "kakao" {
		resp := kakao.RefreshToken(&request)
		addHeaders(resp.Headers, headers)
		return resp, nil
	}

	resp := events.APIGatewayProxyResponse{Headers: headers}

	return resp, nil
}
