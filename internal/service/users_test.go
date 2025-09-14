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
	"golang.org/x/crypto/bcrypt"
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
	tx, err := t.Begin()
	if err != nil {
		return err
	}

	err = fn(tx)

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

func TestActivationAccount(t *testing.T) {
	transMock := transactionMock{status: initial}
	tknRequest := tokenInvitation{
		token: "this is token test",
	}
	usersSrv := setupUsersServiceTest(tokenInvitation{}, &transMock)

	t.Run("activate an account by given a token invitation", func(t *testing.T) {
		err := usersSrv.ActivateAccount(context.Background(), tknRequest.token)
		assertNoError(t, err)
	})

	t.Run("activate an account fails because the token invitation not found", func(t *testing.T) {
		err := usersSrv.ActivateAccount(context.Background(), "")
		assertError(t, err, ErrTokenInvitationNotFound.Error())
	})

}

func TestLoginAccount(t *testing.T) {

	usersSrv := setupUsersServiceTest(tokenInvitation{}, nil)

	t.Run("success to login and return the user", func(t *testing.T) {
		req := LoginRequest{
			Email:    "lizzymcalpine@test.com",
			Password: "HelloWorld$123",
		}
		want := &repository.UserModel{
			Email: "lizzymcalpine@test.com",
		}
		resl, err := usersSrv.Login(context.Background(), req)
		assertNoError(t, err)
		if want.Email != resl.Email {
			t.Errorf("should success login but return after login not match, got :%v", resl)
		}
	})

	t.Run("login fails because the user not found", func(t *testing.T) {
		req := LoginRequest{
			Email:    "batman@test.com",
			Password: "HelloWorld$123",
		}
		_, err := usersSrv.Login(context.Background(), req)
		assertError(t, err, ErrUserNotFound.Error())
	})

	t.Run("login fails because user's password not matched", func(t *testing.T) {
		req := LoginRequest{
			Email:    "lizzymcalpine@test.com",
			Password: "a singer",
		}

		_, err := usersSrv.Login(context.Background(), req)
		assertError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())
	})

	t.Run("login fails because the user's account not activated yet", func(t *testing.T) {
		req := LoginRequest{
			Email:    "coolemailname@test.com",
			Password: "HelloWorld$123",
		}
		_, err := usersSrv.Login(context.Background(), req)
		assertError(t, err, ErrUserNotActivated.Error())
	})
}

func TestDeleteAccount(t *testing.T) {
	transMock := transactionMock{status: initial}
	usersSrv := setupUsersServiceTest(tokenInvitation{}, &transMock)

	err := usersSrv.DeleteAccount(context.Background(), 21)
	assertNoError(t, err)

}

func TestFindUser(t *testing.T) {
	usersSrv := setupUsersServiceTest(tokenInvitation{}, nil)

	t.Run("find the user and return it", func(t *testing.T) {
		want := &repository.UserModel{
			Id:       21,
			Username: "lizzy",
			Email:    "lizzymcalpine@test.com",
		}
		resl, err := usersSrv.FindUserById(context.Background(), want.Id)
		assertNoError(t, err)
		if !reflect.DeepEqual(resl, want) {
			t.Errorf("should success get the user but return not match, got :%v", resl)
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

func assertError(t testing.TB, err error, got string) {
	t.Helper()

	if err == nil {
		t.Errorf("expected an error but got none ")
		return
	}

	if err.Error() != got {
		t.Errorf("expected an error with message: %q but got error message :%q", got, err)
	}
}

func assertNoError(t testing.TB, err error) {
	if err != nil {
		t.Errorf("expected an no error but got an error: %v ", err)
	}
}
