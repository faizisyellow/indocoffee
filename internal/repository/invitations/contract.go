package invitations

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/faizisyellow/indocoffee/internal/models"
	"github.com/faizisyellow/indocoffee/internal/utils"
)

type Invitations interface {
	Insert(ctx context.Context, tx *sql.Tx, invt models.InvitationModel) error
	Get(ctx context.Context, tx *sql.Tx, token string) (int, error)
	DeleteByUserId(ctx context.Context, tx *sql.Tx, usrId int) error
}

type Contract struct {
	NewInvitations func() (Invitations, *sql.Tx, func())
}

func (u Contract) Test(t *testing.T) {
	t.Run("create new invitation", func(t *testing.T) {
		var (
			ctx                      = context.Background()
			invitations, tx, cleanup = u.NewInvitations()
			newInvitations           = models.InvitationModel{
				UserId:   1,
				Token:    utils.UUID{}.Generate(),
				ExpireAt: time.Hour * 24,
			}
		)
		t.Cleanup(cleanup)

		err := invitations.Insert(ctx, tx, newInvitations)
		if err != nil {
			t.Error("should not be error")
		}
	})

	t.Run("get invitation by token and return user's id", func(t *testing.T) {
		var (
			ctx                      = context.Background()
			invitations, tx, cleanup = u.NewInvitations()
			initial                  = models.InvitationModel{
				UserId:   1,
				Token:    utils.UUID{}.Generate(),
				ExpireAt: time.Hour * 24,
			}
		)
		t.Cleanup(cleanup)

		err := invitations.Insert(ctx, tx, initial)
		if err != nil {
			t.Error("should not be error")
		}

		id, err := invitations.Get(ctx, tx, initial.Token)
		if err != nil {
			t.Error("should not be error")
		}

		if id != initial.UserId {
			t.Error("expected to be matched")
		}
	})

	t.Run("delete invitations", func(t *testing.T) {
		var (
			ctx                      = context.Background()
			invitations, tx, cleanup = u.NewInvitations()
		)
		t.Cleanup(cleanup)
		if err := invitations.DeleteByUserId(ctx, tx, 1); err != nil {
			t.Error("should not be error")
		}
	})
}
