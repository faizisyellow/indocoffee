package service_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/faizisyellow/indocoffee/internal/db"
	"github.com/faizisyellow/indocoffee/internal/repository/invitations"
	"github.com/faizisyellow/indocoffee/internal/repository/users"
	"github.com/faizisyellow/indocoffee/internal/service"
	"github.com/faizisyellow/indocoffee/internal/utils"
)

func TestUserService(t *testing.T) {
	t.Run("run in memory store", func(t *testing.T) {
		UserServiceTest{
			CreateDependencies: func() (users.Users, invitations.Invitations, utils.Token, db.Transactioner, Cleanup) {
				return &users.InMemoryUsers{}, &invitations.InMemoryInvitations{},
					tokenInvitationFake{token: "lizzy is the goddess of saddness"}, &transactionFake{state: initial}, func() {
						// nothing to clean up
					}
			},
		}.Test(t)
	})
}

type Cleanup func()

type UserServiceTest struct {
	CreateDependencies func() (users.Users, invitations.Invitations, utils.Token, db.Transactioner, Cleanup)
}

func (u UserServiceTest) Test(t *testing.T) {
	t.Run("register new user", func(t *testing.T) {
		var (
			ctx                             = context.Background()
			usr, invt, tkn, tranx, teardown = u.CreateDependencies()
			sut                             = service.UsersServices{usr, invt, tkn, tranx}
			request                         = service.RegisterRequest{
				Username: "lizzy",
				Email:    "lizzymcalpine@test.test",
				Password: "Something123$$",
			}
		)
		t.Cleanup(teardown)

		_, err := sut.RegisterAccount(ctx, request)
		if err != nil {
			t.Error("should not be error")
		}
	})

	t.Run("activate new user", func(t *testing.T) {
		var (
			ctx                             = context.Background()
			usr, invt, tkn, tranx, teardown = u.CreateDependencies()
			sut                             = service.UsersServices{usr, invt, tkn, tranx}
			request                         = service.ActivatedRequest{
				Token: "lizzy is the goddess of sadness",
			}
		)
		t.Cleanup(teardown)

		_, err := sut.RegisterAccount(ctx, service.RegisterRequest{
			Username: "elizabeth",
			Email:    "elizabeth@test.test",
			Password: "Lizzy2442$",
		})
		if err != nil {
			t.Error("should not be error")
		}

		err = sut.ActivateAccount(ctx, request.Token)
		if err != nil {
			t.Error("should not be error")
		}
	})

	t.Run("login user", func(t *testing.T) {
		var (
			ctx                             = context.Background()
			usr, invt, tkn, tranx, teardown = u.CreateDependencies()
			sut                             = service.UsersServices{usr, invt, tkn, tranx}
			request                         = service.LoginRequest{
				Email:    "elizabeth@test.test",
				Password: "Lizzy2442$",
			}
		)
		t.Cleanup(teardown)

		tok, err := sut.RegisterAccount(ctx, service.RegisterRequest{
			Username: "elizabeth",
			Email:    "elizabeth@test.test",
			Password: "Lizzy2442$",
		})
		if err != nil {
			t.Error("should not be error")
		}

		err = sut.ActivateAccount(ctx, tok.Token)
		if err != nil {
			t.Error("should not be error")
		}

		_, err = sut.Login(ctx, request)
		if err != nil {
			t.Error("should not be error")
		}
	})

	t.Run("get user's profile that already activate", func(t *testing.T) {
		var (
			ctx                             = context.Background()
			usr, invt, tkn, tranx, teardown = u.CreateDependencies()
			sut                             = service.UsersServices{usr, invt, tkn, tranx}
		)
		t.Cleanup(teardown)

		_, err := sut.FindUserById(ctx, 2)
		if err != nil {
			t.Error("should not be error")
		}

	})

	t.Run("delete user", func(t *testing.T) {
		var (
			ctx                             = context.Background()
			usr, invt, tkn, tranx, teardown = u.CreateDependencies()
			sut                             = service.UsersServices{usr, invt, tkn, tranx}
		)
		t.Cleanup(teardown)
		err := sut.DeleteAccount(ctx, 2)
		if err != nil {
			t.Error("should not be error")
		}
	})

}

type tokenInvitationFake struct {
	token string
}

func (t tokenInvitationFake) Generate() string {
	return t.token
}

type stateTransaction int

const (
	initial stateTransaction = iota
	begin
	open
	failed
)

func (s *stateTransaction) string(stat stateTransaction) string {
	return []string{"init", "begin", "open", "failed"}[stat]
}

type transactionFake struct {
	state stateTransaction
}

func (t *transactionFake) Begin() (*sql.Tx, error) {
	t.state = begin
	if t.state.string(begin) != "begin" {
		return nil, errors.New("transaction not begin")
	}
	return nil, nil
}

func (t *transactionFake) Commit() error {
	t.state = open
	if t.state.string(open) != "open" {
		return errors.New("transaction errors")
	}
	return nil
}

func (t *transactionFake) Rollback() error {
	t.state = failed
	if t.state.string(failed) != "failed" {
		return errors.New("transaction rollback errors")
	}
	return nil
}

func (t *transactionFake) WithTx(ctx context.Context, fn func(tx *sql.Tx) error) error {
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
