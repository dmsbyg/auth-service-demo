package auth

import "context"

type Service interface {
	Register(ctx context.Context, email string, password []byte) (token string, err error)
	Login(ctx context.Context, email string, password []byte) (token string, err error)
}

type Repository interface {
	CreateUser(ctx context.Context, email string, hashedPassword []byte) (token string, err error)
	GetUserByEmail(ctx context.Context, email string) (user User, err error)
}
