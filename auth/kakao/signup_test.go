package kakao

import (
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"github.com/leegeobuk/Bible-User/dbutil"
	"github.com/leegeobuk/Bible-User/model"
)

func TestSignup(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	db, err := dbutil.ConnectDB()
	defer db.Close()
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	tests := []struct {
		desc string
		user *model.User
		want bool
	}{
		{"UserID: signupfail", &model.User{UserID: "signupfail"}, false},
		{"UserID: signuppass", &model.User{UserID: "signuppass"}, true},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			if result := mockSignup(db, test.user); result != test.want {
				t.Errorf("want: %t, got: %t", test.want, result)
			}
		})
	}
}

func mockSignup(database *gorm.DB, user *model.User) bool {
	if dbutil.IsMember(database, user) {
		return false
	}
	return true
}
