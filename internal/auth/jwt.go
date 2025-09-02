package auth

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type JwtAuthentication struct {
	SecretKey string
	Iss       string
	Sub       string
}

func New(secretKey, iss, sub string) *JwtAuthentication {
	return &JwtAuthentication{
		SecretKey: secretKey,
		Iss:       iss,
		Sub:       sub,
	}
}

func (ja *JwtAuthentication) GenerateToken(claims jwt.Claims) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// endcode token
	signedToken, err := token.SignedString([]byte(ja.SecretKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (ja *JwtAuthentication) VerifyToken(token string) (*jwt.Token, error) {

	return jwt.Parse(
		token,
		ja.VerifyParsedToken,
		jwt.WithIssuer(ja.Iss),
		jwt.WithSubject(ja.Sub),
		jwt.WithExpirationRequired(),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
	)
}

func (ja *JwtAuthentication) VerifyParsedToken(token *jwt.Token) (any, error) {

	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("unexpect signing method %v", token.Header["alg"])
	}

	return []byte(ja.SecretKey), nil
}
