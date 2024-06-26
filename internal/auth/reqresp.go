package auth

type ErrorResponse struct {
	Error string `json:"error"`
}

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=8"`
}

type RegisterResponse struct {
	Token string `json:"token"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=8"`
}

type LoginResponse struct {
	Token string `json:"token"`
}
