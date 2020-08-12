package kakao

import "os"

var (
	kakaoKey = os.Getenv("KAKAO_LOGIN_API_KEY")
	dbUser   = os.Getenv("DB_USER")
	dbPW     = os.Getenv("DB_PW")
	dbHost   = os.Getenv("DB_HOST")
	dbName   = os.Getenv("DB_NAME")
)
