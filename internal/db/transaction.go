package db

import (
	"context"
	"database/sql"
)

type Transactioner interface {
	WithTx(ctx context.Context, fn func(tx *sql.Tx) error) error
}
