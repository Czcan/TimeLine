package test

import (
	"net/url"
	"strings"
	"testing"

	"github.com/fatih/color"
)

type CommentTestCase struct {
	Token            string
	Content          string
	ID               string
	ExpectedComments string
	ExpectedError    string
}

func TestComment(t *testing.T) {
	setup()

	testCases := []CommentTestCase{
		{Token: "123123", ID: "1", Content: "Comment_1", ExpectedComments: `{"code":200,"data":[{"nick_name":"name","content":"Comment_1","avatar_url":"","date":0}],"message":null}`},
		{Token: "123456", ExpectedError: `invalid user`},
		{Token: "123123", ID: "0", ExpectedError: "invalid params"},
	}
	for i, testCase := range testCases {
		body := SingePost(testCase.Token, "/api/comment", url.Values{
			"id":      {testCase.ID},
			"content": {testCase.Content},
		})
		if testCase.ExpectedError != "" && !strings.Contains(body, testCase.ExpectedError) {
			t.Errorf(color.RedString("TestComment #%v: expected error %v but got %v", i+1, testCase.ExpectedError, body))
		}
		if testCase.ExpectedComments != "" && body != testCase.ExpectedComments {
			t.Errorf(color.RedString("TestComment #%v: expected comments %v but got %v", i+1, testCase.ExpectedComments, body))
		}
		color.Green("TestComment #%v: success", i+1)
	}
}
