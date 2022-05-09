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

type NoteListTestCase struct {
	Token         string
	FolderID      string
	ExpectedNotes string
	ExpectedError string
}

func TestNoteList(t *testing.T) {
	setup()
	RunSQL(DB, `
		INSERT INTO folders (id, user_id, name) VALUES (1, 1, 'folder1');
		INSERT INTO folders (id, user_id, name) VALUES (2, 1, 'folder2');
		INSERT INTO notes (id, user_id, folder_id, content) VALUES (1, 1, 1, "note1");
		INSERT INTO notes (id, user_id, folder_id, content) VALUES (2, 1, 2, "note2");
	`)
	testCases := []NoteListTestCase{
		{Token: "123123", FolderID: "1", ExpectedNotes: `{"code":200,"data":[{"id":1,"user_id":1,"folder_id":1,"content":"note1","status":false,"created_at":0}],"message":null}`},
		{Token: "123123", FolderID: "2", ExpectedNotes: `{"code":200,"data":[{"id":2,"user_id":1,"folder_id":2,"content":"note2","status":false,"created_at":0}],"message":null}`},
		{Token: "123123", FolderID: "0", ExpectedError: `invalid param`},
		{Token: "123456", ExpectedError: `invalid user`},
	}
	for i, testCase := range testCases {
		body := SingeGet(testCase.Token, "/api/note/list", url.Values{
			"folder_id": {testCase.FolderID},
		})
		if testCase.ExpectedError != "" && !strings.Contains(body, testCase.ExpectedError) {
			t.Errorf(color.RedString("TestNoteList #%v: expected error %v but got error %v", i+1, testCase.ExpectedError, body))
		}
		if testCase.ExpectedNotes != "" && body != testCase.ExpectedNotes {
			t.Errorf(color.RedString("TestNoteList #%v: expected notes %v but got notes %v", i+1, testCase.ExpectedNotes, body))
		}
		color.Green("TestNoteList #%v: Success", i+1)
	}
}

type NoteCreateTestCase struct {
	Token         string
	FolderID      string
	Content       string
	ExpectedError string
	ExpectedNotes string
}

func TestNoteCreate(t *testing.T) {
	setup()
	RunSQL(DB, `
		INSERT INTO folders (id, user_id, name) VALUES (1, 1, 'folder1');
		INSERT INTO folders (id, user_id, name) VALUES (2, 1, 'folder2');
		INSERT INTO notes (id, user_id, folder_id, content) VALUES (1, 1, 1, "note1");
		INSERT INTO notes (id, user_id, folder_id, content) VALUES (2, 1, 2, "note2");
	`)
	testCases := []NoteCreateTestCase{
		{Token: "123123", FolderID: "1", Content: "note3", ExpectedNotes: `{"code":200,"data":[{"id":1,"user_id":1,"folder_id":1,"content":"note1","status":false,"created_at":0},{"id":3,"user_id":1,"folder_id":1,"content":"note3","status":false,"created_at":0}],"message":null}`},
		{Token: "123123", FolderID: "invalid param", Content: "", ExpectedError: "invalid param"},
		{Token: "invalid token", ExpectedError: `invalid user`},
	}
	for i, testCase := range testCases {
		body := SingePost(testCase.Token, "/api/note/create", url.Values{
			"folder_id": {testCase.FolderID},
			"content":   {testCase.Content},
		})
		if testCase.ExpectedError != "" && !strings.Contains(body, testCase.ExpectedError) {
			t.Errorf(color.RedString("TestNoteCreate #%v: expected error %v but got error %v", i+1, testCase.ExpectedError, body))
		}
		if testCase.ExpectedNotes != "" && body != testCase.ExpectedNotes {
			t.Errorf(color.RedString("TestNoteCreate #%v: expected notes %v but got notes %v", i+1, testCase.ExpectedNotes, body))
		}
		color.Green("TestNoteCreate #%v: Success", i+1)
	}
}

type NoteUpdateTestCase struct {
	Token         string
	NoteID        string
	Status        string
	ExpectedError string
	ExpectedNotes string
}

func TestUpdateNote(t *testing.T) {
	setup()
	RunSQL(DB, `
		INSERT INTO folders (id, user_id, name) VALUES (1, 1, 'folder1');
		INSERT INTO folders (id, user_id, name) VALUES (2, 1, 'folder2');
		INSERT INTO notes (id, user_id, folder_id, content) VALUES (1, 1, 1, "note1");
		INSERT INTO notes (id, user_id, folder_id, content, status) VALUES (2, 1, 2, "note2", 1);
	`)
	testCases := []NoteUpdateTestCase{
		{Token: "123123", NoteID: "1", Status: "1", ExpectedNotes: `{"code":200,"data":"updated successed"}`},
		{Token: "123123", NoteID: "2", Status: "0", ExpectedNotes: `{"code":200,"data":"updated successed"}`},
		{Token: "123456", ExpectedError: `invalid user`},
	}
	for i, testCase := range testCases {
		body := SingePost(testCase.Token, "/api/note/update", url.Values{
			"note_id": {testCase.NoteID},
			"status":  {testCase.Status},
		})
		if testCase.ExpectedError != "" && !strings.Contains(body, testCase.ExpectedError) {
			t.Errorf("TestUpdateNote #%v: expected error %v but got error %v", i+1, testCase.ExpectedError, body)
		}
		if testCase.ExpectedNotes != "" && body != testCase.ExpectedNotes {
			t.Errorf("TestUpdateNote #%v: expected notes %v but got notes %v", i+1, testCase.ExpectedNotes, body)
		}
		color.Green("TestUpdateNote #%v: Success", i+1)
	}
}
