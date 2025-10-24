package loginLimiter

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisLoginLimiter struct {
	Rdb      *redis.Client
	Limit    int
	Duration time.Duration
}

func (r *RedisLoginLimiter) CanLogin(ctx context.Context, ip, email string) (bool, error) {
	var item ItemLoginLimiter
	err := r.Rdb.HGetAll(ctx, GetRedisKey(ip, email)).Scan(&item)
	if err == redis.Nil {
		return true, nil
	} else if err != nil {
		return false, err
	}

	return item.AttemptCount < r.Limit, nil
}

func (r *RedisLoginLimiter) SetAttemptLogin(ctx context.Context, ip, email string, expireHours int) error {
	key := GetRedisKey(ip, email)

	exist, err := r.Rdb.Exists(ctx, key).Result()
	if err != nil {
		return err
	}

	if exist == 0 {
		fields := map[string]any{
			"email":        email,
			"attemptCount": 1,
		}
		if err := r.Rdb.HSet(ctx, key, fields).Err(); err != nil {
			return err
		}
		return r.Rdb.Expire(ctx, key, time.Duration(expireHours)*r.Duration).Err()
	}

	if err := r.Rdb.HIncrBy(ctx, key, "attemptCount", 1).Err(); err != nil {
		return err
	}

	return nil
}

func GetRedisKey(ip, email string) string {
	emailHash := sha1.Sum([]byte(email))
	key := strings.Builder{}
	key.WriteString("login:")
	key.WriteString(ip)
	key.WriteString(":")
	key.WriteString(hex.EncodeToString(emailHash[:]))

	return key.String()
}
