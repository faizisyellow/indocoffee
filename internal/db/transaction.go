package db

import (
	"context"
	"database/sql"
)

type Transactioner interface {
	Begin() (*sql.Tx, error)
	Rollback() error
	Commit() error
	WithTx(ctx context.Context, fn func(tx *sql.Tx) error) error
}
