package service

import (
	"context"
	"database/sql"
	"errors"
	"reflect"
	"testing"

	"github.com/faizisyellow/indocoffee/internal/db"
	"github.com/faizisyellow/indocoffee/internal/mock"
	"github.com/faizisyellow/indocoffee/internal/repository"
)

type tokenInvitation struct {
	token string
}

func (t tokenInvitation) Generate() string {

	return t.token
}

type statusTransaction int

func (s *statusTransaction) string(stat statusTransaction) string {

	return []string{"init", "begin", "open", "failed"}[stat]
}

const (
	initial statusTransaction = iota
	begin
	open
	failed
)

type transactionMock struct {
	status statusTransaction
}

func (t *transactionMock) Begin() (*sql.Tx, error) {
	t.status = begin

	if t.status.string(begin) != "begin" {
		return nil, errors.New("transaction not begin")
	}

	return nil, nil
}

func (t *transactionMock) Rollback() error {
	t.status = failed
	if t.status.string(failed) != "failed" {
		return errors.New("transaction rollback errors")
	}

	return nil
}

func (t *transactionMock) Commit() error {
	t.status = open
	if t.status.string(open) != "open" {
		return errors.New("transaction errors")
	}

	return nil
}

func (t *transactionMock) WithTx(ctx context.Context, fn func(tx *sql.Tx) error) error {
	_, err := t.Begin()
	if err != nil {
		return err
	}

	err = fn(nil)

	if err != nil {
		if err := t.Rollback(); err != nil {
			return err
		}
		return err
	}

	return t.Commit()
}

func TestRegisterAccount(t *testing.T) {
	tkn := tokenInvitation{token: "this is test token"}
	transMock := transactionMock{status: initial}
	usersSrv := setupUsersServiceTest(tkn, &transMock)

	t.Run("register new account", func(t *testing.T) {
		request := RegisterRequest{
			Username: "username",
			Email:    "test@gmail.com",
			Password: "HelloWorld$123",
		}

		want := &RegisterResponse{Token: tkn.token}
		got, err := usersSrv.RegisterAccount(context.Background(), request)
		if err != nil {
			t.Errorf("not expected an error but got one: %q", err)
		}

		if !reflect.DeepEqual(want, got) {
			t.Errorf("want to equal %v, but got: %v", want, got)
		}

	})

	t.Run("account that registered already exist", func(t *testing.T) {
		request := RegisterRequest{
			Username: "lizzy",
			Email:    "lizzymcalpine@test.com",
			Password: "HelloWorld$123",
		}

		got, err := usersSrv.RegisterAccount(context.Background(), request)
		if err == nil {
			t.Error("expected an error but got none")
		}

		if got != nil {
			t.Error("expected response to be empty but has a response")
		}
	})
}

func setupUsersServiceTest(token tokenInvitation, trans db.Transactioner) UsersServices {

	return UsersServices{
		Repository: repository.Repository{
			Users:      &mock.UsersRepositoryMock{},
			Invitation: &mock.InvitationRepositoryMock{},
		},
		Transaction: trans,
		Token:       token,
	}

}
