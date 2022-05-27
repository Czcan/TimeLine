package models

import (
	"gorm.io/gorm"
)

type Note struct {
	ID        int            `json:"id"`
	UserID    int            `json:"user_id"`
	FolderID  int            `json:"folder_id"`
	Content   string         `json:"content"`
	Status    bool           `json:"status"`
	CreatedAt int            `json:"created_at"`
	UpdatedAt int            `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-"`
}

func CreateNote(db *gorm.DB, folderID, userID int, content string) ([]Note, error) {
	note := &Note{FolderID: folderID, Content: content, UserID: userID}
	if err := db.Save(&note).Error; err != nil {
		return nil, err
	}
	notes := GetNoteList(db, folderID, userID)
	return notes, nil
}

func UpdateNoteStatus(db *gorm.DB, noteID, userID int, status bool) error {
	if err := db.Model(&Note{}).Where("id = ? AND user_id = ?", noteID, userID).
		Update("status", status).Error; err != nil {
		return err
	}
	return nil
}

func DeleteNote(db *gorm.DB, noteID int, userID, folderID int) ([]Note, error) {
	if err := db.Where("id = ?", noteID).Delete(&Note{}).Error; err != nil {
		return nil, err
	}
	notes := GetNoteList(db, userID, folderID)
	return notes, nil
}

func GetNoteList(db *gorm.DB, userID int, folderID int) []Note {
	notes := []Note{}
	db.Where("folder_id = ? AND user_id = ? AND deleted_at is null", folderID, userID).Find(&notes)
	return notes
}
