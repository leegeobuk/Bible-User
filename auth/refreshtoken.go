package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/dgrijalva/jwt-go"
	"github.com/leegeobuk/Bible-User/auth/kakao"
)

// RefreshToken returns new access_token and expiring time in seconds
func RefreshToken(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if request.QueryStringParameters["type"] == "kakao" {
		resp := kakao.RefreshToken(&request)
		addHeaders(resp.Headers, corsHeaders)
		return resp, nil
	}

	resp := events.APIGatewayProxyResponse{MultiValueHeaders: map[string][]string{}, Headers: corsHeaders, StatusCode: http.StatusInternalServerError}

	// unmarshal request body
	refreshRequest := &refreshTokenRequest{}
	err := json.Unmarshal([]byte(request.Body), refreshRequest)
	if err != nil {
		resp.Body = err.Error()
		resp.StatusCode = http.StatusInternalServerError
		return resp, nil
	}

	// get refresh_token stored in cookie
	// return error if cookie doesn't exist
	cookieString, ok := request.Headers["Cookie"]
	if !ok {
		resp.Body = errEmptyCookie.Error()
		resp.StatusCode = http.StatusUnauthorized
		return resp, nil
	}

	// validate refresh token
	refreshTokenStr := parseCookie(cookieString)
	claims := &claims{}
	rToken, err := jwt.ParseWithClaims(refreshTokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(refreshSignKey), nil
	})
	if err != nil {
		resp.Body = err.Error()
		resp.StatusCode = http.StatusInternalServerError
		return resp, nil
	}

	// error if refresh token not vaild
	if !rToken.Valid {
		resp.StatusCode = http.StatusUnauthorized
		return resp, nil
	}

	// reissue access token
	aTokenStr, err := generateAccessToken(claims.userID, accessDur)
	if err != nil {
		resp.Body = err.Error()
		resp.StatusCode = http.StatusInternalServerError
		return resp, nil
	}

	refreshTokenResp := &refreshTokenResponse{
		AccessToken: aTokenStr,
		ExpiresIn:   strconv.Itoa(int(accessDur.Seconds())),
	}

	data, err := json.Marshal(refreshTokenResp)
	if err != nil {
		resp.Body = err.Error()
		resp.StatusCode = http.StatusInternalServerError
		return resp, nil
	}

	// reissue refresh token and set as cookie if less than 30 days are left before expiration
	if claims.ExpiresAt-time.Now().Local().Unix() <= int64(refreshDur.Seconds()/2) {
		rTokenStr, err := generateRefreshToken(claims.userID, refreshDur)
		if err != nil {
			resp.Body = err.Error()
			resp.StatusCode = http.StatusInternalServerError
			return resp, nil
		}

		cookie := createRefreshCookie(rTokenStr, refreshDur)
		setCookie(resp.MultiValueHeaders, cookie)
	}

	resp.Body = string(data)
	resp.StatusCode = http.StatusOK

	return resp, nil
}
