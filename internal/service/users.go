package service

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/faizisyellow/indocoffee/internal/db"
	"github.com/faizisyellow/indocoffee/internal/repository"
	"github.com/faizisyellow/indocoffee/internal/utils"
)

type UsersServices struct {
	Repository repository.Repository
	Db         *sql.DB
	TransFnc   db.TransFnc
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
		return nil, err
	}

	err = us.TransFnc(us.Db, ctx, func(tx *sql.Tx) error {

		var newAccount repository.UserModel
		newAccount.Email = req.Email
		newAccount.Username = req.Username

		if err = newAccount.Password.ParseFromPassword(req.Password); err != nil {
			return err
		}

		usrId, err := us.Repository.Users.Insert(ctx, tx, newAccount)
		if err != nil {
			duplicateKey := "Error 1062"
			switch {
			case strings.Contains(err.Error(), duplicateKey):
				return ErrUserAlreadyExist
			default:
				return err
			}

		}

		tokenIvt := utils.GenerateTokenUuid()

		invt := repository.InvitationModel{
			UserId:   usrId,
			Token:    tokenIvt,
			ExpireAt: time.Hour * 24,
		}

		err = us.Repository.Invitation.Insert(ctx, tx, invt)
		if err != nil {
			return err
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

	err := us.TransFnc(us.Db, ctx, func(tx *sql.Tx) error {

		usrId, err := us.Repository.Invitation.Get(ctx, tx, token)
		if err != nil {
			return ErrTokenInvitationNotFound
		}

		user, err := us.Repository.Users.GetById(ctx, usrId)
		if err != nil {
			return ErrUserRegisteredNotFound
		}

		user.IsActive = utils.BoolToPoint(true)

		err = us.Repository.Users.Update(ctx, tx, user)
		if err != nil {
			return err
		}

		err = us.Repository.Invitation.DeleteByUserId(ctx, tx, user.Id)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil

}

func (us *UsersServices) Login(ctx context.Context, req LoginRequest) (*repository.UserModel, error) {

	user, err := us.Repository.Users.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, ErrUserNotFound
	}

	if !*user.IsActive {
		return nil, ErrUserNotActivated
	}

	err = user.Password.ComparePassword(req.Password)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (us *UsersServices) DeleteAccount(ctx context.Context, usrid int) error {

	return us.TransFnc(us.Db, ctx, func(tx *sql.Tx) error {

		err := us.Repository.Users.Delete(ctx, tx, usrid)
		if err != nil {
			return err
		}

		return nil
	})

}

func (us *UsersServices) FindUserById(ctx context.Context, usrid int) (*repository.UserModel, error) {

	user, err := us.Repository.Users.GetById(ctx, usrid)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrUserNotFound
		default:
			return nil, err
		}

	}

	return &user, nil
}
