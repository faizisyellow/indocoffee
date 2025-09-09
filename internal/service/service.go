package service

import (
	"context"
	"database/sql"

	"github.com/faizisyellow/indocoffee/internal/db"
	"github.com/faizisyellow/indocoffee/internal/repository"
	"github.com/faizisyellow/indocoffee/internal/service/dto"
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
		Create(ctx context.Context, req dto.CreateRoleRequest) (string, error)
		FindAll(ctx context.Context) ([]repository.RolesModel, error)
		FindById(ctx context.Context, id int) (repository.RolesModel, error)
		Update(ctx context.Context, id int, req dto.UpdateRoleRequest) error
		Delete(ctx context.Context, id int) error
		Remove(ctx context.Context) error
	}

	BeansService interface {
		Create(ctx context.Context, req dto.CreateBeanRequest) (string, error)
		FindAll(ctx context.Context) ([]repository.BeansModel, error)
		FindById(ctx context.Context, id int) (repository.BeansModel, error)
		Update(ctx context.Context, id int, req dto.UpdateBeanRequest) error
		Delete(ctx context.Context, id int) error
		Remove(ctx context.Context) error
	}

	FormsService interface {
		Create(ctx context.Context, req dto.CreateFormRequest) (string, error)
		FindAll(ctx context.Context) ([]repository.FormsModel, error)
		FindById(ctx context.Context, id int) (repository.FormsModel, error)
		Update(ctx context.Context, id int, req dto.UpdateFormRequest) error
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
