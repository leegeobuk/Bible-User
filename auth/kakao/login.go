package kakao

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

const (
	tokenURL = "https://kauth.kakao.com/oauth/token"
	userURL  = "https://kapi.kakao.com/v2/user/me"
)

// Login authenticates user and decide whether to login or not
func Login(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	nilResp := events.APIGatewayProxyResponse{}

	// get token from Kakao Login API
	kakaoToken, err := getToken(request)
	if err != nil {
		return nilResp, err
	}

	// request member info from Kakao logic
	kakaoUser, err := getUserInfo(kakaoToken)
	// _, err = getUserInfo(kakaoToken)
	if err != nil {
		return nilResp, err
	}

	// finding account in db logic

	// resp := &kakaoLoginResponse{}
	// body, err := json.Marshal(resp)
	body, err := json.Marshal(kakaoUser)
	if err != nil {
		return nilResp, err
	}

	return events.APIGatewayProxyResponse{
		Headers: map[string]string{
			"Access-Control-Allow-Headers": "Content-Type",
			"Access-Control-Allow-Origin":  "*",
		},
		StatusCode: http.StatusOK,
		Body:       string(body),
	}, nil
}

func getToken(request events.APIGatewayProxyRequest) (*kakaoTokenResponse, error) {
	// unmarshal request
	loginRequest := &kakaoLoginRequest{}
	err := json.Unmarshal([]byte(request.Body), loginRequest)
	if err != nil {
		return nil, err
	}

	// post request to Kakao Login API
	tokenURL := createTokenURL(*loginRequest)
	resp, err := http.Post(tokenURL, "", nil)
	if err != nil {
		return nil, err
	}

	// unmarshal response from Kakao Login API
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	token := &kakaoTokenResponse{}
	err = json.Unmarshal(data, token)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func createTokenURL(req kakaoLoginRequest) string {
	url := tokenURL
	url += "?grant_type=" + req.GrantType + "&client_id=" + kakaoKey
	url += "&redirect_uri=" + req.RedirectURI + "&code=" + req.Code
	return url
}

func getUserInfo(token *kakaoTokenResponse) (*kakaoUserResponse, error) {
	// create request and set header
	req, err := http.NewRequest("GET", userURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)

	// send request
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	// unmarshal response
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	kakaoUser := &kakaoUserResponse{}
	err = json.Unmarshal(data, kakaoUser)
	if err != nil {
		return nil, err
	}

	return kakaoUser, nil
}
