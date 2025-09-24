package invitations

import (
	"context"
	"database/sql"

	"github.com/faizisyellow/indocoffee/internal/models"
)

type InMemoryInvitations struct {
	Invitation []models.InvitationModel
}

func (i *InMemoryInvitations) Insert(ctx context.Context, _ *sql.Tx, invt models.InvitationModel) error {

	i.Invitation = append(i.Invitation, invt)
	return nil
}

func (i *InMemoryInvitations) Get(ctx context.Context, _ *sql.Tx, token string) (int, error) {
	inv := models.InvitationModel{}

	for _, invitation := range i.Invitation {
		if invitation.Token == token {
			inv = invitation
		}
	}

	return inv.UserId, nil
}

func (in *InMemoryInvitations) DeleteByUserId(ctx context.Context, _ *sql.Tx, usrId int) error {
	for i, invitation := range in.Invitation {
		if invitation.UserId == usrId {
			in.Invitation = append(in.Invitation[:i], in.Invitation[i+1:]...)
			break
		}
	}

	return nil
}
