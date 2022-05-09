package test

import (
	"errors"
	"time"

	"github.com/Czcan/TimeLine/utils/jsonwt"
	"github.com/golang-jwt/jwt"
)

type JWTClientMock struct{}

func (client *JWTClientMock) GetToken(Uid string) (string, error) {
	return "123123", nil
}

func (client *JWTClientMock) ParseToken(token string) (*jsonwt.Token, error) {
	if token != "123123" {
		return nil, errors.New("invalid token")
	}
	claim := &jsonwt.Token{
		"123123",
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 2).Unix(),
		},
	}
	return claim, nil
}
