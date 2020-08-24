package kakao

import (
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"github.com/leegeobuk/Bible-User/model"
)

func TestLogin(t *testing.T) {
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
		{"id: 0", &model.User{Model: gorm.Model{ID: 0}}, false},
		{"id: 1", &model.User{Model: gorm.Model{ID: 1}}, true},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			if result := isMember(test.user, db); result != test.want {
				t.Errorf("want: %t, got: %t ", test.want, result)
			}
		})
	}
}
