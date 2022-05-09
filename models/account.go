package models

import "time"

type Account struct {
	ID        int
	UserID    int
	Images    string
	Title     string
	Content   string
	Likers    int
	Follwers  int
	CreatedAt time.Time
	UpdatedAt time.Time
}
