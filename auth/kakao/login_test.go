package kakao

import (
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"github.com/leegeobuk/Bible-User/dbutil"
	"github.com/leegeobuk/Bible-User/model"
)

func TestLogin(t *testing.T) {
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
		{"UserID: loginpass", &model.User{UserID: "loginpass"}, true},
		{"UserID: loginfail", &model.User{UserID: "loginfail"}, false},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			if result := mockLogin(db, test.user); result != test.want {
				t.Errorf("want: %t, got: %t", test.want, result)
			}
		})
	}
}

func mockLogin(database *gorm.DB, user *model.User) bool {
	if err := dbutil.FindMember(database, user); err == nil {
		return true
	}
	return false
}
