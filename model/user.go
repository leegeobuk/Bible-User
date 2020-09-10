package model

import "github.com/jinzhu/gorm"

// User is a user model
type User struct {
	gorm.Model `gorm:"embedded"`
	UserID     string `gorm:"unique,not null"`
	Password   string
	Name       string `gorm:"not null"`
	Type       string `gorm:"not null"`
}
