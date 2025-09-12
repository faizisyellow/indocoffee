package service

import (
	"context"
	"database/sql"
	"testing"

	"github.com/faizisyellow/indocoffee/internal/repository"
)

type tokenInvitation struct {
	token string
}

func (t tokenInvitation) Generate() string {

	return t.token
}

func TestRegisterAccount(t *testing.T) {

	userMock := &UserRepositoryMock{}
	token := tokenInvitation{token: "this is test token"}

	userServ := UsersServices{
		Repository: repository.Repository{
			Users: userMock,
		},
		Db:       nil,
		TransFnc: nil,
		Token:    token,
	}

	request := RegisterRequest{
		Username: "test69",
		Email:    "test@gmail.com",
		Password: "HelloWorld$123",
	}

	want := &RegisterResponse{Token: token.token}
	got, err := userServ.RegisterAccount(context.Background(), request)
	if err != nil {
		t.Errorf("got error %q but want none", err)
	}

	if got != want {
		t.Errorf("want to equal %v, but got: %v", want, got)
	}
}

type UserRepositoryMock struct{}

func (u *UserRepositoryMock) Insert(ctx context.Context, tx *sql.Tx, usr repository.UserModel) (int, error) {

	return 0, nil
}

func (u *UserRepositoryMock) GetById(ctx context.Context, id int) (repository.UserModel, error) {

	return repository.UserModel{}, nil
}

func (u *UserRepositoryMock) GetByEmail(ctx context.Context, email string) (repository.UserModel, error) {
	return repository.UserModel{}, nil
}

func (u *UserRepositoryMock) Update(ctx context.Context, tx *sql.Tx, usr repository.UserModel) error {

	return nil
}

func (u *UserRepositoryMock) Delete(ctx context.Context, tx *sql.Tx, id int) error {

	return nil
}
