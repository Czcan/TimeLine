package test

import (
	"net/url"
	"strings"
	"testing"
	"time"

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
		{Email: "test1@qq.com", Password: "123456", ExpectedUser: `{"code":200,"data":{"token":"123123","email":"test1@qq.com","nick_name":"name","avatar":"","gender":0,"age":0},"message":null}`},
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
		color.Green("TestUsersAuth #%v: Success", i+1)
	}
}

type UserRegisterTestCase struct {
	Email            string
	Password         string
	PasswordAgain    string
	ExpectedResponse string
	ExpectedError    string
}

func TestUserRegister(t *testing.T) {
	setup()
	testCases := []UserRegisterTestCase{
		{Email: "test1@qq.com", Password: "123456", PasswordAgain: "123456", ExpectedError: "用户已存在"},
		{Email: "test3@qq.com", Password: "123456", PasswordAgain: "123123", ExpectedError: "incorrect password"},
		{Email: "test5@qq.com", Password: "123456", PasswordAgain: "123456", ExpectedResponse: "注册成功"},
		{Email: "", Password: "123456", ExpectedError: "email is empty"},
	}
	time.Sleep(time.Second * 1)
	for i, testCase := range testCases {
		body := Post("/api/register", url.Values{
			"email":     {testCase.Email},
			"password":  {testCase.Password},
			"password1": {testCase.PasswordAgain},
		})
		if testCase.ExpectedError != "" && !strings.Contains(body, testCase.ExpectedError) {
			t.Errorf(color.RedString("TestUserRegister #%v: Expected error %v but got %v", i+1, testCase.ExpectedError, body))
		}
		if testCase.ExpectedResponse != "" && !strings.Contains(body, testCase.ExpectedResponse) {
			t.Errorf(color.RedString("TestUserRegister #%v: Expected resp %v but got %v", i+1, testCase.ExpectedResponse, body))
		}
		color.Green("TestUserRegister #%v: Success", i+1)
	}
}
