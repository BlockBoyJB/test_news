package service

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"time"
)

const (
	TTL = time.Hour * 24 * 7
)

type authService struct {
	key []byte
}

func newAuthService(key string) *authService {
	return &authService{
		key: []byte(key),
	}
}

func (s *authService) Create() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(TTL).Unix(),
		IssuedAt:  time.Now().Unix(),
	})
	signedToken, err := token.SignedString(s.key)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func (s *authService) Validate(tokenString string) bool {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("incorrect sign method")
		}
		return s.key, nil
	})
	if err != nil {
		return false
	}
	return token.Valid
}
