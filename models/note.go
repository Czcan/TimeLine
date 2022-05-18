package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Note struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	FolderID  int       `json:"folder_id"`
	Content   string    `json:"content"`
	Status    bool      `json:"status"`
	CreatedAt int       `json:"created_at"`
	UpdatedAt int       `json:"updated_at"`
	DeletedAt time.Time `json:"-"`
}

func GetNoteList(db *gorm.DB, userID int, folderID int) []Note {
	notes := []Note{}
	db.Where("folder_id = ? AND user_id = ?", folderID, userID).Find(&notes)
	if len(notes) == 0 {
		return nil
	}
	return notes
}


