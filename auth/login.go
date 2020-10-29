package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/leegeobuk/Bible-User/auth/kakao"
	"github.com/leegeobuk/Bible-User/dbutil"
	"github.com/leegeobuk/Bible-User/model"
	"golang.org/x/crypto/bcrypt"
)

const (
	accessTokenValidHours = 6
	refreshTokenValidDays = 60
)

var (
	hourDur    = time.Hour
	dayDur     = 24 * hourDur
	accessDur  = accessTokenValidHours * hourDur
	refreshDur = refreshTokenValidDays * dayDur
)

// Login authenticates user and decide whether to login or not
func Login(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if request.QueryStringParameters["type"] == "kakao" {
		resp := kakao.Login(&request)
		addHeaders(resp.Headers, corsHeaders)
		return resp, nil
	}

	resp := events.APIGatewayProxyResponse{MultiValueHeaders: map[string][]string{}, Headers: corsHeaders, StatusCode: http.StatusInternalServerError}

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
	accessTokenString, err := generateAccessToken(user.UserID, accessDur)
	if err != nil {
		resp.Body = err.Error()
		return resp, nil
	}

	// generate refresh token
	refreshTokenString, err := generateRefreshToken(user.UserID, refreshDur)
	if err != nil {
		resp.Body = err.Error()
		return resp, nil
	}

	// set response and cookie
	res := &loginResponse{
		AccessToken: accessTokenString,
		ExpiresIn:   strconv.Itoa(int(accessDur.Seconds())),
		Type:        "app",
	}
	data, err := json.Marshal(res)

	cookie := createRefreshCookie(refreshTokenString, refreshDur)
	setCookie(resp.MultiValueHeaders, cookie)

	resp.Body = string(data)
	resp.StatusCode = http.StatusOK

	return resp, nil
}
