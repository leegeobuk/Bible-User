package kakao

import (
	"github.com/jinzhu/gorm"
	"github.com/leegeobuk/Bible-User/model"
)

// kakaoLoginRequest is a request from client for Kakao login
type kakaoLoginRequest struct {
	GrantType   string `json:"grantType"`
	RedirectURI string `json:"redirectUri"`
	Code        string `json:"code"`
}

// kakaoTokenAPIDTO is a token response from Kakao Login API
type kakaoTokenAPIDTO struct {
	TokenType             string `json:"token_type"`
	AccessToken           string `json:"access_token"`
	ExpiresIn             int    `json:"expires_in"`
	RefreshToken          string `json:"refresh_token"`
	RefreshTokenExpiresIn int    `json:"refresh_token_expires_in"`
	Scope                 string `json:"scope"`
}

// kakaoUserAPIDTO is a user info response from Kakao User API
type kakaoUserAPIDTO struct {
	ID           uint         `json:"id"`
	Properties   properties   `json:"properties"`
	KakaoAccount kakaoAccount `json:"kakao_account"`
}

func (r *kakaoUserAPIDTO) toKakaoUser() *model.User {
	return &model.User{
		Model:    gorm.Model{ID: r.ID},
		Type:     "kakao",
		Nickname: r.KakaoAccount.Profile.Nickname,
	}
}

type properties struct {
	Nickname          string `json:"nickname"`
	ProfileImageURL   string `json:"profile_image_url"`
	ThumbnailImageURL string `json:"thumbnail_image_url"`
}

type kakaoAccount struct {
	Profile  profile `json:"profile"`
	Email    string  `json:"email"`
	AgeRange string  `json:"age_range"`
	Birthday string  `json:"birthday"`
	Gender   string  `json:"gender"`
}

type profile struct {
	Nickname              string `json:"nickname"`
	ProfileImage          string `json:"profile_image"`
	ThumbnailImageURL     string `json:"thumbnail_image_url"`
	ProfileNeedsAgreement string `json:"profile_needs_agreement"`
}

type kakaoTokenDTO struct {
	AccessToken string `json:"accessToken"`
	ExpiresIn   int    `json:"expiresIn"`
}

// refreshTokenRequest is a request from client to refresh access_token
type refreshTokenRequest struct {
	GrantType    string `json:"grantType"`
	RefreshToken string `json:"refreshToken,omitempty"`
}

type kakaoRefreshTokenAPIDTO struct {
	TokenType             string `json:"token_type"`
	AccessToken           string `json:"access_token"`
	ExpiresIn             int    `json:"expires_in"`
	RefreshToken          string `json:"refresh_token,omitempty"`
	RefreshTokenExpiresIn int    `json:"refresh_token_expires_in,omitempty"`
}
