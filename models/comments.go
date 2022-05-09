package models

import "time"

type Comments struct {
	ID        int
	AccoutID  int
	UserID    int
	Content   string
	CreatedAt time.Time
}