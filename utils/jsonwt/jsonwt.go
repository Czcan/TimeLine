package jsonwt

import (
	"errors"
	"time"

	"github.com/Czcan/TimeLine/config"
	"github.com/golang-jwt/jwt"
)

type JWTValidate interface {
	GetToken(Uid string) (string, error)
	ParseToken(token string) (*Token, error)
}

type Token struct {
	Uid string `json:"uid"`
	jwt.StandardClaims
}

var _ JWTValidate = (*JWTClient)(nil)

type JWTClient struct {
	Issuer     string
	ExpireTime time.Duration
	SecretKey  []byte
}

func New(secret []byte, expire time.Duration, issure string) *JWTClient {
	c := config.MustGetAppConfig()
	if secret == nil {
		secret = []byte(c.SecretKey)
	}
	if expire == 0 {
		expire = time.Hour * 2
	}
	return &JWTClient{
		Issuer:     issure,
		ExpireTime: expire,
		SecretKey:  secret,
	}
}

func (client *JWTClient) GetToken(uid string) (string, error) {
	c := Token{
		uid,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(client.ExpireTime).Unix(),
			Issuer:    client.Issuer,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return token.SignedString(client.SecretKey)
}

func (client *JWTClient) ParseToken(token string) (*Token, error) {
	tokenClaim, err := jwt.ParseWithClaims(token, &Token{}, func(token *jwt.Token) (interface{}, error) {
		return client.SecretKey, nil
	})
	if err != nil {
		return nil, err
	}
	if c, ok := tokenClaim.Claims.(*Token); ok && tokenClaim.Valid {
		return c, nil
	}
	return nil, errors.New("invalid token")
}
