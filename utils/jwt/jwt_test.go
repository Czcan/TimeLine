package jwt

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type JWTTokenTestCase struct {
	Uid           string
	Secret        []byte
	ExpectedUid   string
	ExpectedError string
}

func TestValidate(t *testing.T) {
	c := New([]byte("123123"), time.Hour*2, "jwt")
	testCases := []JWTTokenTestCase{
		{Uid: "123", Secret: []byte("123123"), ExpectedUid: "123"},
		{Uid: "456", Secret: []byte("123"), ExpectedError: "signature is invalid"},
	}
	for _, testCase := range testCases {
		token, err := c.GetToken(testCase.Uid)
		if err == nil {
			a := strings.Split(token, ".")
			assert.Len(t, a, 3)
		} else {
			t.Errorf("get token failed: %v", err.Error())
		}
		c.SecretKey = testCase.Secret
		claims, err := c.ParseToken(token)
		if err == nil {
			assert.Equal(t, testCase.Uid, claims.Uid)
		} else {
			assert.EqualError(t, err, testCase.ExpectedError)
		}
	}
}
