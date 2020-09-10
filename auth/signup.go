package auth

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/leegeobuk/Bible-User/auth/kakao"
	"github.com/leegeobuk/Bible-User/dbutil"
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
	req := &signupRequest{}
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
	db, err := dbutil.ConnectDB()
	if err != nil {
		resp.Body = err.Error()
		return resp, nil
	}

	// unathorize user if already a member
	user := req.toUser()
	if dbutil.IsMember(db, user) {
		resp.Body = errAccountExist.Error()
		resp.StatusCode = http.StatusUnauthorized
		return resp, nil
	}

	// add account to db
	if err := db.Create(user).Error; err != nil {
		resp.Body = err.Error()
		return resp, nil
	}

	resp.StatusCode = http.StatusOK

	return resp, nil
}

func validate(req *signupRequest) error {
	var errInvalidEmail = errors.New("error invalid email address")
	// var errInvalidPassword = errors.New("error invalid password")

	// validate email
	if !strings.Contains(req.Email, "@") {
		return errInvalidEmail
	}

	// validate pw

	return nil
}
