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

	RolesService interface {
		Create(ctx context.Context, req CreateRoleRequest) (string, error)
		FindAll(ctx context.Context) ([]ResponseRolesFindAll, error)
		FindById(ctx context.Context, id int) (ResponseRolesById, error)
		Update(ctx context.Context, id int, nw repository.RolesModel) error
		Delete(ctx context.Context, id int) error
		Remove(ctx context.Context) error
	}

	BeansService interface {
		Create(ctx context.Context, req RequestCreateBean) (string, error)
		FindAll(ctx context.Context) ([]ResponseFindAll, error)
		FindById(ctx context.Context, id int) (repository.BeansModel, error)
		Update(ctx context.Context, id int, req RequestUpdateBean) error
		Delete(ctx context.Context, id int) error
		Remove(ctx context.Context) error
	}

	FormsService interface {
		Create(ctx context.Context, req CreateFormRequest) (string, error)
		FindAll(ctx context.Context) ([]ResponseFormsFindAll, error)
		FindById(ctx context.Context, id int) (ResponseFormsById, error)
		Update(ctx context.Context, id int, nw repository.FormsModel) error
		Delete(ctx context.Context, id int) error
		Remove(ctx context.Context) error
	}
}

var CONFLICT_CODE = "Error 1062 (23000)"

func New(store repository.Repository, txfnc db.TransFnc, Db *sql.DB) *Service {
	return &Service{
		UsersService: &UsersServices{Repository: store, TransFnc: txfnc, Db: Db},
		RolesService: &RolesServices{Repository: store},
		BeansService: &BeansServices{Repository: store},
		FormsService: &FormsServices{Repository: store},
	}
}
