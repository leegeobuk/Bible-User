package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/leegeobuk/Bible-User/auth/kakao"
)

func main() {
	lambda.Start(kakao.Login)
}
