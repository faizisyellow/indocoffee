package service

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/faizisyellow/indocoffee/internal/db"
	"github.com/faizisyellow/indocoffee/internal/repository"
	errorService "github.com/faizisyellow/indocoffee/internal/service/error"
	"github.com/faizisyellow/indocoffee/internal/utils"
)

type UsersServices struct {
	Repository  repository.Repository
	Token       utils.Token
	Transaction db.Transactioner
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
)

func (us *UsersServices) RegisterAccount(ctx context.Context, req RegisterRequest) (*RegisterResponse, error) {

	var response = new(RegisterResponse)

	err := utils.IsPasswordValid(req.Password)
	if err != nil {
		//Todo: handle error client
		return nil, errorService.New(err, err)
	}

	err = us.Transaction.WithTx(ctx, func(tx *sql.Tx) error {

		var newAccount repository.UserModel
		newAccount.Email = req.Email
		newAccount.Username = req.Username

		if err = newAccount.Password.ParseFromPassword(req.Password); err != nil {
			//Todo: handle error client
			return errorService.New(err, err)
		}

		usrId, err := us.Repository.Users.Insert(ctx, tx, newAccount)
		if err != nil {
			switch {
			case strings.Contains(err.Error(), CONFLICT_CODE):
				return errorService.New(ErrUserAlreadyExist, err)
			default:
				//Todo: handle error client
				return errorService.New(err, err)
			}

		}

		tokenIvt := us.Token.Generate()

		invt := repository.InvitationModel{
			UserId:   usrId,
			Token:    tokenIvt,
			ExpireAt: time.Hour * 24,
		}

		err = us.Repository.Invitation.Insert(ctx, tx, invt)
		if err != nil {
			//Todo: handle error client
			return errorService.New(err, err)
		}

		// register and invite success, send to response
		response.Token = tokenIvt

		return nil
	})

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (us *UsersServices) ActivateAccount(ctx context.Context, token string) error {

	return us.Transaction.WithTx(ctx, func(tx *sql.Tx) error {

		usrId, err := us.Repository.Invitation.Get(ctx, tx, token)
		if err != nil {
			return errorService.New(ErrTokenInvitationNotFound, err)
		}

		user, err := us.Repository.Users.GetById(ctx, usrId)
		if err != nil {
			return errorService.New(ErrUserRegisteredNotFound, err)
		}

		user.IsActive = utils.BoolToPoint(true)

		err = us.Repository.Users.Update(ctx, tx, user)
		if err != nil {
			//TODO: handle error to client
			return errorService.New(err, err)
		}

		err = us.Repository.Invitation.DeleteByUserId(ctx, tx, user.Id)
		if err != nil {
			//TODO: handle error to client
			return errorService.New(err, err)
		}

		return nil
	})

}

func (us *UsersServices) Login(ctx context.Context, req LoginRequest) (*repository.UserModel, error) {

	user, err := us.Repository.Users.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, errorService.New(ErrUserNotFound, err)
	}

	if !*user.IsActive {
		return nil, errorService.New(ErrUserNotActivated, err)
	}

	err = user.Password.ComparePassword(req.Password)
	if err != nil {
		//TODO: handle error to client
		return nil, errorService.New(err, err)
	}

	return &user, nil
}

func (us *UsersServices) DeleteAccount(ctx context.Context, usrid int) error {

	return us.Transaction.WithTx(ctx, func(tx *sql.Tx) error {

		err := us.Repository.Users.Delete(ctx, tx, usrid)
		if err != nil {
			//TODO: handle error to client
			return errorService.New(err, err)
		}

		return nil
	})

}

func (us *UsersServices) FindUserById(ctx context.Context, usrid int) (*repository.UserModel, error) {

	user, err := us.Repository.Users.GetById(ctx, usrid)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, errorService.New(ErrUserNotFound, err)
		default:
			//TODO: handle error to client
			return nil, errorService.New(err, err)
		}

	}

	return &user, nil
}
