package test

import (
	"strings"
	"testing"

	"github.com/fatih/color"
)

type UserSendMailTestCase struct {
	Email         string
	ExpectedError string
}

type UserRegisterTestCase struct {
	Email         string
	Code          string
	Password      string
	ExpectedError string
}

// func TestUserSendMail(t *testing.T) {
// 	testCases := []UserSendMailTestCase{
// 		{Email: "", ExpectError: "邮箱不能为空"},
// 		{Email: "1048196021@qq.com", ExpectError: ""},
// 	}

// 	for i, testCase := range testCases {
// 		body := httpPost("/users/auth", url.Values{"email": {testCase.Email}})

// 		if testCase.ExpectedError == "" && strings.Contains(body, "error") {
// 			t.Error(color.RedString("TestUserSendMail #%v: Expected don't get error but got %v\n", i+1, body))
// 		}

// 		if testCase.ExpectedError != "" && !strings.Contains(body, testCase.ExpectedError) {
// 			t.Error(color.RedString("TestUserSendMail #%v: Expected get error %v but got %v\n", i+1, testCase.ExpectedError, body))
// 		}

// 		color.Green("TestUserSendMail #%v: Success!", i+1)
// 	}
// }

func TestUserRegister(t *testing.T) {
	testCases := []UserRegisterTestCase{
		{Email: "", Code: "2537", Password: "123456", ExpectedError: "邮箱不能为空"},
		{Email: "1048196021@qq.com", Code: "", Password: "123456", ExpectedError: "验证码不能为空"},
		{Email: "1048196021@qq.com", Code: "2537", Password: "", ExpectedError: "密码不能为空"},
	}

	for i, testCase := range testCases {
		// body := Post("/users/register", url.Values{"email": {testCase.Email}, "auth_code": {testCase.Code}, "password": {testCase.Password}})

		if testCase.ExpectedError == "" && strings.Contains(body, "error") {
			t.Error(color.RedString("TestUserSendMail #%v: Expected don't get error but got %v\n", i+1, body))
		}

		if testCase.ExpectedError != "" && !strings.Contains(body, testCase.ExpectedError) {
			t.Error(color.RedString("TestUserSendMail #%v: Expected get error %v but got %v\n", i+1, testCase.ExpectedError, body))
		}
	}
}
