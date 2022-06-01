package test

import (
	"fmt"
	"net/url"
	"strings"
	"testing"

	"github.com/fatih/color"
)

type LikersTestCase struct {
	Token           string
	ID              string
	Liker           string
	ExpectedError   string
	ExpectedReponse string
	ExpectedAccount string
}

func TestLikers(t *testing.T) {
	setup()
	RunSQL(DB, `
		INSERT INTO accounts (id, title, content, images, likers, follwers) VALUES (1, "Account_1", "Account_1", "1,2,3", 5, 6);
	`)
	testCases := []LikersTestCase{
		{Token: "123456", ID: "0", ExpectedError: `invalid user`},
		{Token: "123123", ID: "1", Liker: "1", ExpectedReponse: `{"code":200,"data":6,"message":null}`, ExpectedAccount: `1,Account_1,Account_1,1,2,3,6,6`},
		{Token: "123123", ID: "1", Liker: "1", ExpectedError: `{"code":400,"data":null,"message":"invalid operation"}`},
	}
	for i, testCase := range testCases {
		body := SingeGet(testCase.Token, fmt.Sprintf("/api/liker?id=%s&liker=%s", testCase.ID, testCase.Liker), nil)
		if testCase.ExpectedError != "" && !strings.Contains(body, testCase.ExpectedError) {
			t.Errorf("TestLikers #%v: expected error %v but got %v", i+1, body, testCase.ExpectedError)
		}
		if testCase.ExpectedReponse != "" && body != testCase.ExpectedReponse {
			t.Errorf("TestLikers #%v: expected response %v but got %v", i+1, testCase.ExpectedReponse, body)
		}
		account := GetRecords(DB, "accounts", "id, title, content, images, likers, follwers", "where id = 1")
		if testCase.ExpectedAccount != "" && account != testCase.ExpectedAccount {
			t.Errorf("TestLikers #%v: expected account %v but got %v", i+1, testCase.ExpectedAccount, account)
		}
		color.Green("TestLikers #%v: success", i+1)
	}
}

type FollwerTestCase struct {
	Token              string
	ID                 string
	Follwer            string
	ExpectedError      string
	ExpectedResponse   string
	ExpectedAccount    string
	ExpectedCollection string
	ExpectedLiker      string
}

func TestFollwer(t *testing.T) {
	setup()
	RunSQL(DB, `
		INSERT INTO accounts (id, title, content, images, likers, follwers) VALUES (1, "Account_1", "Account_1", "1,2,3", 5, 6);
		INSERT INTO accounts (id, title, content, images, likers, follwers) VALUES (2, "Account_2", "Account_2", "1,2,3", 5, 6);
		INSERT INTO collections (id, user_id, account_id) VALUES (1, 1, 2);
		INSERT INTO collections (id, user_id, account_id) VALUES (2, 1, 3);
		INSERT INTO likers (id, user_id, account_id, is_follwer) VALUES (1, 1, 2, 1);
	`)
	testCases := []FollwerTestCase{
		{Token: "123456", ExpectedError: `invalid user`},
		{Token: "123123", ID: "0", ExpectedError: `invalid params`},
		{Token: "123123", ID: "1", Follwer: "1", ExpectedResponse: `{"code":200,"data":7,"message":null}`, ExpectedAccount: `1,Account_1,Account_1,1,2,3,5,7`, ExpectedCollection: `1,1,2; 2,1,3; 3,1,1`, ExpectedLiker: `1,1,2,1; 2,1,1,1`},
		{Token: "123123", ID: "1", Follwer: "1", ExpectedError: `{"code":400,"data":null,"message":"invalid operation"}`},
		{Token: "123123", ID: "2", Follwer: "0", ExpectedAccount: `2,Account_2,Account_2,1,2,3,5,5`, ExpectedCollection: `2,1,3; 3,1,1`, ExpectedLiker: `1,1,2,0; 2,1,1,1`},
	}
	for i, testCase := range testCases {
		body := SingeGet(testCase.Token, "/api/follwer", url.Values{
			"id":      {testCase.ID},
			"follwer": {testCase.Follwer},
		})
		if testCase.ExpectedError != "" && !strings.Contains(body, testCase.ExpectedError) {
			t.Errorf(color.RedString("TestFollwer #%v: expected error %v but got %v", i+1, testCase.ExpectedError, body))
		}
		if testCase.ExpectedResponse != "" && body != testCase.ExpectedResponse {
			t.Errorf(color.RedString("TestFollwer #%v: expected response %v but got %v", i+1, testCase.ExpectedResponse, body))
		}
		account := ""
		if testCase.ID != "0" {
			account = GetRecords(DB, "accounts", "id, title, content, images, likers, follwers", fmt.Sprintf("where id = %s", testCase.ID))
		}
		if testCase.ExpectedAccount != "" && account != testCase.ExpectedAccount {
			t.Errorf(color.RedString("TestFollwer #%v: expected account %v but got %v", i+1, testCase.ExpectedAccount, account))
		}
		collection := GetRecords(DB, "collections", "id, user_id, account_id")
		if testCase.ExpectedCollection != "" && collection != testCase.ExpectedCollection {
			t.Errorf(color.RedString("TestFollwer #%v: expected collection %v but got %v", i+1, testCase.ExpectedCollection, collection))
		}
		liker := GetRecords(DB, "likers", "id, user_id, account_id, is_liked, is_follwer")
		if testCase.ExpectedLiker != "" && liker != testCase.ExpectedLiker {
			t.Errorf(color.RedString("TestFollwer #%v: expected liker %v but got %v", i+1, testCase.ExpectedLiker, liker))
		}
		color.Green("TestFollwer #%v: success", i+1)
	}
}
