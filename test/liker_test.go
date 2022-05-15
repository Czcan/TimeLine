package test

import (
	"fmt"
	"net/url"
	"strings"
	"testing"

	"github.com/fatih/color"
)

type LikersTestCase struct {
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
		{ID: "0", ExpectedError: `invalid params`},
		{ID: "1", Liker: "1", ExpectedReponse: `{"code":200,"data":6,"message":null}`, ExpectedAccount: `1,Account_1,Account_1,1,2,3,6,6`},
	}
	for i, testCase := range testCases {
		body := Get(fmt.Sprintf("/api/liker?id=%s&liker=%s", testCase.ID, testCase.Liker))
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
}

func TestFollwer(t *testing.T) {
	setup()
	RunSQL(DB, `
		INSERT INTO accounts (id, title, content, images, likers, follwers) VALUES (1, "Account_1", "Account_1", "1,2,3", 5, 6);
	`)
	testCases := []FollwerTestCase{
		{Token: "123456", ExpectedError: `invalid user`},
		{Token: "123123", ID: "0", ExpectedError: `invalid parmas`},
		{Token: "123123", ID: "1", Follwer: "1", ExpectedResponse: `dsada`, ExpectedAccount: `dsada`, ExpectedCollection: `sadas`},
	}
	for i, testCase := range testCases {
		body := SingeGet(testCase.Token, "/api/follwer", url.Values{
			"id":      {testCase.ID},
			"follwer": {testCase.Follwer},
		})
		if testCase.ExpectedError != "" && !strings.Contains(body, testCase.ExpectedError) {
			t.Errorf("TestFollwer #%v: expected error %v but got %v", i+1, testCase.ExpectedError, body)
		}
		if testCase.ExpectedResponse != "" && body != testCase.ExpectedResponse {
			t.Errorf("TestFollwer #%v: expected response %v but got %v", i+1, testCase.ExpectedResponse, body)
		}
		account := GetRecords(DB, "accounts", "id, title, content, images, likers, follwers", "where id = 1")
		if testCase.ExpectedAccount != "" && account != testCase.ExpectedAccount {
			t.Errorf("TestFollwer #%v: expected account %v but got %v", i+1, testCase.ExpectedAccount, account)
		}
		collection := GetRecords(DB, "collections", "id, user_id, account_id")
		if testCase.ExpectedCollection != "" && collection != testCase.ExpectedCollection {
			t.Errorf("TestFollwer #%v: expected collection %v but got %v", i+1, testCase.ExpectedCollection, collection)
		}
	}
}
