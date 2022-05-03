package models

type Liker struct {
	ID      int
	UserID  int
	NoteID  int
	IsLiked bool
}
