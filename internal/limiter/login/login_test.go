package loginLimiter_test

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	loginLimiter "github.com/faizisyellow/indocoffee/internal/limiter/login"
	"github.com/redis/go-redis/v9"
)

func TestRedisLoginLimiter(t *testing.T) {
	if getEnvironment(t) != "development" {
		t.Skip("skipping test: only runs in development environment")
	}

	loginLimiter.Contract{
		NewLoginLimiter: func() (loginLimiter.LoginLimiter, func()) {
			rdb, err := setupTestRedis(t)
			if err != nil {
				t.Fatal(err)
			}

			if err := setUpInitialData(rdb); err != nil {
				t.Fatal(err)
			}

			redisLoginLimiter := loginLimiter.RedisLoginLimiter{
				Rdb:      rdb,
				Limit:    3,
				Duration: time.Second,
			}

			cleanup := func() {
				ctx := context.Background()
				defer rdb.Close()

				iter := rdb.Scan(ctx, 0, "login:*", 0).Iterator()
				for iter.Next(ctx) {
					if err := rdb.Del(ctx, iter.Val()).Err(); err != nil {
						t.Errorf("failed to delete redis key %q: %v", iter.Val(), err)
					}
				}

				if err := iter.Err(); err != nil {
					t.Errorf("redis scan error: %v", err)
				}
			}

			return &redisLoginLimiter, cleanup
		},
	}.Test(t)

}

func setupTestRedis(t *testing.T) (*redis.Client, error) {
	t.Helper()

	db, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		return nil, err
	}

	return redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Username: os.Getenv("REDIS_USERNAME"),
		Password: os.Getenv("REDIS_PW"),
		DB:       db,
	}), nil
}

func setUpInitialData(rdb *redis.Client) error {
	ctx := context.Background()

	initialData := map[string]loginLimiter.ItemLoginLimiter{
		loginLimiter.GetRedisKey("u123", "alice@example.com"): {
			Email:        "alice@example.com",
			AttemptCount: 2,
		},
		loginLimiter.GetRedisKey("x9ab", "bob@example.com"): {
			Email:        "bob@example.com",
			AttemptCount: 1,
		},
		loginLimiter.GetRedisKey("z77y", "charlie@example.com"): {
			Email:        "charlie@example.com",
			AttemptCount: 3,
		},
	}

	for key, item := range initialData {
		err := rdb.HSet(ctx, key, map[string]any{
			"email":        item.Email,
			"attemptCount": item.AttemptCount,
		}).Err()
		if err != nil {
			return fmt.Errorf("failed to seed key %s: %w", key, err)
		}
	}

	return nil
}

func getEnvironment(t *testing.T) string {
	t.Helper()

	return os.Getenv("ENV")
}
