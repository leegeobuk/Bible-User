package kakao

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

var errEmptyCookie = errors.New("error empty cookie from request")
var errEmptyToken = errors.New("error empty access_token from Kakao API")

// RefreshToken returns new access_token using refresh_token
func RefreshToken(ctx context.Context, request *events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	resp := events.APIGatewayProxyResponse{Headers: headers, StatusCode: http.StatusInternalServerError}

	// unmarshal request body
	refreshRequest := &refreshTokenRequest{}
	err := json.Unmarshal([]byte(request.Body), refreshRequest)
	if err != nil {
		return resp, err
	}

	// get refresh_token stored in cookie
	cookieString, ok := request.Headers["Cookie"]
	if !ok {
		resp.StatusCode = http.StatusUnauthorized
		resp.Body = errEmptyCookie.Error()
		return resp, nil
	}

	// get new access_token from Kakao Login API
	refreshRequest.RefreshToken = parseCookie(cookieString)
	refreshedToken, err := getNewToken(refreshRequest)
	if err != nil {
		if err == errEmptyToken {
			resp.StatusCode = http.StatusUnauthorized
			resp.Body = err.Error()
			return resp, nil
		}
		return resp, err
	}

	// marshal kakaoTokenDTO
	refreshTokenResp := &kakaoTokenDTO{
		AccessToken: refreshedToken.AccessToken,
		ExpiresIn:   refreshedToken.ExpiresIn,
	}
	data, err := json.Marshal(refreshTokenResp)
	if err != nil {
		return resp, err
	}

	// set cookie if new refresh_token is returned as well
	if refreshedToken.RefreshToken != "" {
		resp.Headers = copyHeaders(headers)
		cookie := createRefreshCookie(refreshedToken.RefreshToken, refreshedToken.RefreshTokenExpiresIn)
		setCookie(resp.Headers, cookie)
	}

	resp.Body = string(data)
	resp.StatusCode = http.StatusOK

	return resp, nil
}

func getNewToken(refreshRequest *refreshTokenRequest) (*kakaoRefreshTokenAPIDTO, error) {
	// request to Kakao Login API for new access_token
	tokenURL := createRefreshTokenURL(*refreshRequest)
	resp, err := http.Post(tokenURL, "", nil)
	if err != nil {
		return nil, err
	}

	// unmarshal response for new token
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	refreshToken := &kakaoRefreshTokenAPIDTO{}
	err = json.Unmarshal(data, refreshToken)
	if err != nil {
		return nil, err
	}

	// error if Kakako Login API returns wrong value
	if refreshToken.AccessToken == "" {
		return nil, errEmptyToken
	}

	return refreshToken, nil
}

func parseCookie(cookieString string) string {
	i := strings.Index(cookieString, "=")
	return cookieString[i+1:]
}

func createRefreshTokenURL(req refreshTokenRequest) string {
	kakaoKey := os.Getenv("KAKAO_LOGIN_API_KEY")
	return fmt.Sprintf(
		"%s?grant_type=%s&client_id=%s&refresh_token=%s", tokenBaseURL, req.GrantType, kakaoKey, req.RefreshToken,
	)
}
