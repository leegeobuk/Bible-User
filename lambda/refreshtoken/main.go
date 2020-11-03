package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/leegeobuk/Bible-User/auth/app"
	"github.com/leegeobuk/Bible-User/auth/kakao"
)

func main() {
	lambda.Start(refreshToken)
}

func refreshToken(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if request.QueryStringParameters["type"] == "kakao" {
		resp := kakao.RefreshToken(&request)
		return resp, nil
	}

	resp := app.RefreshToken(&request)
	return resp, nil
}
