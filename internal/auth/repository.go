package auth

import (
	"context"
	"database/sql"
	"errors"

	"github.com/dmsbyg/auth-service-demo/internal/common"
	"github.com/jmoiron/sqlx"
	"github.com/mattn/go-sqlite3"
)

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) CreateUser(ctx context.Context, id, email string, hashedPassword []byte) (err error) {
	stmt, err := r.db.PrepareContext(ctx, "INSERT INTO users(id, email, password) VALUES(?, ?, ?)")
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, id, email, hashedPassword)
	if err != nil {
		return wrapError(err)
	}

	return nil
}

func (r *repository) GetUserByEmail(ctx context.Context, email string) (user User, err error) {
	stmt, err := r.db.PreparexContext(ctx, "SELECT * FROM users WHERE email = ?")
	if err != nil {
		return User{}, wrapError(err)
	}

	rows := stmt.QueryRowxContext(ctx, email)
	if err := rows.StructScan(&user); err != nil {
		return User{}, wrapError(err)
	}

	return user, nil
}

func wrapError(err error) error {
	var sqlite3Err sqlite3.Error
	if errors.As(err, &sqlite3Err) {
		if sqlite3Err.ExtendedCode == sqlite3.ErrConstraintUnique {
			return common.DuplicateError{Entity: "email"} // TODO: how to extract entity name from error
		}
	}

	if errors.Is(err, sql.ErrNoRows) {
		return common.ErrResourceNotFound
	}

	return err
}
