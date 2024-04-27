package authmodel

type AuthRequest struct {
	Login    string `validate:"required,min=4"`
	Password string `validate:"required,min=4"`
}

type AuthResponse struct {
	AccessToken  string
	RefreshToken string
}
