package models

import "time"

type NoteFolder struct {
	ID        int
	Name      string
	UserID    int
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Note struct {
	ID         int
	UserID     int
	FolderID   int
	Content    string
	IsFinished bool
	CreatedAt  time.Time
}
