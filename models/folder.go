package models

import (
	"time"

	"github.com/Czcan/TimeLine/utils/database"
	"github.com/jinzhu/gorm"
)

type Folder struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	UserID    int       `json:"user_id"`
	CreatedAt int       `json:"create_at"`
	UpdatedAt int       `json:"updated_at"`
	DeletedAt time.Time `json:"-"`
}

func GetFolderList(db *gorm.DB, userID int) []Folder {
	folders := []Folder{}
	db.Where("user_id = ?", userID).Find(&folders)
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
