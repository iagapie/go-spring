package auth

type SignInDTO struct {
	Email    string `json:"email,omitempty" validate:"required,email,min=3,max=255"`
	Password string `json:"password,omitempty" validate:"required,min=8,max=64"`
}

type RefreshTokenDTO struct {
	Token string `json:"token,omitempty" validate:"required,uuid4"`
}

type TokensResponse struct {
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}
