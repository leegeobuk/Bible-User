package dbutil

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" // imported for gorm dialect
	"github.com/leegeobuk/Bible-User/model"
)

// ConnectDB connects to database
func ConnectDB() (*gorm.DB, error) {
	dbUser := os.Getenv("DB_USER")
	dbPW := os.Getenv("DB_PW")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	args := fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPW, dbHost, dbName)
	return gorm.Open("mysql", args)
}

// IsMember checks if the user is in the database
func IsMember(db *gorm.DB, user *model.User) bool {
	return !db.Find(user, model.User{UserID: user.UserID}).RecordNotFound()
}
