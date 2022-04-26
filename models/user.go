package models

import "time"

type User struct {
	ID        int
	Avatar    string
	NickName  string
	Uid       string
	Gender    bool
	Age       int
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}
