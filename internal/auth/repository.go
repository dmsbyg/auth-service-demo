package auth

import (
	"context"
	"database/sql"
)

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) CreateUser(ctx context.Context, email string, hashedPassword []byte) (token string, err error) {
	panic("not implemented") // TODO: Implement
}

func (r *repository) GetUserByEmail(ctx context.Context, email string) (user User, err error) {
	panic("not implemented") // TODO: Implement
}
