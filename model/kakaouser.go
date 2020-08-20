package model

import "github.com/jinzhu/gorm"

// KakaoUser is a user model for those who signed up through Kakao
type KakaoUser struct {
	gorm.Model
	Nickname     string
	RefreshToken string
}
