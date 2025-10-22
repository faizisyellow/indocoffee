package loginLimiter

import (
	"context"
	"testing"
)

type LoginLimiter interface {
	CanLogin(ctx context.Context, key, email string) (bool, error)
	SetAttemptLogin(ctx context.Context, key, email string, expireAt int) error
}

type ItemLoginLimiter struct {
	Email        string `redis:"email"`
	AttemptCount int    `redis:"attemptCount"`
}

type Contract struct {
	NewLoginLimiter func() (LoginLimiter, func())
}

func (c Contract) Test(t *testing.T) {
	t.Run("check attempt success login", func(t *testing.T) {
		testCases := []struct {
			name           string
			key            string
			requestPayload string
		}{
			{
				name:           "should passed to login because email is correct",
				key:            "x02b",
				requestPayload: "lizzy@example.com",
			},
			{
				name:           "should passed to login when the key is exist but request different email",
				key:            "x02b",
				requestPayload: "morelizzy@example.com",
			},
			{
				name:           "should passed to login at second try",
				key:            "x9ab",
				requestPayload: "bob@example.com",
			},
		}

		for _, tc := range testCases {
			var (
				ctx                    = context.Background()
				loginLimiter, teardown = c.NewLoginLimiter()
			)
			t.Cleanup(func() {
				teardown()
			})

			canLogin, err := loginLimiter.CanLogin(ctx, tc.key, tc.requestPayload)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !canLogin {
				t.Fatalf("expected user to be allowed to try login")
			}
		}

	})

	t.Run("check attempt fail login because it's limited", func(t *testing.T) {
		var (
			ctx                    = context.Background()
			loginLimiter, teardown = c.NewLoginLimiter()
		)
		t.Cleanup(func() {
			teardown()
		})

		canLogin, err := loginLimiter.CanLogin(ctx, "z77y", "charlie@example.com")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if canLogin {
			t.Fatalf("expected user to be not allowed to  login")
		}
	})

	t.Run("should block after exceeding max attempts", func(t *testing.T) {
		t.Run("should block after limit is on 3", func(t *testing.T) {
			var (
				ctx                    = context.Background()
				loginLimiter, teardown = c.NewLoginLimiter()
			)
			t.Cleanup(func() {
				teardown()
			})

			err := loginLimiter.SetAttemptLogin(ctx, "u123", "alice@example.com", 5)
			if err != nil {
				t.Fatalf("expected to be not error but got: %v", err.Error())
			}

			canLogin, err := loginLimiter.CanLogin(ctx, "u123", "alice@example.com")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if canLogin {
				t.Fatalf("expected user to be not allowed to  login")
			}
		})

		t.Run("should block after try several times", func(t *testing.T) {
			ctx := context.Background()
			loginLimiter, teardown := c.NewLoginLimiter()
			t.Cleanup(teardown)

			for i := 0; i < 3; i++ {
				err := loginLimiter.SetAttemptLogin(ctx, "x9ab", "bob@example.com", 5)
				if err != nil {
					t.Fatalf("expected no error, got: %v", err)
				}
			}

			canLogin, err := loginLimiter.CanLogin(ctx, "x9ab", "bob@example.com")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if canLogin {
				t.Fatalf("expected user to be blocked after exceeding max attempts")
			}
		})

	})
}
