package service

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/faizisyellow/indocoffee/internal/db"
	"github.com/faizisyellow/indocoffee/internal/models"
	"github.com/faizisyellow/indocoffee/internal/repository"
	"github.com/faizisyellow/indocoffee/internal/repository/invitations"
	"github.com/faizisyellow/indocoffee/internal/repository/users"
	"github.com/faizisyellow/indocoffee/internal/service/dto"
	errorService "github.com/faizisyellow/indocoffee/internal/service/error"
	"github.com/faizisyellow/indocoffee/internal/utils"
)

type UsersServices struct {
	UsersStore       users.Users
	InvitationsStore invitations.Invitations
	Token            utils.Token
	Transaction      db.Transactioner
}

type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=16"`
	Email    string `json:"email" validate:"required,email,min=6,max=32"`
	Password string `json:"password" validate:"required,max=18"`
}

type RegisterResponse struct {
	Token string `json:"token"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type ActivatedRequest struct {
	Token string `json:"token"`
}

var (
	ErrTokenInvitationNotFound = errors.New("invitation not found, please register first")
	ErrUserRegisteredNotFound  = errors.New("user not found, please register first")
	ErrUserNotFound            = errors.New("user not found")
	ErrUserNotActivated        = errors.New("user not activated, please activate first")
	ErrUserAlreadyExist        = errors.New("this user already exists")
	ErrUserInternal            = errors.New("server incounter internal error")
)

const CUSTOMER_ROLE = 3

func (us *UsersServices) RegisterAccount(ctx context.Context, req RegisterRequest) (*RegisterResponse, error) {

	var response = new(RegisterResponse)

	err := utils.IsPasswordValid(req.Password)
	if err != nil {
		return nil, errorService.New(err, err)
	}

	return response, us.Transaction.WithTx(ctx, func(tx *sql.Tx) error {

		var newAccount models.User
		newAccount.Email = req.Email
		newAccount.Username = req.Username
		newAccount.RoleId = CUSTOMER_ROLE

		if err = newAccount.Password.ParseFromPassword(req.Password); err != nil {
			return errorService.New(err, err)
		}

		usrId, err := us.UsersStore.Insert(ctx, tx, newAccount)
		if err != nil {
			switch {
			case strings.Contains(err.Error(), CONFLICT_CODE):
				return errorService.New(ErrUserAlreadyExist, err)
			default:
				return errorService.New(ErrUserInternal, err)
			}

		}

		tokenIvt := us.Token.Generate()

		invt := models.InvitationModel{
			UserId:   usrId,
			Token:    tokenIvt,
			ExpireAt: time.Hour * 24,
		}

		err = us.InvitationsStore.Insert(ctx, tx, invt)
		if err != nil {
			return errorService.New(ErrUserInternal, err)
		}

		// register and invite success, send to response
		response.Token = tokenIvt

		return nil
	})

}

func (us *UsersServices) ActivateAccount(ctx context.Context, token string) error {

	return us.Transaction.WithTx(ctx, func(tx *sql.Tx) error {

		usrId, err := us.InvitationsStore.Get(ctx, tx, token)
		if err != nil {
			return errorService.New(ErrTokenInvitationNotFound, err)
		}

		user, err := us.UsersStore.GetById(ctx, usrId)
		if err != nil {
			return errorService.New(ErrUserRegisteredNotFound, err)
		}

		user.IsActive = utils.BoolToPoint(true)

		err = us.UsersStore.Update(ctx, tx, user)
		if err != nil {
			return errorService.New(ErrUserInternal, err)
		}

		err = us.InvitationsStore.DeleteByUserId(ctx, tx, user.Id)
		if err != nil {
			return errorService.New(ErrUserInternal, err)
		}

		return nil
	})

}

func (us *UsersServices) Login(ctx context.Context, req LoginRequest) (*models.User, error) {

	user, err := us.UsersStore.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, errorService.New(ErrUserNotFound, err)
	}

	if !*user.IsActive {
		return nil, errorService.New(ErrUserNotActivated, ErrUserNotActivated)
	}

	err = user.Password.ComparePassword(req.Password)
	if err != nil {
		return nil, errorService.New(err, err)
	}

	return &user, nil
}

func (us *UsersServices) DeleteAccount(ctx context.Context, usrid int) error {

	return us.Transaction.WithTx(ctx, func(tx *sql.Tx) error {

		err := us.UsersStore.Delete(ctx, tx, usrid)
		if err != nil {
			return errorService.New(ErrUserInternal, err)
		}

		return nil
	})

}

func (us *UsersServices) FindUserById(ctx context.Context, usrid int) (*models.User, error) {

	user, err := us.UsersStore.GetById(ctx, usrid)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, errorService.New(ErrUserNotFound, err)
		default:
			return nil, errorService.New(ErrUserInternal, err)
		}

	}

	return &user, nil
}

func (us *UsersServices) FindUsersCart(ctx context.Context, usrId int) (dto.GetUsersCartResponse, error) {

	userWithCart, err := us.UsersStore.GetUsersCart(ctx, usrId)
	if err != nil {
		return dto.GetUsersCartResponse{}, errorService.New(ErrUserInternal, err)
	}

	var (
		carts []dto.CartItemDetail
	)

	for _, crt := range userWithCart.Carts {
		cart := dto.CartItemDetail{
			Id:       crt.Id,
			Quantity: crt.Quantity,
			Product: dto.CartProductDTO{
				Id:      crt.Product.Id,
				Roasted: crt.Product.Roasted,
				Image:   crt.Product.Image,
				Stock:   crt.Product.Quantity,
				Price:   crt.Product.Price,
				Bean:    dto.CartBeanDTO{Name: crt.Product.BeansModel.Name},
				Form:    dto.CartFormDTO{Name: crt.Product.FormsModel.Name},
			},
		}

		carts = append(carts, cart)
	}

	response := dto.GetUsersCartResponse{
		Id:       userWithCart.Id,
		Username: userWithCart.Username,
		Carts:    carts,
	}

	return response, nil
}

func (u *UsersServices) FindUsersOrders(ctx context.Context, r repository.PaginatedOrdersQuery, usrid int) ([]models.Order, error) {
	orders, err := u.UsersStore.GetUsersOrders(ctx, r, usrid)
	if err != nil {
		return nil, errorService.New(ErrOrdersInternal, err)
	}

	return orders, nil
}
