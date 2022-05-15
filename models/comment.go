package models

import "time"

type Comment struct {
	ID        int
	AccountID int
	UserID    int
	Content   string
	CreatedAt time.Time
}
