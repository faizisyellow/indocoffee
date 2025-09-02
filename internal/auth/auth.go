package auth

import "github.com/golang-jwt/jwt/v5"

type Authenticator interface {
	GenerateToken(jwt.Claims) (string, error)
	VerifyToken(token string) (*jwt.Token, error)
}
