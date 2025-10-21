package auth_test

import (
	"testing"
	"time"

	"github.com/faizisyellow/indocoffee/internal/auth"
	"github.com/golang-jwt/jwt/v5"
)

func TestGenerateToken(t *testing.T) {
	t.Run("return a token", func(t *testing.T) {
		jwtService := auth.New("hello-world", "it just test", "test")
		result, err := jwtService.GenerateToken(jwt.MapClaims{
			"iss": "it just test",
			"sub": "test",
			"iat": time.Now().Unix(),
			"exp": time.Now().Add(time.Second * 2).Unix(),
			"id":  1,
		})
		if err != nil {
			t.Errorf("exptected no error :%v", err.Error())
			return
		}

		if result == "" {
			t.Error("exptected return a token but got empty string")
			return
		}
	})

	t.Run("get token multiple times", func(t *testing.T) {
		jwtService := auth.New("hello-world", "it just test", "test")

		resultOne, err := jwtService.GenerateToken(jwt.MapClaims{
			"iss": "it just test",
			"sub": "test",
			"iat": time.Now().Unix(),
			"exp": time.Now().Add(time.Second * 2).Unix(),
			"id":  1,
		})
		if err != nil {
			t.Errorf("expected no error: %v", err.Error())
			return
		}

		if resultOne == "" {
			t.Error("expected return a token but got empty string")
			return
		}

		// Wait to ensure different timestamp
		time.Sleep(time.Second * 1)

		resultTwo, err := jwtService.GenerateToken(jwt.MapClaims{
			"iss": "it just test",
			"sub": "test",
			"iat": time.Now().Unix(),
			"exp": time.Now().Add(time.Second * 2).Unix(),
			"id":  1,
		})
		if err != nil {
			t.Errorf("expected no error: %v", err.Error())
			return
		}

		if resultTwo == "" {
			t.Error("expected return a token but got empty string")
			return
		}

		if resultOne == resultTwo {
			t.Error("expected different token but got the same")
			return
		}
	})
}
