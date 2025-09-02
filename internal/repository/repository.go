package repository

import (
	"context"
	"database/sql"
	"time"
)

type Repository struct {
	Users interface {
		Insert(ctx context.Context, tx *sql.Tx, usr UserModel) (int, error)
		GetById(ctx context.Context, id int) (UserModel, error)
		GetByEmail(ctx context.Context, email string) (UserModel, error)
		Update(ctx context.Context, tx *sql.Tx, usr UserModel) error
		Delete(ctx context.Context, tx *sql.Tx, id int) error
	}

	Invitation interface {
		Insert(ctx context.Context, tx *sql.Tx, invt InvitationModel) error
		Get(ctx context.Context, tx *sql.Tx, token string) (int, error)
		DeleteByUserId(ctx context.Context, tx *sql.Tx, usrid int) error
	}
}

const QueryTimeout = time.Second * 5

func New(db *sql.DB) *Repository {

	return &Repository{
		Users:      &UsersRepository{Db: db},
		Invitation: &InvitationRepository{Db: db},
	}
}
