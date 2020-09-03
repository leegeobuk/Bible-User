package kakao

import (
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"github.com/leegeobuk/Bible-User/model"
)

func TestSignup(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	db, err := connectDB()
	defer db.Close()
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
			if result := mockSignup(test.user, db); result != test.want {
				t.Errorf("want: %t, got: %t", test.want, result)
			}
		})
	}
}

func mockSignup(user *model.User, db *gorm.DB) bool {
	if isMember(user, db) {
		return false
	}
	return true
}
