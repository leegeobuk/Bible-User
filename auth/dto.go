package auth

import (
	"strings"

	"github.com/leegeobuk/Bible-User/model"
)

type loginRequest struct {
	Email     string `json:"email"`
	PW        string `json:"password"`
	ConfirmPW string `json:"confirmPassword"`
}

func (r loginRequest) toUser() *model.User {
	return &model.User{
		Nickname: r.Email[:strings.Index(r.Email, "@")],
		Type:     "app",
	}
}
