package users_test

import (
	"testing"

	"github.com/faizisyellow/indocoffee/internal/repository/users"
)

func TestInMemoryUsers(t *testing.T) {
	users.Contract{
		NewUsers: &users.InMemoryUsers{},
	}.Test(t)
}
