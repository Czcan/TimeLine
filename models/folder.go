package models

import (
	"github.com/Czcan/TimeLine/utils/database"
	"gorm.io/gorm"
)

type Folder struct {
	ID        int            `json:"id"`
	Name      string         `json:"name"`
	UserID    int            `json:"user_id"`
	CreatedAt int            `json:"created_at"`
	UpdatedAt int            `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-"`
}

func GetFolderList(db *gorm.DB, userID int) []Folder {
	folders := []Folder{}
	db.Where("user_id = ? AND deleted_at is null", userID).Find(&folders)
	if len(folders) == 0 {
		return nil
	}
	return folders
}

func DeletedFolderAndNote(db *gorm.DB, folderID int) error {
	database.Transaction(db, func(tx *gorm.DB) error {
		if err := db.Where("id = ?", folderID).Delete(&Folder{}).Error; err != nil {
			return err
		}
		if err := db.Where("folder_id = ?", folderID).Delete(&Note{}).Error; err != nil {
			return err
		}
		return nil
	})
	return nil
}
