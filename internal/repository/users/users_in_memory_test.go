package users_test

import (
	"database/sql"
	"testing"

	"github.com/faizisyellow/indocoffee/internal/repository/users"
)

func TestInMemoryUsers(t *testing.T) {
	users.Contract{
		NewUsers: func() (users.Users, *sql.Tx, func()) {
			return &users.InMemoryUsers{}, nil, func() {}
		},
	}.Test(t)
}
