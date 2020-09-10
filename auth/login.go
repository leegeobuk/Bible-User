package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/leegeobuk/Bible-User/auth/kakao"
	"github.com/leegeobuk/Bible-User/dbutil"
	"github.com/leegeobuk/Bible-User/model"
	"golang.org/x/crypto/bcrypt"
)

// Login authenticates user and decide whether to login or not
func Login(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if request.QueryStringParameters["type"] == "kakao" {
		resp := kakao.Login(&request)
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

	// connect to db
	db, err := dbutil.ConnectDB()
	if err != nil {
		resp.Body = err.Error()
		return resp, nil
	}

	// unauthorize if id doesn't exist
	user := &model.User{UserID: req.Email, Password: req.PW}
	savedUser := &model.User{UserID: user.UserID}
	if err := dbutil.FindMember(db, savedUser); err != nil {
		resp.Body = err.Error()
		resp.StatusCode = http.StatusUnauthorized
		return resp, nil
	}

	// decrypt pw
	fmt.Println("user pw", user.Password)
	fmt.Println("hashed pw", savedUser.Password)
	err = bcrypt.CompareHashAndPassword([]byte(savedUser.Password), []byte(user.Password))
	if err != nil {
		// unauthorize if pw doesn't match
		if err == bcrypt.ErrMismatchedHashAndPassword {
			resp.StatusCode = http.StatusUnauthorized
		}
		resp.Body = err.Error()
		return resp, nil
	}

	resp.StatusCode = http.StatusOK

	return resp, nil
}
