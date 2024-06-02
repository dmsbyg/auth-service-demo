package auth

import (
	"context"
	"errors"

	"github.com/dmsbyg/auth-service-demo/internal/auth/token"
	"github.com/dmsbyg/auth-service-demo/internal/common"
	"github.com/dmsbyg/auth-service-demo/pkg/logger"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type service struct {
	repository Repository
	tokenMaker token.Maker
	logger     logger.Logger
}

func NewService(repository Repository, tokenMaker token.Maker, l logger.Logger) Service {
	return &service{
		repository: repository,
		tokenMaker: tokenMaker,
		logger:     l,
	}
}

func (s *service) Register(ctx context.Context, email string, password []byte) (token string, err error) {
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		s.logger.Errorf("hashing password error: %v", err)
		return "", common.NewInternalAppError(err)
	}
	uuid := uuid.New()
	err = s.repository.CreateUser(ctx, uuid.String(), email, hashedPassword)
	if err != nil {
		return "", err
	}

	token, err = s.tokenMaker.Make(uuid.String(), email)
	if err != nil {
		s.logger.Errorf("make token error: %v", err)
		return "", err
	}

	return token, err
}

func (s *service) Login(ctx context.Context, email string, password []byte) (token string, err error) {
	user, err := s.repository.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, common.ErrResourceNotFound) {
			return "", common.ErrUnauthorized
		}
		return "", common.NewInternalAppError(err)
	}
	err = bcrypt.CompareHashAndPassword(user.Password, password)
	if err != nil {
		s.logger.Errorf("compare hash and password error: %v", err)
		return "", common.ErrUnauthorized
	}

	token, err = s.tokenMaker.Make(user.ID, user.Email)
	if err != nil {
		s.logger.Errorf("make token error: %v", err)
		return "", err
	}

	return token, nil
}
