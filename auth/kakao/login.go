package kakao

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" // imported for gorm dialect
	"github.com/leegeobuk/Bible-User/model"
)

const (
	tokenURL = "https://kauth.kakao.com/oauth/token"
	userURL  = "https://kapi.kakao.com/v2/user/me"
)

var header = map[string]string{
	"Access-Control-Allow-Headers": "Content-Type",
	"Access-Control-Allow-Origin":  "*",
}

// Login authenticates user and decide whether to login or not
func Login(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	resp := events.APIGatewayProxyResponse{Headers: header}

	// get token from Kakao Login API
	kakaoToken, err := getToken(request)
	if err != nil {
		return resp, err
	}

	// request member info from Kakao logic
	kakaoUserResp, err := getUserInfo(kakaoToken)
	if err != nil {
		return resp, err
	}

	// finding account in db logic
	db, err := connectDB()
	defer db.Close()
	if err != nil {
		return resp, err
	}

	user := kakaoUserResp.toKakaoUser()

	// if account not found
	if isValidAccount(user, db) {
		resp.StatusCode = http.StatusOK
	} else {
		resp.StatusCode = http.StatusUnauthorized
	}

	return resp, nil
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
	kakaoKey := os.Getenv("KAKAO_LOGIN_API_KEY")
	return fmt.Sprintf("%s?grant_type=%s&client_id=%s&redirect_uri=%s&code=%s",
		tokenURL, req.GrantType, kakaoKey, req.RedirectURI, req.Code)
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

func connectDB() (*gorm.DB, error) {
	dbUser := os.Getenv("DB_USER")
	dbPW := os.Getenv("DB_PW")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	args := fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPW, dbHost, dbName)
	return gorm.Open("mysql", args)
}

func isValidAccount(user *model.KakaoUser, db *gorm.DB) bool {
	if db.First(user, user.ID).RecordNotFound() {
		return false
	}
	return true
}
