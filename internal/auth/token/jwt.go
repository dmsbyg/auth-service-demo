package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/dmsbyg/auth-service-demo/internal/common"
	"github.com/golang-jwt/jwt/v5"
)

const (
	secretKeyMinLength = 10
)

var ErrKeyTooShort = fmt.Errorf("token shall be minimum %d characters long", secretKeyMinLength)

type JwtClaims struct {
	Email string `json:"email"`
	ID    string `json:"id"`
	jwt.RegisteredClaims
}

func NewJWTMaker(secretkey string, duration time.Duration) (*jwtMaker, error) {
	if len(secretkey) < secretKeyMinLength {
		return nil, ErrKeyTooShort
	}
	return &jwtMaker{
		key:           secretkey,
		tokenDuration: duration,
	}, nil
}

type jwtMaker struct {
	key           string
	tokenDuration time.Duration
}

// Make creates a new jwt token for a specific userID & userEmail
func (j jwtMaker) Make(userID, userEmail string) (tokenString string, err error) {
	if userID == "" || userEmail == "" {
		return "", errors.New("user ID or user email cannot be empty")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, JwtClaims{
		Email: userEmail,
		ID:    userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "bmbl-auth",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.tokenDuration)),
		},
	})

	return token.SignedString([]byte(j.key))
}

func (j jwtMaker) Verify(tokenString string) (payload JwtClaims, err error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, err
		}
		return []byte(j.key), nil
	}
	var claims JwtClaims
	token, err := jwt.ParseWithClaims(tokenString, &claims, keyFunc)
	if err != nil {
		return JwtClaims{}, err
	}

	if !token.Valid {
		return JwtClaims{}, common.ErrUnauthorized
	}

	return claims, nil
}
