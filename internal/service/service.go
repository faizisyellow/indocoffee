package service

import (
	"context"

	"github.com/faizisyellow/indocoffee/internal/db"
	"github.com/faizisyellow/indocoffee/internal/models"
	"github.com/faizisyellow/indocoffee/internal/repository/beans"
	"github.com/faizisyellow/indocoffee/internal/repository/forms"
	"github.com/faizisyellow/indocoffee/internal/repository/invitations"
	"github.com/faizisyellow/indocoffee/internal/repository/roles"
	"github.com/faizisyellow/indocoffee/internal/repository/users"
	"github.com/faizisyellow/indocoffee/internal/service/dto"
	"github.com/faizisyellow/indocoffee/internal/utils"
	"github.com/google/uuid"
)

type Service struct {
	UsersService interface {
		RegisterAccount(ctx context.Context, req RegisterRequest) (*RegisterResponse, error)
		ActivateAccount(ctx context.Context, token string) error
		Login(ctx context.Context, req LoginRequest) (*models.User, error)
		DeleteAccount(ctx context.Context, id int) error
		FindUserById(ctx context.Context, id int) (*models.User, error)
	}

	RolesService interface {
		Create(ctx context.Context, req dto.CreateRoleRequest) (string, error)
		FindAll(ctx context.Context) ([]models.RolesModel, error)
		FindById(ctx context.Context, id int) (models.RolesModel, error)
		Update(ctx context.Context, id int, req dto.UpdateRoleRequest) error
		Delete(ctx context.Context, id int) error
		Remove(ctx context.Context) error
	}

	BeansService interface {
		Create(ctx context.Context, req dto.CreateBeanRequest) (string, error)
		FindAll(ctx context.Context) ([]models.BeansModel, error)
		FindById(ctx context.Context, id int) (models.BeansModel, error)
		Update(ctx context.Context, id int, req dto.UpdateBeanRequest) error
		Delete(ctx context.Context, id int) error
		Remove(ctx context.Context) error
	}

	FormsService interface {
		Create(ctx context.Context, req dto.CreateFormRequest) (string, error)
		FindAll(ctx context.Context) ([]models.FormsModel, error)
		FindById(ctx context.Context, id int) (models.FormsModel, error)
		Update(ctx context.Context, id int, req dto.UpdateFormRequest) error
		Delete(ctx context.Context, id int) error
		Remove(ctx context.Context) error
	}
}

var CONFLICT_CODE = "Error 1062 (23000)"

func New(
	usersStore users.Users,
	invitationsStore invitations.Invitations,
	beansStore beans.Beans,
	formsStore forms.Forms,
	rolesStore roles.Roles,
	tx db.Transactioner,
) *Service {
	return &Service{
		UsersService: &UsersServices{
			UsersStore:       usersStore,
			InvitationsStore: invitationsStore,
			Token:            utils.UUID{Plaintoken: uuid.New().String()},
			Transaction:      tx,
		},
		BeansService: &BeansServices{BeansStore: beansStore},
		FormsService: &FormsServices{FormsStore: formsStore},
		RolesService: &RolesServices{RolesStore: rolesStore},
	}
}
