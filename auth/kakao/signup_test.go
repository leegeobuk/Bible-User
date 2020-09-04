package kakao

import (
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"github.com/leegeobuk/Bible-User/db"
	"github.com/leegeobuk/Bible-User/model"
)

func TestSignup(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	database, err := db.ConnectDB()
	defer database.Close()
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	tests := []struct {
		desc string
		user *model.User
		want bool
	}{
		{"id: 1", &model.User{Model: gorm.Model{ID: 1}}, false},
		{"id: 2", &model.User{Model: gorm.Model{ID: 2}}, true},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			if result := mockSignup(database, test.user); result != test.want {
				t.Errorf("want: %t, got: %t", test.want, result)
			}
		})
	}
}

func mockSignup(database *gorm.DB, user *model.User) bool {
	if db.IsMember(database, user) {
		return false
	}
	return true
}
