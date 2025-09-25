package users

import (
	"context"
	"database/sql"

	"github.com/faizisyellow/indocoffee/internal/models"
	"github.com/faizisyellow/indocoffee/internal/repository"
)

// This is Users repository to access Users
// From database.
type UsersRepository struct {
	Db *sql.DB
}

// Insert inserts new usr to database.
// Returns usr's id and nil on success, or -1 and an error on failure.
func (u *UsersRepository) Insert(ctx context.Context, tx *sql.Tx, usr models.User) (int, error) {

	query := `
	INSERT INTO users(username,email,password,role_id)
	VALUES(?,?,?,?)
	`

	ctx, cancel := context.WithTimeout(ctx, repository.QueryTimeout)
	defer cancel()

	res, err := tx.ExecContext(ctx, query, usr.Username, usr.Email, usr.Password.HashedText, usr.RoleId)
	if err != nil {
		return -1, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return -1, err
	}

	return int(id), nil
}

// GetById gets a User by User's id from database.
// Returns a User and nil on success or empty User and an error on failure.
func (u *UsersRepository) GetById(ctx context.Context, id int) (models.User, error) {

	var user models.User

	query := `SELECT id,username,email,is_active,password,created_at FROM users WHERE id = ?`

	ctx, cancel := context.WithTimeout(ctx, repository.QueryTimeout)
	defer cancel()

	err := u.Db.QueryRowContext(ctx, query, id).Scan(
		&user.Id,
		&user.Username,
		&user.Email,
		&user.IsActive,
		&user.Password.HashedText,
		&user.CreatedAt,
	)
	if err != nil {
		return user, err
	}

	return user, nil
}

// GetById gets a User by User's email from database.
// Returns a User and nil on success or empty User and an error on failure.
func (u *UsersRepository) GetByEmail(ctx context.Context, email string) (models.User, error) {

	var user models.User

	query := `SELECT id,username,email,password,is_active FROM users WHERE email = ?`

	ctx, cancel := context.WithTimeout(ctx, repository.QueryTimeout)
	defer cancel()

	err := u.Db.QueryRowContext(ctx, query, email).Scan(
		&user.Id,
		&user.Username,
		&user.Email,
		&user.Password.HashedText,
		&user.IsActive,
	)
	if err != nil {
		return user, err
	}

	return user, nil

}

// Update updates a User and
// Ensure usr has id.
// Returns nil on success or an error on failure.
func (u *UsersRepository) Update(ctx context.Context, tx *sql.Tx, usr models.User) error {

	query := `UPDATE users SET username = ?,  email = ?, password = ?, is_active = ? WHERE id = ?`

	ctx, cancel := context.WithTimeout(ctx, repository.QueryTimeout)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, &usr.Username, &usr.Email, &usr.Password.HashedText, &usr.IsActive, usr.Id)
	if err != nil {
		return err
	}

	return nil
}

// Delete delete a User by User's id
// Returns nil on success or an error on failure.
func (u *UsersRepository) Delete(ctx context.Context, tx *sql.Tx, id int) error {

	query := `DELETE FROM users WHERE id = ?`

	ctx, cancel := context.WithTimeout(ctx, repository.QueryTimeout)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}
