package kakao

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
)

// RefreshToken returns new access_token using refresh_token
func RefreshToken(request *events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	resp := events.APIGatewayProxyResponse{Headers: map[string]string{}, MultiValueHeaders: map[string][]string{}, StatusCode: http.StatusInternalServerError}

	// unmarshal request body
	refreshRequest := &refreshTokenRequest{}
	err := json.Unmarshal([]byte(request.Body), refreshRequest)
	if err != nil {
		resp.Body = err.Error()
		return resp
	}

	// get refresh_token stored in cookie
	cookieString, ok := request.Headers["Cookie"]
	if !ok {
		resp.Body = errEmptyCookie.Error()
		resp.StatusCode = http.StatusUnauthorized
		return resp
	}

	// get new access_token from Kakao Login API
	refreshRequest.RefreshToken = parseCookie(cookieString)
	refreshedToken, err := getNewToken(refreshRequest)
	if err != nil {
		resp.Body = err.Error()
		if err == errEmptyToken {
			resp.StatusCode = http.StatusUnauthorized
		}
		return resp
	}

	// marshal refreshTokenResponse
	refreshTokenResp := &refreshTokenResponse{
		AccessToken: refreshedToken.AccessToken,
		ExpiresIn:   refreshedToken.ExpiresIn,
	}

	data, err := json.Marshal(refreshTokenResp)
	if err != nil {
		resp.Body = err.Error()
		resp.StatusCode = http.StatusInternalServerError
		return resp
	}

	// set cookie if new refresh_token is returned as well
	if refreshedToken.RefreshToken != "" {
		cookie := createRefreshCookie(refreshedToken.RefreshToken, refreshedToken.RefreshTokenExpiresIn)
		setCookie(resp.MultiValueHeaders, cookie)
	}

	resp.Body = string(data)
	resp.StatusCode = http.StatusOK

	return resp
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

func createRefreshTokenURL(req refreshTokenRequest) string {
	kakaoKey := os.Getenv("KAKAO_LOGIN_API_KEY")
	return fmt.Sprintf(
		"%s?grant_type=%s&client_id=%s&refresh_token=%s", tokenBaseURL, req.GrantType, kakaoKey, req.RefreshToken,
	)
}
