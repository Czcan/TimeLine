package test

import (
	"errors"
	"time"

	"github.com/Czcan/TimeLine/utils/jwt"
	jwtPkg "github.com/golang-jwt/jwt"
)

type JWTClientMock struct{}

func (client *JWTClientMock) GetToken(Uid string) (string, error) {
	return "123123", nil
}

func (client *JWTClientMock) ParseToken(token string) (*jwt.Token, error) {
	if token != "123123" {
		return nil, errors.New("invalid token")
	}
	claim := &jwt.Token{
		"123123",
		jwtPkg.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 2).Unix(),
		},
	}
	return claim, nil
}
