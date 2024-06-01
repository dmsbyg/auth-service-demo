package auth

import (
	"context"

	"github.com/dmsbyg/auth-service-demo/internal/auth/token"
	"github.com/dmsbyg/auth-service-demo/internal/common"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type service struct {
	repository Repository
	tokenMaker token.Maker
}

func NewService(repository Repository, tokenMaker token.Maker) Service {
	return &service{
		repository: repository,
		tokenMaker: tokenMaker,
	}
}

func (s *service) Register(ctx context.Context, email string, password []byte) (token string, err error) {
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return "", common.InternalAppError{Err: err}
	}
	uuid := uuid.New()
	err = s.repository.CreateUser(ctx, uuid.String(), email, hashedPassword)
	if err != nil {
		return "", err
	}

	token, err = s.tokenMaker.Make(uuid.String(), email)
	if err != nil {
		return "", err
	}

	return token, err
}

func (s *service) Login(ctx context.Context, email string, password []byte) (token string, err error) {
	panic("not implemented") // TODO: Implement
}
