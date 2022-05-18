package test

import (
	"net/url"
	"strings"
	"testing"

	"github.com/fatih/color"
)

type FolderListTestCase struct {
	Token           string
	ExpectedFolders string
	ExpectedError   string
}

func TestFolderList(t *testing.T) {
	setup()
	RunSQL(DB, `
		INSERT INTO folders (id, user_id, name) VALUES (1, 1, 'folder1');
		INSERT INTO folders (id, user_id, name) VALUES (2, 1, 'folder2');
	`)
	testCases := []FolderListTestCase{
		{Token: "123123", ExpectedFolders: `{"code":200,"data":[{"id":1,"name":"folder1","user_id":1,"create_at":0,"updated_at":0},{"id":2,"name":"folder2","user_id":1,"create_at":0,"updated_at":0}],"message":null}`},
		{Token: "123456", ExpectedError: `invalid user`},
	}
	for i, testCase := range testCases {
		body := SingeGet(testCase.Token, "/api/folder/list", nil)
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
	Name            string
	ExpectedFolders string
	ExpectedError   string
}

func TestFolderCreate(t *testing.T) {
	setup()
	RunSQL(DB, `
		INSERT INTO folders (id, user_id, name) VALUES (1, 1, 'folder1');
		INSERT INTO folders (id, user_id, name) VALUES (2, 1, 'folder2');
	`)
	testCases := []FolderCreateTestCase{
		{Token: "123123", Name: "folder3", ExpectedFolders: `{"code":200,"data":[{"id":1,"name":"folder1","user_id":1,"create_at":0,"updated_at":0},{"id":2,"name":"folder2","user_id":1,"create_at":0,"updated_at":0},{"id":3,"name":"folder3","user_id":1,"create_at":0,"updated_at":0}],"message":null}`},
		{Token: "123123", ExpectedError: `invalid params`},
		{Token: "123456", ExpectedError: `invalid user`},
	}
	for i, testCase := range testCases {
		body := SingePost(testCase.Token, "/api/folder/create", url.Values{
			"name": {testCase.Name},
		})
		if testCase.ExpectedError != "" && !strings.Contains(body, testCase.ExpectedError) {
			t.Errorf(color.RedString("TestFolderCreate #%v: expected error %v but got %v", i+1, testCase.ExpectedError, body))
		}
		if testCase.ExpectedFolders != "" && body != testCase.ExpectedFolders {
			t.Errorf(color.RedString("TestFolderCreate #%v: expected folders %v but got %v", i+1, testCase.ExpectedFolders, body))
		}
		color.Green("TestFolderCreate #%v: Success", i+1)
	}
}

type FolderDeletedTestCase struct {
	Token          string
	ID             string
	ExpectedFolder string
	ExpectedNote   string
	ExpectedError  string
}

func TestDeletedFolder(t *testing.T) {
	setup()
	RunSQL(DB, `
		INSERT INTO folders (id, user_id, name) VALUES (1, 1, 'folder1');
		INSERT INTO folders (id, user_id, name) VALUES (2, 1, 'folder2');
		INSERT INTO notes (id, user_id, folder_id, content) VALUES (1, 1, 1, "note1");
		INSERT INTO notes (id, user_id, folder_id, content) VALUES (2, 1, 1, "note2");
		INSERT INTO notes (id, user_id, folder_id, content) VALUES (3, 1, 2, "note3");
	`)
	testCases := []FolderDeletedTestCase{
		{Token: "123123", ID: "1", ExpectedFolder: `{"code":200,"data":[{"id":2,"name":"folder2","user_id":1,"create_at":0,"updated_at":0}],"message":null}`, ExpectedNote: `3,1,2,note3`},
		{Token: "invalid user", ExpectedError: "invalid user"},
		{Token: "123123", ID: "invalid id", ExpectedError: "invalid params"},
	}
	for i, testCase := range testCases {
		body := SingeDelete(testCase.Token, "/api/folder/deleted", url.Values{
			"folder_id": {testCase.ID},
		})
		if testCase.ExpectedError != "" && !strings.Contains(body, testCase.ExpectedError) {
			t.Errorf(color.RedString("TestDeletedFolder #%v: expected error %v but got %v", i+1, testCase.ExpectedError, body))
		}
		if testCase.ExpectedFolder != "" && body != testCase.ExpectedFolder {
			t.Errorf(color.RedString("TestDeletedFolder #%v: expected folders %v but got %v", i+1, testCase.ExpectedFolder, body))
		}
		notes := GetRecords(DB, "notes", "id, user_id, folder_id, content", "WHERE deleted_at is null")
		if testCase.ExpectedNote != "" && notes != testCase.ExpectedNote {
			t.Errorf(color.RedString("TestDeletedFolder #%v: expected notes %v but got %v", i+1, testCase.ExpectedNote, notes))
		}
		color.Green("TestDeletedFolder #%v: Success", i+1)
	}
}
