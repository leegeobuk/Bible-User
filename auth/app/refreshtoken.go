package app

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/dgrijalva/jwt-go"
	"github.com/leegeobuk/Bible-User/auth"
)

// RefreshToken returns new access_token and expiring time in seconds
func RefreshToken(request *events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	resp := auth.Response(request)

	// unmarshal request body
	refreshRequest := &refreshTokenRequest{}
	err := json.Unmarshal([]byte(request.Body), refreshRequest)
	if err != nil {
		resp.Body = err.Error()
		resp.StatusCode = http.StatusInternalServerError
		return resp
	}

	// get refresh_token stored in cookie
	// return error if cookie doesn't exist
	cookieString, ok := request.Headers["Cookie"]
	if !ok {
		resp.Body = auth.ErrEmptyCookie.Error()
		resp.StatusCode = http.StatusUnauthorized
		return resp
	}

	// validate refresh token
	refreshTokenStr := auth.ParseCookie(cookieString)
	claims := &claims{}
	rToken, err := jwt.ParseWithClaims(refreshTokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(refreshSignKey), nil
	})
	if err != nil {
		resp.Body = err.Error()
		resp.StatusCode = http.StatusInternalServerError
		return resp
	}

	// error if refresh token not vaild
	if !rToken.Valid {
		resp.StatusCode = http.StatusUnauthorized
		return resp
	}

	// reissue access token
	aTokenStr, err := generateAccessToken(claims.userID, accessDur)
	if err != nil {
		resp.Body = err.Error()
		resp.StatusCode = http.StatusInternalServerError
		return resp
	}

	refreshTokenResp := &refreshTokenResponse{
		AccessToken: aTokenStr,
		ExpiresIn:   strconv.Itoa(int(accessDur.Seconds())),
	}

	data, err := json.Marshal(refreshTokenResp)
	if err != nil {
		resp.Body = err.Error()
		resp.StatusCode = http.StatusInternalServerError
		return resp
	}

	// reissue refresh token and set as cookie if less than 30 days are left before expiration
	if claims.ExpiresAt-time.Now().Local().Unix() <= int64(refreshDur.Seconds()/2) {
		rTokenStr, err := generateRefreshToken(claims.userID, refreshDur)
		if err != nil {
			resp.Body = err.Error()
			resp.StatusCode = http.StatusInternalServerError
			return resp
		}

		cookie := auth.CreateRefreshCookie(rTokenStr, refreshDur)
		auth.SetCookie(resp.MultiValueHeaders, cookie)
	}

	resp.Body = string(data)
	resp.StatusCode = http.StatusOK

	return resp
}
