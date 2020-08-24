package model

import "github.com/jinzhu/gorm"

// User is a user model for those who signed up through Kakao
type User struct {
	gorm.Model
	Nickname     string
	Type         string
	RefreshToken string
}
