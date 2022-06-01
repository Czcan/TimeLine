package test

import (
	"net/url"
	"strings"
	"testing"

	"github.com/fatih/color"
)

type TaskListTestCase struct {
	Token         string
	FolderID      string
	ExpectedTask  string
	ExpectedError string
}

func TestTaskList(t *testing.T) {
	setup()
	RunSQL(DB, `
		INSERT INTO tasks (id, user_id, folder_id, content, description) VALUES (1, 1, 1, "content1", "desc1");
		INSERT INTO tasks (id, user_id, folder_id, content, description) VALUES (2, 1, 2, "content2", "desc2");
		INSERT INTO tasks (id, user_id, folder_id, content, description, deleted_at) VALUES (3, 1, 2, "content3", "desc3", "2022-05-18 15:36:53")
	`)
	testCases := []TaskListTestCase{
		{Token: "123123", FolderID: "1", ExpectedTask: `{"code":200,"data":[{"id":1,"folder_id":1,"content":"content1","description":"desc1","status":0,"start_at":0}],"message":null}`},
		{Token: "123123", FolderID: "2", ExpectedTask: `{"code":200,"data":[{"id":2,"folder_id":2,"content":"content2","description":"desc2","status":0,"start_at":0}],"message":null}`},
		{Token: "123123", FolderID: "0", ExpectedError: `invalid param`},
		{Token: "123456", ExpectedError: `invalid user`},
	}
	for i, testCase := range testCases {
		body := SingeGet(testCase.Token, "/api/task/list", url.Values{
			"folder_id": {testCase.FolderID},
		})
		if testCase.ExpectedError != "" && !strings.Contains(body, testCase.ExpectedError) {
			t.Errorf(color.RedString("TestTaskList #%v: expected error %v but got error %v", i+1, testCase.ExpectedError, body))
		}
		if testCase.ExpectedTask != "" && body != testCase.ExpectedTask {
			t.Errorf(color.RedString("TestTaskList #%v: expected tasks %v but got notes %v", i+1, testCase.ExpectedTask, body))
		}
		color.Green("TestTaskList #%v: Success", i+1)
	}
}

type TaskCreateTestCase struct {
	Token         string
	FolderID      string
	Content       string
	Desc          string
	StartAt       string
	ExpectedError string
	ExpectedTasks string
}

func TestTaskCreate(t *testing.T) {
	setup()
	RunSQL(DB, `
		INSERT INTO folders (id, user_id, name) VALUES (1, 1, 'folder1');
		INSERT INTO folders (id, user_id, name) VALUES (2, 1, 'folder2');
		INSERT INTO tasks (id, user_id, folder_id, content, description) VALUES (1, 1, 1, "content1", "desc1");
		INSERT INTO tasks (id, user_id, folder_id, content, description) VALUES (2, 1, 2, "content2", "desc2");
	`)
	testCases := []TaskCreateTestCase{
		{Token: "123123", FolderID: "1", Content: "content3", Desc: "desc3", StartAt: "202061", ExpectedTasks: `{"code":200,"data":[{"id":3,"folder_id":1,"content":"content3","description":"desc3","status":0,"start_at":202061},{"id":1,"folder_id":1,"content":"content1","description":"desc1","status":0,"start_at":0}],"message":null}`},
		{Token: "123123", FolderID: "invalid param", Content: "", ExpectedError: "invalid param"},
		{Token: "invalid token", ExpectedError: `invalid user`},
	}
	for i, testCase := range testCases {
		body := SingePost(testCase.Token, "/api/task/create", url.Values{
			"folder_id": {testCase.FolderID},
			"content":   {testCase.Content},
			"desc":      {testCase.Desc},
			"start_at":  {testCase.StartAt},
		})
		if testCase.ExpectedError != "" && !strings.Contains(body, testCase.ExpectedError) {
			t.Errorf(color.RedString("TestTaskCreate #%v: expected error %v but got error %v", i+1, testCase.ExpectedError, body))
		}
		body = ExtractDate(body, "data.#.created_at", "0")
		if testCase.ExpectedTasks != "" && body != testCase.ExpectedTasks {
			t.Errorf(color.RedString("TestTaskCreate #%v: expected tasks %v but got notes %v", i+1, testCase.ExpectedTasks, body))
		}
		color.Green("TestTaskCreate #%v: Success", i+1)
	}
}

type DeletedTaskTestCase struct {
	Token         string
	ID            string
	FolderID      string
	ExpectedTasks string
	ExpectedError string
}

func TestDeletedTask(t *testing.T) {
	setup()
	RunSQL(DB, `
		INSERT INTO tasks (id, user_id, folder_id, content, description) VALUES (1, 1, 1, "content1", "desc1");
		INSERT INTO tasks (id, user_id, folder_id, content, description) VALUES (2, 1, 1, "content2", "desc2");
	`)
	testCases := []DeletedTaskTestCase{
		{Token: "123123", ID: "1", FolderID: "1", ExpectedTasks: `{"code":200,"data":[{"id":2,"folder_id":1,"content":"content2","description":"desc2","status":0,"start_at":0}],"message":null}`},
		{Token: "invalid user", ExpectedError: "invalid user"},
		{Token: "123123", ID: "invalid id", ExpectedError: "invalid params"},
	}
	for i, testCase := range testCases {
		body := SingeDelete(testCase.Token, "/api/task/deleted", url.Values{
			"task_id":   {testCase.ID},
			"folder_id": {testCase.FolderID},
		})
		if testCase.ExpectedError != "" && !strings.Contains(body, testCase.ExpectedError) {
			t.Errorf(color.RedString("TestDeletedTask #%v: expected error %v but got %v", i+1, testCase.ExpectedError, body))
		}
		if testCase.ExpectedTasks != "" && body != testCase.ExpectedTasks {
			t.Errorf(color.RedString("TestDeletedTask #%v: expected notes %v but got %v", i+1, testCase.ExpectedTasks, body))
		}
		color.Green("TestDeletedTask #%v: success", i+1)
	}
}
