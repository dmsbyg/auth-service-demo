package auth

import "context"

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{repository: repository}
}

func (s *service) Register(ctx context.Context, email string, password []byte) (token string, err error) {
	panic("not implemented") // TODO: Implement
}

func (s *service) Login(ctx context.Context, email string, password []byte) (token string, err error) {
	panic("not implemented") // TODO: Implement
}
