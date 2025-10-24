package repo

import (
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Querier interface {
	Exec(sql string, args ...any) (pgconn.CommandTag, error)
	QueryRow(sql string, args ...any) pgx.Row
	Query(sql string, args ...any) (pgx.Rows, error)
}

const (
	codeErrUniqueViolation     = "23505"
	codeErrForeignKeyViolation = "23503"
)

var (
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")
)
