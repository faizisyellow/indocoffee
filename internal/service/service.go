package service

import (
	"context"
	"database/sql"

	"github.com/faizisyellow/indocoffee/internal/db"
	"github.com/faizisyellow/indocoffee/internal/repository"
)

type Service struct {
	UsersService interface {
		RegisterAccount(ctx context.Context, req RegisterRequest) (*RegisterResponse, error)
		ActivateAccount(ctx context.Context, token string) error
		Login(ctx context.Context, req LoginRequest) (*repository.UserModel, error)
		DeleteAccount(ctx context.Context, id int) error
		FindUserById(ctx context.Context, id int) (*repository.UserModel, error)
	}
}

func New(store repository.Repository, txfnc db.TransFnc, Db *sql.DB) *Service {
	return &Service{
		UsersService: &UsersServices{Repository: store, TransFnc: txfnc, Db: Db},
	}
}
