package models

import (
	"github.com/jinzhu/gorm"
)

type Folder struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	UserID    int    `json:"user_id"`
	CreatedAt int    `json:"create_at"`
	UpdatedAt int    `json:"updated_at"`
}

type Note struct {
	ID        int    `json:"id"`
	UserID    int    `json:"user_id"`
	FolderID  int    `json:"folder_id"`
	Content   string `json:"content"`
	Status    bool   `json:"status"`
	CreatedAt int    `json:"created_at"`
}

func GetFolderList(db *gorm.DB, userID int) []Folder {
	folders := []Folder{}
	db.Where("user_id = ?", userID).Find(&folders)
	if len(folders) == 0 {
		return nil
	}
	return folders
}

func GetNoteList(db *gorm.DB, userID int, folderID int) []Note {
	notes := []Note{}
	db.Where("folder_id = ? AND user_id = ?", folderID, userID).Find(&notes)
	if len(notes) == 0 {
		return nil
	}
	return notes
}
