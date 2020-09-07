package auth

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/leegeobuk/Bible-User/auth/kakao"
	"github.com/leegeobuk/Bible-User/db"
)

// Signup validates user information and saves it if valid
func Signup(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if request.QueryStringParameters["type"] == "kakao" {
		resp := kakao.Signup(&request)
		addHeaders(resp.Headers, headers)
		return resp, nil
	}

	resp := events.APIGatewayProxyResponse{Headers: headers, StatusCode: http.StatusInternalServerError}

	// unmarshal request
	req := &loginRequest{}
	err := json.Unmarshal([]byte(request.Body), req)
	if err != nil {
		resp.Body = err.Error()
		return resp, nil
	}

	// validate email and password
	err = validate(req)
	if err != nil {
		resp.Body = err.Error()
		resp.StatusCode = http.StatusBadRequest
		return resp, nil
	}

	// connect to database
	database, err := db.ConnectDB()
	if err != nil {
		resp.Body = err.Error()
		return resp, nil
	}

	// unathorize user if already a member
	user := req.toUser()
	if db.IsMember(database, user) {
		resp.Body = errAccountExist.Error()
		resp.StatusCode = http.StatusUnauthorized
		return resp, nil
	}

	// add account to db
	if err := database.Create(user).Error; err != nil {
		resp.Body = err.Error()
		return resp, nil
	}

	resp.Headers = headers
	resp.StatusCode = http.StatusOK

	return resp, nil
}

func validate(req *loginRequest) error {
	var errInvalidEmail = errors.New("error invalid email address")
	var errPasswordNotMatch = errors.New("error password and confirm password not match")
	// var errInvalidPassword = errors.New("error invalid password")

	// validate email
	if !strings.Contains(req.Email, "@") {
		return errInvalidEmail
	}

	// check if password and confirm password are same
	if req.PW != req.ConfirmPW {
		return errPasswordNotMatch
	}

	// validate pw

	return nil
}
