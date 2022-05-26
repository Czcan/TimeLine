package test

import (
	"fmt"
	"sync"
	"testing"

	"github.com/fatih/color"
)

func TestLikerIdenpt(t *testing.T) {
	setup()
	RunSQL(DB, `
		INSERT INTO accounts (id, title, content, images, likers, follwers) VALUES (1, "Account_1", "Account_1", "1,2,3", 5, 6);
		INSERT INTO likers (id, user_id, account_id, is_liked) VALUES (1, 1, 1, 0);
	`)
	var wg sync.WaitGroup
	wg.Add(2)
	for i := 0; i < 2; i++ {
		go func() {
			defer wg.Done()
			body := SingeGet("123123", fmt.Sprintf("/api/liker?id=%s&liker=%s", "1", "1"), nil)
			likers := GetRecords(DB, "likers", "id, user_id, account_id, is_liked, updated_at")
			color.Green(body)
			color.Green(likers)
		}()
	}
	wg.Wait()
}
