package kakao

// kakaoLoginRequest is a request from client for Kakao login
type kakaoLoginRequest struct {
	GrantType   string `json:"grant_type"`
	RedirectURI string `json:"redirect_uri"`
	Code        string `json:"code"`
}

// kakaoTokenResponse is a token response from Kakao Login API
type kakaoTokenResponse struct {
	TokenType             string `json:"token_type"`
	AccessToken           string `json:"access_token"`
	ExpiresIn             int    `json:"expires_in"`
	RefreshToken          string `json:"refresh_token"`
	RefreshTokenExpiresIn int    `json:"refresh_token_expires_in"`
	Scope                 string `json:"scope"`
}

// kakaoUserResponse is a user info response from Kakao User API
type kakaoUserResponse struct {
	ID           int          `json:"id"`
	Properties   properties   `json:"properties"`
	KakaoAccount kakaoAccount `json:"kakao_account"`
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

type kakaoLoginResponse struct {
	status int
	msg    string
}
