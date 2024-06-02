package auth

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/dmsbyg/auth-service-demo/internal/common"
	"github.com/dmsbyg/auth-service-demo/pkg/logger"
	"github.com/dmsbyg/auth-service-demo/utils"
	"github.com/jmoiron/sqlx"
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
		r.logger.Errorf("prepare context error: %s", err)

		return err
	}

	_, err = stmt.ExecContext(ctx, id, email, hashedPassword)
	if err != nil {
		r.logger.Errorf("exec context error: %s", err)
		return wrapError(err)
	}

	return nil
}

func (r *repository) GetUserByEmail(ctx context.Context, email string) (user User, err error) {
	stmt, err := r.db.PreparexContext(ctx, "SELECT * FROM users WHERE email = ?")
	if err != nil {
		r.logger.Errorf("prepare context error: %s", err)
		return User{}, wrapError(err)
	}

	rows := stmt.QueryRowxContext(ctx, email)
	if err := rows.StructScan(&user); err != nil {
		r.logger.Errorf("scan struct error: %s", err)
		return User{}, wrapError(err)
	}

	return user, nil
}

type ErrMessage struct {
	Code         int `json:"Code"`
	ExtendedCode int `json:"ExtendedCode"`
	SystemErrno  int `json:"SystemErrno"`
}

var SQLiteErrCodes = map[string]int{
	"uniqueConstraint": 2067,
}

func wrapError(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return common.ErrResourceNotFound
	}

	err = captureSQLiteErr(err)

	return err
}

func captureSQLiteErr(err error) error {
	parsedErr, marshalErr := json.Marshal(err)
	if marshalErr != nil {
		return err
	}

	var errMsg ErrMessage
	unmarshalErr := json.Unmarshal(parsedErr, &errMsg)
	if unmarshalErr != nil {
		return err
	}

	if errMsg.ExtendedCode == SQLiteErrCodes["uniqueConstraint"] {
		return common.DuplicateError{Entity: utils.GetDuplicateColumnName(err.Error())}
	}
	return err
}
