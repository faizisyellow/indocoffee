package users

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

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

	query := `SELECT id,username,email,is_active,password,role_id,created_at FROM users WHERE id = ?`

	ctx, cancel := context.WithTimeout(ctx, repository.QueryTimeout)
	defer cancel()

	err := u.Db.QueryRowContext(ctx, query, id).Scan(
		&user.Id,
		&user.Username,
		&user.Email,
		&user.IsActive,
		&user.Password.HashedText,
		&user.RoleId,
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

	user := models.User{}
	user.Role = &models.RolesModel{}

	query := `
		SELECT
			users.id,
			username,
			email,
			password,
			is_active,
			roles.name
		FROM users JOIN roles ON users.role_id = roles.id
	 	WHERE email = ?`

	ctx, cancel := context.WithTimeout(ctx, repository.QueryTimeout)
	defer cancel()

	err := u.Db.QueryRowContext(ctx, query, email).Scan(
		&user.Id,
		&user.Username,
		&user.Email,
		&user.Password.HashedText,
		&user.IsActive,
		&user.Role.Name,
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

func (u *UsersRepository) GetUsersCart(ctx context.Context, id int) (models.User, error) {
	query := `
	SELECT
		users.id,
		users.username,
		cart_items.id,
		cart_items.quantity,
		products.id,
		products.roasted,
		products.price,
		products.image,
		products.quantity AS product_quantity,
		beans.name AS bean,
		forms.name AS form
	FROM users
	LEFT JOIN cart_items ON cart_items.user_id = users.id
	LEFT JOIN products ON products.id = cart_items.product_id
	LEFT JOIN beans ON beans.id = products.bean_id
	LEFT JOIN forms ON forms.id = products.form_id
	WHERE users.id = ? AND cart_items.status='open';
	`

	ctx, cancel := context.WithTimeout(ctx, repository.QueryTimeout)
	defer cancel()

	rows, err := u.Db.QueryContext(ctx, query, id)
	if err != nil {
		return models.User{}, err
	}
	defer rows.Close()

	var user models.User
	usersMap := make(map[int]*models.User)

	for rows.Next() {
		var (
			cart            models.Cart
			cartId          sql.NullInt64
			quantity        sql.NullInt64
			productId       sql.NullInt64
			roasted         sql.NullString
			price           sql.NullFloat64
			image           sql.NullString
			productQuantity sql.NullInt64
			beanName        sql.NullString
			formName        sql.NullString
		)

		if err := rows.Scan(
			&user.Id,
			&user.Username,
			&cartId,
			&quantity,
			&productId,
			&roasted,
			&price,
			&image,
			&productQuantity,
			&beanName,
			&formName,
		); err != nil {
			return models.User{}, err
		}

		if _, ok := usersMap[user.Id]; !ok {
			usersMap[user.Id] = &user
		}

		// Only append cart if it exists (cartId is valid)
		if cartId.Valid {
			cart.Id = int(cartId.Int64)
			cart.Quantity = int(quantity.Int64)
			cart.ProductId = int(productId.Int64)
			cart.Product.Roasted = roasted.String
			cart.Product.Price = price.Float64
			cart.Product.Image = image.String
			cart.Product.Quantity = int(productQuantity.Int64)
			cart.Product.BeansModel.Name = beanName.String
			cart.Product.FormsModel.Name = formName.String

			usersMap[user.Id].Carts = append(usersMap[user.Id].Carts, cart)
		}
	}

	for _, usr := range usersMap {
		user = *usr
	}

	return user, rows.Err()
}

func (o *UsersRepository) GetUsersOrders(ctx context.Context, qry repository.PaginatedOrdersQuery, usrid int) ([]models.Order, error) {
	query := `
		SELECT
			id,
			idempotency_key,
			customer_id,
			customer_name,
			customer_email,
			status,
			total_price,
			phone_number,
			alternative_phone_number,
			street,
			city,
			created_at,
			items,
			cart_ids
		FROM (
			SELECT
				id,
				idempotency_key,
				customer_id,
				customer_name,
				customer_email,
				status,
				total_price,
				phone_number,
				alternative_phone_number,
				street,
				city,
				created_at,
				items,
				cart_ids
			FROM orders
			WHERE status LIKE concat('%',?,'%') AND customer_id = ?
		) AS filtered_orders
		ORDER BY created_at ` + qry.Sort + `
		LIMIT ?
		OFFSET ?
	`

	ctx, cancel := context.WithTimeout(ctx, repository.QueryTimeout)
	defer cancel()

	rowsResult, err := o.Db.QueryContext(ctx, query, qry.Status, usrid, qry.Limit, qry.Offset)
	if err != nil {
		return nil, err
	}
	defer rowsResult.Close()

	var orders []models.Order

	for rowsResult.Next() {
		var (
			order       models.Order
			itemsJSON   sql.NullString
			cartIdsJSON sql.NullString
		)

		if err := rowsResult.Scan(
			&order.Id,
			&order.IdempotencyKey,
			&order.CustomerId,
			&order.CustomerName,
			&order.CustomerEmail,
			&order.Status,
			&order.TotalPrice,
			&order.PhoneNumber,
			&order.AlternativePhoneNumber,
			&order.Street,
			&order.City,
			&order.CreatedAt,
			&itemsJSON,
			&cartIdsJSON,
		); err != nil {
			return nil, err
		}

		// Unmarshal only if not NULL
		if itemsJSON.Valid && itemsJSON.String != "" {
			err = json.Unmarshal([]byte(itemsJSON.String), &order.Items)
			if err != nil {
				return nil, fmt.Errorf("failed to unmarshal items: %w", err)
			}
		}

		if cartIdsJSON.Valid && cartIdsJSON.String != "" {
			err = json.Unmarshal([]byte(cartIdsJSON.String), &order.CartIds)
			if err != nil {
				return nil, fmt.Errorf("failed to unmarshal items: %w", err)
			}
		}

		orders = append(orders, order)

	}

	return orders, rowsResult.Err()
}
