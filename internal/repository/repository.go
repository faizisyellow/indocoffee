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

	Roles interface {
		Insert(ctx context.Context, nw RolesModel) error
		GetAll(ctx context.Context) ([]RolesModel, error)
		GetById(ctx context.Context, id int) (RolesModel, error)
		Update(ctx context.Context, nw RolesModel) error
		Delete(ctx context.Context, id int) error
		Destroy(ctx context.Context, id int) error
	}

	Beans interface {
		Insert(ctx context.Context, nw BeansModel) error
		GetAll(ctx context.Context) ([]BeansModel, error)
		GetById(ctx context.Context, id int) (BeansModel, error)
		Update(ctx context.Context, nw BeansModel) error
		Delete(ctx context.Context, id int) error
		Destroy(ctx context.Context, id int) error
	}

	Forms interface {
		Insert(ctx context.Context, nw FormsModel) error
		GetAll(ctx context.Context) ([]FormsModel, error)
		GetById(ctx context.Context, id int) (FormsModel, error)
		Update(ctx context.Context, nw FormsModel) error
		Delete(ctx context.Context, id int) error
		Destroy(ctx context.Context, id int) error
	}
}

const QueryTimeout = time.Second * 5

func New(db *sql.DB) *Repository {

	return &Repository{
		Users:      &UsersRepository{Db: db},
		Invitation: &InvitationRepository{Db: db},
		Roles:      &RolesRepository{Db: db},
		Beans:      &BeansRepository{Db: db},
		Forms:      &FormsRepository{Db: db},
	}
}
