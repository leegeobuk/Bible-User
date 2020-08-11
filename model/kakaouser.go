package model

import "time"

// KakaoUser is a user model for those who signed up through Kakao
type KakaoUser struct {
	id       int
	nickname string
	signedAt	time.Time
}
