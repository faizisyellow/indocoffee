package users

import (
	"context"
	"database/sql"
	"reflect"
	"testing"
	"time"

	"github.com/faizisyellow/indocoffee/internal/models"
	"github.com/faizisyellow/indocoffee/internal/utils"
)

type Users interface {
	Insert(ctx context.Context, tx *sql.Tx, usr models.User) (int, error)
	GetById(ctx context.Context, id int) (models.User, error)
	GetByEmail(ctx context.Context, email string) (models.User, error)
	Update(ctx context.Context, tx *sql.Tx, usr models.User) error
	Delete(ctx context.Context, tx *sql.Tx, id int) error
}

type Contract struct {
	NewUsers func() (Users, *sql.Tx, func())
}

func (u Contract) Test(t *testing.T) {
	t.Run("create new user", func(t *testing.T) {
		var (
			ctx                = context.Background()
			users, tx, cleanup = u.NewUsers()
			userPayload        = models.User{
				Username: "lizzy",
				Email:    "lizzy@test.test",
			}
		)
		t.Cleanup(cleanup)

		err := userPayload.Password.ParseFromPassword("Test1234$")
		if err != nil {
			t.Error("should not be error")
		}

		_, err = users.Insert(ctx, tx, userPayload)
		if err != nil {
			t.Error("should not be error")
		}
	})

	t.Run("get user by id", func(t *testing.T) {
		var (
			ctx                = context.Background()
			users, tx, cleanup = u.NewUsers()
			expected           = models.User{
				Id:        1,
				Username:  "lizzy",
				Email:     "lizzy@test.test",
				IsActive:  utils.BoolToPoint(false),
				CreatedAt: time.Time{},
			}
		)
		t.Cleanup(cleanup)

		_, err := users.Insert(ctx, tx, expected)
		if err != nil {
			t.Error("should not be error")
		}

		result, err := users.GetById(ctx, expected.Id)
		if err != nil {
			t.Error("should not be error")
		}

		if !reflect.DeepEqual(expected, result) {
			t.Error("want to be matched")
		}
	})

	t.Run("get user by email", func(t *testing.T) {
		var (
			ctx                = context.Background()
			users, tx, cleanup = u.NewUsers()
			expected           = models.User{
				Id:        1,
				Username:  "lizzy",
				Email:     "lizzy@test.test",
				IsActive:  utils.BoolToPoint(false),
				CreatedAt: time.Time{},
			}
		)
		t.Cleanup(cleanup)

		_, err := users.Insert(ctx, tx, expected)
		if err != nil {
			t.Error("should not be error")
		}

		result, err := users.GetByEmail(ctx, expected.Email)
		if err != nil {
			t.Error("should not be error")
		}

		if !reflect.DeepEqual(expected, result) {
			t.Error("want to be matched")
		}
	})

	t.Run("update user given updated field", func(t *testing.T) {
		var (
			ctx                = context.Background()
			users, tx, cleanup = u.NewUsers()
			initialUser        = models.User{
				Id:        1,
				Username:  "lizzy",
				Email:     "lizzy@test.test",
				IsActive:  utils.BoolToPoint(false),
				CreatedAt: time.Time{},
			}
			payloadUpdateUser = models.User{
				Username: "LizzyMCalpine",
			}
			expectedAfterUpdateUser = models.User{
				Id:        1,
				Username:  "LizzyMCalpine",
				Email:     "lizzy@test.test",
				IsActive:  utils.BoolToPoint(false),
				CreatedAt: time.Time{},
			}
		)

		t.Cleanup(cleanup)

		_, err := users.Insert(ctx, tx, initialUser)
		if err != nil {
			t.Error("should not be error")
		}

		extUsr, err := users.GetById(ctx, initialUser.Id)
		if err != nil {
			t.Error("should not be error")
		}
		// parsing
		extUsr.Username = payloadUpdateUser.Username

		err = users.Update(ctx, nil, extUsr)

		updatedUsr, err := users.GetById(ctx, initialUser.Id)

		if !reflect.DeepEqual(updatedUsr, expectedAfterUpdateUser) {
			t.Error("should not be error")
		}
	})

	t.Run("delete user by id", func(t *testing.T) {
		var (
			ctx                = context.Background()
			users, tx, cleanup = u.NewUsers()
			initialUser        = models.User{
				Id:        1,
				Username:  "lizzy",
				Email:     "lizzy@test.test",
				IsActive:  utils.BoolToPoint(false),
				CreatedAt: time.Time{},
			}
		)

		t.Cleanup(cleanup)

		_, err := users.Insert(ctx, tx, initialUser)
		if err != nil {
			t.Error("should not be error")
		}

		if err := users.Delete(ctx, nil, initialUser.Id); err != nil {
			t.Error("should not be error")
		}
	})
}
