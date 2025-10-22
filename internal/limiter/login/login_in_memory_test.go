package loginLimiter_test

import (
	"testing"

	loginLimiter "github.com/faizisyellow/indocoffee/internal/limiter/login"
)

func TestLoginLimiterInMemory(t *testing.T) {
	loginLimiter.Contract{func() (loginLimiter.LoginLimiter, func()) {
		tc := &loginLimiter.InMemoryLoginLimiter{
			Limit: 3,
			LoginItemStore: map[string]loginLimiter.ItemLoginLimiter{
				"u123": {
					Email:        "alice@example.com",
					AttemptCount: 2,
				},
				"x9ab": {
					Email:        "bob@example.com",
					AttemptCount: 1,
				},
				"z77y": {
					Email:        "charlie@example.com",
					AttemptCount: 3,
				},
			},
		}

		dropItems := func() {
			tc.LoginItemStore = nil
		}

		return tc, dropItems
	}}.Test(t)
}
