package auth

import (
	"github.com/leegeobuk/Bible-User/model"
)

type signupRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	PW    string `json:"pw"`
}

func (r signupRequest) toUser() *model.User {
	return &model.User{
		UserID:   r.Email,
		Password: r.PW,
		Name:     r.Name,
		Type:     "app",
	}
}

type loginRequest struct {
	Email string `json:"email"`
	PW    string `json:"pw"`
}
