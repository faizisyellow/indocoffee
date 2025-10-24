package loginLimiter

import (
	"context"
	"reflect"
)

type InMemoryLoginLimiter struct {
	Limit          int
	LoginItemStore map[string]ItemLoginLimiter
}

func (ll *InMemoryLoginLimiter) CanLogin(_ context.Context, key, email string) (bool, error) {
	if reflect.DeepEqual(ll.LoginItemStore[key], ItemLoginLimiter{}) || ll.LoginItemStore[key].Email != email {
		return true, nil
	}

	return ll.LoginItemStore[key].AttemptCount < ll.Limit, nil
}

func (ll *InMemoryLoginLimiter) SetAttemptLogin(_ context.Context, key, email string, _ int) error {
	item, exists := ll.LoginItemStore[key]

	if exists && item.Email == email {
		item.AttemptCount++
		ll.LoginItemStore[key] = item
		return nil
	}

	ll.LoginItemStore[key] = ItemLoginLimiter{
		Email:        email,
		AttemptCount: 1,
	}

	return nil
}
