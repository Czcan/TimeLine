package test

import (
	"net/url"
	"strings"
	"testing"

	"github.com/fatih/color"
)

type FolderListTestCase struct {
	Token           string
	Kind            string
	ExpectedFolders string
	ExpectedError   string
}

func TestFolderList(t *testing.T) {
	setup()
	RunSQL(DB, `
		INSERT INTO folders (id, user_id, name, kind) VALUES (1, 1, 'folder1', 0);
		INSERT INTO folders (id, user_id, name, kind) VALUES (2, 1, 'folder2', 1);
	`)
	testCases := []FolderListTestCase{
		{Token: "123123", Kind: "1", ExpectedFolders: `{"code":200,"data":[{"id":2,"name":"folder2","user_id":1,"kind":1,"created_at":0}],"message":null}`},
		{Token: "123123", Kind: "0", ExpectedFolders: `{"code":200,"data":[{"id":1,"name":"folder1","user_id":1,"kind":0,"created_at":0}],"message":null}`},
		{Token: "123456", ExpectedError: `invalid user`},
	}
	for i, testCase := range testCases {
		body := SingeGet(testCase.Token, "/api/folder/list", url.Values{
			"kind": {testCase.Kind},
		})
		if testCase.ExpectedError != "" && !strings.Contains(body, testCase.ExpectedError) {
			t.Errorf(color.RedString("TestFolderList #%v: Expected error %v but got %v", i+1, testCase.ExpectedError, body))
		}
		if testCase.ExpectedFolders != "" && body != testCase.ExpectedFolders {
			t.Errorf(color.RedString("TestFolderList #%v: Expected folders %v but got %v", i+1, testCase.ExpectedFolders, body))
		}
		color.Green("TestFolderList #%v: Success", i+1)
	}
}

type FolderCreateTestCase struct {
	Token           string
	Kind            string
	Name            string
	ExpectedFolders string
	ExpectedError   string
}

func TestFolderCreate(t *testing.T) {
	setup()
	RunSQL(DB, `
		INSERT INTO folders (id, user_id, name, kind) VALUES (1, 1, 'folder1', 1);
		INSERT INTO folders (id, user_id, name, kind) VALUES (2, 1, 'folder2', 0);
	`)
	testCases := []FolderCreateTestCase{
		{Token: "123123", Kind: "1", Name: "folder3", ExpectedFolders: `{"code":200,"data":[{"id":1,"name":"folder1","user_id":1,"kind":1,"created_at":0},{"id":3,"name":"folder3","user_id":1,"kind":1,"created_at":0}],"message":null}`},
		{Token: "123123", Kind: "0", Name: "folder4", ExpectedFolders: `{"code":200,"data":[{"id":2,"name":"folder2","user_id":1,"kind":0,"created_at":0},{"id":4,"name":"folder4","user_id":1,"kind":0,"created_at":0}],"message":null}`},
		{Token: "123123", ExpectedError: `invalid params`},
		{Token: "123456", ExpectedError: `invalid user`},
	}
	for i, testCase := range testCases {
		body := SingePost(testCase.Token, "/api/folder/create", url.Values{
			"name": {testCase.Name},
			"kind": {testCase.Kind},
		})
		if testCase.ExpectedError != "" && !strings.Contains(body, testCase.ExpectedError) {
			t.Errorf(color.RedString("TestFolderCreate #%v: expected error %v but got %v", i+1, testCase.ExpectedError, body))
		}
		body = ExtractDate(body, "data.#.created_at", "0")
		if testCase.ExpectedFolders != "" && body != testCase.ExpectedFolders {
			t.Errorf(color.RedString("TestFolderCreate #%v: expected folders %v but got %v", i+1, testCase.ExpectedFolders, body))
		}
		color.Green("TestFolderCreate #%v: Success", i+1)
	}
}

type FolderDeletedTestCase struct {
	Token          string
	ID             string
	Kind           string
	ExpectedFolder string
	ExpectedError  string
}

func TestDeletedFolder(t *testing.T) {
	setup()
	RunSQL(DB, `
		INSERT INTO folders (id, user_id, name, kind) VALUES (1, 1, 'folder1', 1);
		INSERT INTO folders (id, user_id, name, kind) VALUES (3, 1, 'folder3', 1);
		INSERT INTO folders (id, user_id, name) VALUES (2, 1, 'folder2');
		INSERT INTO notes (id, user_id, folder_id, content) VALUES (1, 1, 1, "note1");
		INSERT INTO notes (id, user_id, folder_id, content) VALUES (2, 1, 1, "note2");
		INSERT INTO notes (id, user_id, folder_id, content) VALUES (3, 1, 2, "note3");
	`)
	testCases := []FolderDeletedTestCase{
		{Token: "123123", Kind: "1", ID: "1", ExpectedFolder: `{"code":200,"data":[{"id":3,"name":"folder3","user_id":1,"kind":1,"created_at":0}],"message":null}`},
		{Token: "invalid user", ExpectedError: "invalid user"},
		{Token: "123123", ID: "invalid id", ExpectedError: "invalid params"},
	}
	for i, testCase := range testCases {
		body := SingeDelete(testCase.Token, "/api/folder/deleted", url.Values{
			"folder_id": {testCase.ID},
			"kind":      {testCase.Kind},
		})
		if testCase.ExpectedError != "" && !strings.Contains(body, testCase.ExpectedError) {
			t.Errorf(color.RedString("TestDeletedFolder #%v: expected error %v but got %v", i+1, testCase.ExpectedError, body))
		}
		if testCase.ExpectedFolder != "" && body != testCase.ExpectedFolder {
			t.Errorf(color.RedString("TestDeletedFolder #%v: expected folders %v but got %v", i+1, testCase.ExpectedFolder, body))
		}
		color.Green("TestDeletedFolder #%v: Success", i+1)
	}
}
