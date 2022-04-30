package test

import "github.com/Czcan/TimeLine/utils/jwt"

type JWTClientMock struct{}

func (client *JWTClientMock) GetToken(Uid string) (string, error) {
	return "123123", nil
}

func (client *JWTClientMock) ParseToken(token string) (*jwt.Token, error) {
	return nil, nil
}
