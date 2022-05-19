package test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/fatih/color"
)

type AccountListTestCase struct {
	ExpectedError    string
	ExpectedResponse string
}

func TestAccountList(t *testing.T) {
	setup()
	RunSQL(DB, `
		INSERT INTO accounts (id, title, content, images, likers, follwers) VALUES (1, "Account_1", "Account_1", "1,2,3", 5, 6);
		INSERT INTO accounts (id, title, content, images, likers, follwers) VALUES (2, "Account_2", "Account_2", "1,2", 10, 6);
		INSERT INTO accounts (id, title, content, images, likers, follwers) VALUES (3, "Account_3", "Account_3", "1", 11, 6);
	`)
	testCases := []AccountListTestCase{
		{ExpectedResponse: `{"code":200,"data":[{"id":3,"user_id":0,"title":"Account_3","content":"Account_3","likers":11,"follwers":6,"created_at":0,"images":["/accountimg/3/1.jpg"]},{"id":2,"user_id":0,"title":"Account_2","content":"Account_2","likers":10,"follwers":6,"created_at":0,"images":["/accountimg/2/1.jpg","/accountimg/2/2.jpg"]},{"id":1,"user_id":0,"title":"Account_1","content":"Account_1","likers":5,"follwers":6,"created_at":0,"images":["/accountimg/1/1.jpg","/accountimg/1/2.jpg","/accountimg/1/3.jpg"]}],"message":null}`},
	}
	for i, testCase := range testCases {
		body := Get("/api/account/home")
		if testCase.ExpectedResponse != "" && body != testCase.ExpectedResponse {
			t.Errorf("TestAccountList #%v: expected accounts %v but got %v", i+1, testCase.ExpectedResponse, body)
		}
		color.Green("TestAccountList #%v: success", i+1)
	}
}

type AccountDetailTestCase struct {
	ID              string
	ExpectedError   string
	ExpectedReponse string
}

func TestAccountDetail(t *testing.T) {
	setup()
	RunSQL(DB, `
		INSERT INTO accounts (id, title, content, images, likers, follwers) VALUES (1, "Account_1", "Account_1", "1,2,3", 5, 6);
		INSERT INTO comments (id, account_id, user_id, content) VALUES (1, 1, 1, "Comment_1");
		INSERT INTO comments (id, account_id, user_id, content) VALUES (2, 1, 2, "Comment_2");
	`)
	testCases := []AccountDetailTestCase{
		{ID: "0", ExpectedError: `invalid param`},
		{ID: "1", ExpectedReponse: `{"code":200,"data":{"account":{"id":1,"title":"Account_1","content":"Account_1","likers":5,"follwers":5,"created_at":0,"images":["/accountimg/1/1.jpg","/accountimg/1/2.jpg","/accountimg/1/3.jpg"]},"comments":[{"nick_name":"name","content":"Comment_1","avatar_url":"/images/1.jpg","date":0},{"nick_name":"","content":"Comment_2","avatar_url":"/images/2.jpg","date":0}]},"message":null}`},
	}
	for i, testCase := range testCases {
		body := Get(fmt.Sprintf("/api/account/detail/%s", testCase.ID))
		if testCase.ExpectedError != "" && !strings.Contains(body, testCase.ExpectedError) {
			t.Errorf(color.RedString("TestAccountDetail #%v: expected error %v but got %v", i+1, testCase.ExpectedError, body))
		}
		if testCase.ExpectedReponse != "" && body != testCase.ExpectedReponse {
			t.Errorf(color.RedString("TestAccountDetail #%v: expected account_detail %v but got %v", i+1, testCase.ExpectedReponse, body))
		}
		color.Green("TestAccountDetail #%v: success", i+1)
	}
}

type CreateAccountTestCase struct {
}

func TestCreateAccount(t *testing.T) {

}
