package auth

import (
	"context"
	"database/sql"
	"errors"

	"github.com/dmsbyg/auth-service-demo/internal/common"
	"github.com/dmsbyg/auth-service-demo/pkg/logger"
	"github.com/dmsbyg/auth-service-demo/utils"
	"github.com/jmoiron/sqlx"
	"github.com/mattn/go-sqlite3"
)

type repository struct {
	db     *sqlx.DB
	logger logger.Logger
}

func NewRepository(db *sqlx.DB, l logger.Logger) Repository {
	return &repository{
		db:     db,
		logger: l,
	}
}

func (r *repository) CreateUser(ctx context.Context, id, email string, hashedPassword []byte) (err error) {
	stmt, err := r.db.PrepareContext(ctx, "INSERT INTO users(id, email, password) VALUES(?, ?, ?)")
	if err != nil {
		r.logger.Errorf("prepare context error: %w", err)

		return err
	}

	_, err = stmt.ExecContext(ctx, id, email, hashedPassword)
	if err != nil {
		r.logger.Errorf("exec context error: %w", err)
		return wrapError(err)
	}

	return nil
}

func (r *repository) GetUserByEmail(ctx context.Context, email string) (user User, err error) {
	stmt, err := r.db.PreparexContext(ctx, "SELECT * FROM users WHERE email = ?")
	if err != nil {
		r.logger.Errorf("prepare context error: %w", err)
		return User{}, wrapError(err)
	}

	rows := stmt.QueryRowxContext(ctx, email)
	if err := rows.StructScan(&user); err != nil {
		r.logger.Errorf("scan struct error: %w", err)
		return User{}, wrapError(err)
	}

	return user, nil
}

func wrapError(err error) error {
	var sqlite3Err sqlite3.Error
	if errors.As(err, &sqlite3Err) {
		if sqlite3Err.ExtendedCode == sqlite3.ErrConstraintUnique {
			return common.DuplicateError{Entity: utils.GetDuplicateColumnName(sqlite3Err)}
		}
	}

	if errors.Is(err, sql.ErrNoRows) {
		return common.ErrResourceNotFound
	}

	return err
}
