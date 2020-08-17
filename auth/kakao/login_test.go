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
		user *model.KakaoUser
		want bool
	}{
		{"id: 0", &model.KakaoUser{Model: gorm.Model{ID: 0}}, false},
		{"id: 1", &model.KakaoUser{Model: gorm.Model{ID: 1}}, true},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			if result := isMember(tc.user, db); result != tc.want {
				t.Errorf("want: %t, got: %t ", tc.want, result)
			}
		})
	}
}
