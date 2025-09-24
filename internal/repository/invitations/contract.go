package invitations

import (
	"context"
	"database/sql"
	"testing"

	"github.com/faizisyellow/indocoffee/internal/models"
)

type Invitations interface {
	Insert(ctx context.Context, tx *sql.Tx, invt models.InvitationModel) error
	Get(ctx context.Context, tx *sql.Tx, token string) (int, error)
	DeleteByUserId(ctx context.Context, tx *sql.Tx, usrId int) error
}

type Contract struct {
	NewInvitations Invitations
}

func (u Contract) Test(t *testing.T) {

}
