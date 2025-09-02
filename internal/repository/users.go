package repository

import (
	"context"
	"database/sql"
	"reflect"
	"time"

	"github.com/faizisyellow/indocoffee/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type UserModel struct {
	Id        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  Hashed    `json:"-"`
	IsActive  *bool     `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

// Merge merges the changed field from usr
// to the original UserModel
func (us *UserModel) Merge(usr UserModel) {

	currentUser := reflect.ValueOf(us).Elem()
	newUser := reflect.ValueOf(&usr).Elem()

	for i := 0; i < newUser.NumField(); i++ {
		valueNewUserField := newUser.Field(i).Interface()
		valueCurrentUserField := currentUser.Field(i)

		if value, ok := valueNewUserField.(int); ok {
			if value != 0 {
				valueCurrentUserField.SetInt(int64(value))
			}
		}

		if value, ok := valueNewUserField.(string); ok {
			if value != "" {
				valueCurrentUserField.SetString(value)
			}
		}

		if value, ok := valueNewUserField.(*bool); ok {
			if value != nil {
				valueCurrentUserField.Elem().SetBool(*value)
			}
		}

		if value, ok := valueNewUserField.(Hashed); ok {
			if value.Text != nil {

				var structUserField reflect.Value
				if valueCurrentUserField.Kind() == reflect.Ptr {
					structUserField = valueCurrentUserField.Elem()
				} else {
					structUserField = valueCurrentUserField
				}

				for i := 0; i < structUserField.NumField(); i++ {
					field := structUserField.Field(i)
					fieldType := structUserField.Type().Field(i)

					switch fieldType.Name {
					case "Text":
						if field.CanSet() {
							field.Set(reflect.ValueOf(value.Text))
						}
					case "HashedText":
						if field.CanSet() {
							// regenerate with new password
							hashed, _ := utils.GeneratePassword(*value.Text)
							field.Set(reflect.ValueOf(hashed))
						}
					}
				}
			}
		}

	}
}

type Hashed struct {
	Text       *string
	HashedText []byte
}

// Parse parses password to bcrypt hash.
// Parsed password will be saved to HashedText
// and, keep the original plaintext password to text.
func (h *Hashed) ParseFromPassword(password string) error {

	hashed, err := utils.GeneratePassword(password)
	if err != nil {
		return err
	}

	h.Text = &password
	h.HashedText = hashed

	return nil
}

// ComparePassword compares plaintext password to its hashed password.
// Returns nil on success, or an error on failure.
func (h *Hashed) ComparePassword(password string) error {

	return bcrypt.CompareHashAndPassword(h.HashedText, []byte(password))
}

// This is Users repository to access Users
// From database.
type UsersRepository struct {
	Db *sql.DB
}

// Insert inserts new usr to database.
// Returns usr's id and nil on success, or -1 and an error on failure.
func (u *UsersRepository) Insert(ctx context.Context, tx *sql.Tx, usr UserModel) (int, error) {

	query := `
	INSERT INTO users(username,email,password)
	VALUES(?,?,?)
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	res, err := tx.ExecContext(ctx, query, usr.Username, usr.Email, usr.Password.HashedText)
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
func (u *UsersRepository) GetById(ctx context.Context, id int) (UserModel, error) {

	var user UserModel

	query := `SELECT id,username,email,is_active,password,created_at FROM users WHERE id = ?`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
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
func (u *UsersRepository) GetByEmail(ctx context.Context, email string) (UserModel, error) {

	var user UserModel

	query := `SELECT id,username,email,password,is_active FROM users WHERE email = ?`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
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
func (u *UsersRepository) Update(ctx context.Context, tx *sql.Tx, usr UserModel) error {

	query := `UPDATE users SET username = ?,  email = ?, password = ?, is_active = ? WHERE id = ?`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
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

	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}
