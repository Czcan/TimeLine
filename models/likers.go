package models

import "time"

type Liker struct {
	ID        int
	UserID    int
	NoteID    int
	IsLiked   bool
	CreatedAt time.Time
}

type Collection struct {
	ID        int
	UserID    int
	AccountID int
	CreatedAt time.Time
}
