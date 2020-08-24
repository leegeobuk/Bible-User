package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/leegeobuk/Bible-User/auth"
)

func main() {
	lambda.Start(auth.RefreshToken)
}
