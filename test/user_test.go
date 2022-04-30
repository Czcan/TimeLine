package test

import (
	"net/url"
	"strings"
	"testing"

	"github.com/fatih/color"
)

type UserAuthTestCase struct {
	Email         string
	Password      string
	ExpectedUser  string
	ExpectedError string
}

func TestUserAuth(t *testing.T) {
	setup()

	testCases := []UserAuthTestCase{
		{Email: "test1@qq.com", Password: "123456", ExpectedUser: `{"code":200,"data":{"token":"123123","email":"test1@qq.com","nick_name":"123123","avatar":"","gender":0,"age":0}}`},
		{Email: "", Password: "123456", ExpectedError: "email or password is empty"},
		{Email: "test2@qq.com", Password: "123123", ExpectedError: "email or password is error"},
	}

	for i, testCase := range testCases {
		body := Post("/api/auth", url.Values{
			"email":    {testCase.Email},
			"password": {testCase.Password},
		})
		if testCase.ExpectedError != "" && !strings.Contains(body, testCase.ExpectedError) {
			t.Errorf(color.RedString("TestUsersAuth #%v: Expected error %v but got %v", i+1, testCase.ExpectedError, body))
		}
		if testCase.ExpectedUser != "" && testCase.ExpectedUser != body {
			t.Errorf(color.RedString("TestUsersAuth #%v: Expected user %v but got %v", i+1, testCase.ExpectedUser, body))
		}
	}
}
