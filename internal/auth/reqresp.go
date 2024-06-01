package auth

type ErrorResponse struct {
	Error string `json:"error"`
}

type RegisterRequest struct {
	Email    string `json:"email" validate:"required"`
	Password []byte `json:"password" validate:"required,min=8"`
}

type RegisterResponse struct {
	Token string `json:"token"`
}
