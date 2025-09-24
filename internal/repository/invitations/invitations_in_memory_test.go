package invitations_test

import (
	"database/sql"
	"testing"

	"github.com/faizisyellow/indocoffee/internal/repository/invitations"
)

func TestInMemoryInvitations(t *testing.T) {
	invitations.Contract{
		NewInvitations: func() (invitations.Invitations, *sql.Tx, func()) {
			return &invitations.InMemoryInvitations{}, nil, func() {}
		},
	}.Test(t)
}
