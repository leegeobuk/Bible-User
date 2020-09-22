package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/dgrijalva/jwt-go"
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
	defer db.Close()
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

	// compare pw
	err = bcrypt.CompareHashAndPassword([]byte(savedUser.Password), []byte(user.Password))
	if err != nil {
		// unauthorize if pw doesn't match
		if err == bcrypt.ErrMismatchedHashAndPassword {
			resp.StatusCode = http.StatusUnauthorized
		}
		resp.Body = err.Error()
		return resp, nil
	}

	// generate access token
	accessClaims := &claims{user.UserID, jwt.StandardClaims{ExpiresAt: int64(6 * time.Hour.Seconds())}}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(os.Getenv("ACCESS_SIGN_KEY")))
	if err != nil {
		resp.Body = err.Error()
		return resp, nil
	}

	// generate refresh token
	refreshClaims := jwt.StandardClaims{ExpiresAt: int64(60 * 24 * time.Hour.Seconds())}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(os.Getenv("REFRESH_SIGN_KEY")))
	if err != nil {
		resp.Body = err.Error()
		return resp, nil
	}

	// set response and cookie
	res := &loginResponse{AccessToken: accessTokenString, ExpiresIn: strconv.FormatInt(accessClaims.ExpiresAt, 10)}
	data, err := json.Marshal(res)

	cookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshTokenString,
		Expires:  time.Now().Local().Add(time.Duration(refreshClaims.ExpiresAt) * time.Second),
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
		HttpOnly: true,
	}

	resp.Headers["Set-Cookie"] = cookie.String()
	resp.Body = string(data)
	resp.StatusCode = http.StatusOK

	return resp, nil
}
