package mock

import (
	"context"
	"database/sql"
	"errors"

	"github.com/faizisyellow/indocoffee/internal/repository"
)

type UsersRepositoryMock struct {
}

func (u *UsersRepositoryMock) Insert(ctx context.Context, tx *sql.Tx, usr repository.UserModel) (int, error) {
	existingUser := repository.UserModel{
		Email: "lizzymcalpine@test.com",
	}

	if usr.Email == existingUser.Email {
		return 0, errors.New("account already exist")
	}

	return usr.Id, nil
}

func (u *UsersRepositoryMock) GetById(ctx context.Context, id int) (repository.UserModel, error) {

	return repository.UserModel{}, nil
}

func (u *UsersRepositoryMock) GetByEmail(ctx context.Context, email string) (repository.UserModel, error) {
	return repository.UserModel{}, nil
}

func (u *UsersRepositoryMock) Update(ctx context.Context, tx *sql.Tx, usr repository.UserModel) error {

	return nil
}

func (u *UsersRepositoryMock) Delete(ctx context.Context, tx *sql.Tx, id int) error {

	return nil
}

type InvitationRepositoryMock struct {
}

func (i *InvitationRepositoryMock) Insert(ctx context.Context, tx *sql.Tx, invt repository.InvitationModel) error {

	return nil
}
func (i *InvitationRepositoryMock) Get(ctx context.Context, tx *sql.Tx, token string) (int, error) {

	return 0, nil
}
func (i *InvitationRepositoryMock) DeleteByUserId(ctx context.Context, tx *sql.Tx, usrid int) error {

	return nil
}
