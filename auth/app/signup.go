package app

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/leegeobuk/Bible-User/auth"
	"github.com/leegeobuk/Bible-User/dbutil"
)

// Signup validates new user information and saves it if valid
func Signup(request *events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	resp := auth.Response(request)
	
	// unmarshal request
	req := &signupRequest{}
	err := json.Unmarshal([]byte(request.Body), req)
	if err != nil {
		resp.Body = err.Error()
		return resp
	}

	// validate email and password
	err = validate(req)
	if err != nil {
		resp.Body = err.Error()
		resp.StatusCode = http.StatusBadRequest
		return resp
	}

	// connect to database
	db, err := dbutil.ConnectDB()
	defer db.Close()
	if err != nil {
		resp.Body = err.Error()
		return resp
	}

	// unathorize user if already a member
	user := req.toUser()
	if err := dbutil.FindMember(db, user); err == nil {
		resp.Body = auth.ErrAccountExist.Error()
		resp.StatusCode = http.StatusUnauthorized
		return resp
	}

	// add account to db
	if err := db.Create(user).Error; err != nil {
		resp.Body = err.Error()
		return resp
	}

	resp.StatusCode = http.StatusOK
	
	return resp
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