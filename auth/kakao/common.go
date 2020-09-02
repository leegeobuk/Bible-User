package kakao

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" // imported for gorm dialect
	"github.com/leegeobuk/Bible-User/model"
)

const (
	origin       = "http://localhost:3000"
	tokenBaseURL = "https://kauth.kakao.com/oauth/token"
	userURL      = "https://kapi.kakao.com/v2/user/me"
)

var (
	headers = map[string]string{
		"Access-Control-Allow-Credentials": "true",
		"Access-Control-Allow-Headers":     "Content-Type",
		"Access-Control-Allow-Origin":      origin,
	}
	errEmptyCookie = errors.New("error empty cookie from request")
	errEmptyToken  = errors.New("error empty access_token from Kakao API")
)

func copyHeaders(headers map[string]string) map[string]string {
	m := make(map[string]string)
	for k, v := range headers {
		m[k] = v
	}

	return m
}

func getToken(request *events.APIGatewayProxyRequest) (*kakaoTokenAPIDTO, error) {
	// unmarshal request body
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

	token := &kakaoTokenAPIDTO{}
	err = json.Unmarshal(data, token)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func createTokenURL(req kakaoLoginRequest) string {
	kakaoKey := os.Getenv("KAKAO_LOGIN_API_KEY")
	return fmt.Sprintf(
		"%s?grant_type=%s&client_id=%s&redirect_uri=%s&code=%s",
		tokenBaseURL, req.GrantType, kakaoKey, req.RedirectURI, req.Code,
	)
}

func getUserInfo(token *kakaoTokenAPIDTO) (*kakaoUserAPIDTO, error) {
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

	kakaoUser := &kakaoUserAPIDTO{}
	err = json.Unmarshal(data, kakaoUser)
	if err != nil {
		return nil, err
	}

	return kakaoUser, nil
}

func createRefreshCookie(value string, seconds int) *http.Cookie {
	return &http.Cookie{
		Name:     "refresh_token",
		Value:    value,
		Expires:  time.Now().Local().Add(time.Duration(seconds) * time.Second),
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
		HttpOnly: true,
	}
}

func setCookie(h map[string]string, c *http.Cookie) {
	cookieString := c.String()
	h["Set-Cookie"] = cookieString
}

func parseCookie(cookieString string) string {
	i := strings.Index(cookieString, "=")
	return cookieString[i+1:]
}

func connectDB() (*gorm.DB, error) {
	dbUser := os.Getenv("DB_USER")
	dbPW := os.Getenv("DB_PW")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	args := fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPW, dbHost, dbName)
	return gorm.Open("mysql", args)
}

func isMember(user *model.User, db *gorm.DB) bool {
	return !db.First(user, user.ID).RecordNotFound()
}
