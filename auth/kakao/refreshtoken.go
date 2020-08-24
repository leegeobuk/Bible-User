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

var errEmptyCookie = errors.New("error empty refresh_token cookie")
var errEmptyToken = errors.New("error empty access_token from Kakao API")

// RefreshToken returns new access_token using refresh_token
func RefreshToken(ctx context.Context, request *events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	resp := events.APIGatewayProxyResponse{Headers: headers, StatusCode: http.StatusInternalServerError}

	// get new access_token from Kakao Login API
	refreshedToken, err := getNewToken(request)
	if err != nil {
		if err == errEmptyCookie {
			resp.StatusCode = http.StatusOK
			return resp, err
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

	resp.Body = string(data)
	resp.StatusCode = http.StatusOK

	// set HttpOnly cookie if new refresh_token is returned as well
	if refreshedToken.RefreshToken != "" {
		cookie := createRefreshCookie(refreshedToken.RefreshToken, refreshedToken.RefreshTokenExpiresIn)
		setCookie(resp.Headers, cookie)
	}

	return resp, nil
}

func getNewToken(request *events.APIGatewayProxyRequest) (*kakaoRefreshTokenAPIDTO, error) {
	// unmarshal request body
	refreshRequest := &refreshTokenRequest{}
	err := json.Unmarshal([]byte(request.Body), refreshRequest)
	if err != nil {
		return nil, err
	}

	cookieString, ok := request.Headers["Cookie"]
	if !ok {
		return nil, errEmptyCookie
	}

	// request to Kakao Login API for new access_token
	refreshRequest.RefreshToken = parseCookie(cookieString)
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
