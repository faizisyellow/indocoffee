package repository

import (
	"context"
	"database/sql"
	"time"
)

type InvitationModel struct {
	UserId   int           `json:"user_id"`
	Token    string        `json:"token"`
	ExpireAt time.Duration `json:"expire_at"`
}

type InvitationRepository struct {
	Db *sql.DB
}

// Insert inserts new invitation to database.
// Returns  nil on success or an error on failure.
func (ir *InvitationRepository) Insert(ctx context.Context, tx *sql.Tx, invt InvitationModel) error {

	query := `INSERT INTO invitations(user_id,token,expire_at)
	VALUES(?,?,?)
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, invt.UserId, invt.Token, time.Now().Add(invt.ExpireAt))
	if err != nil {
		return err
	}

	return nil
}

// Get gets an invitation if the invitation is not expired
// Returns a the user's id on success or an error on failure.
func (ir *InvitationRepository) Get(ctx context.Context, tx *sql.Tx, token string) (int, error) {

	query := `SELECT user_id FROM invitations WHERE token = ? AND expire_at > ?;`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	userId := 0

	row := tx.QueryRowContext(ctx, query, token, time.Now())
	err := row.Scan(&userId)

	if err != nil {
		return 0, err
	}

	return userId, nil
}

// DeleteByUserId Deletes an invitation by user's id.
func (ir *InvitationRepository) DeleteByUserId(ctx context.Context, tx *sql.Tx, usrId int) error {

	query := `DELETE FROM invitations WHERE user_id = ?`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, usrId)
	if err != nil {
		return err
	}

	return nil
}
