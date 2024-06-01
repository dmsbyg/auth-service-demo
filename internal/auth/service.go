package auth

import (
	"context"

	"github.com/dmsbyg/auth-service-demo/internal/common"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{repository: repository}
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

	token = "some-generated-token" // TODO: generate token

	return token, err
}

func (s *service) Login(ctx context.Context, email string, password []byte) (token string, err error) {
	panic("not implemented") // TODO: Implement
}
