package models

import (
	"gorm.io/gorm"
)

type Folder struct {
	ID        int            `json:"id"`
	Name      string         `json:"name"`
	UserID    int            `json:"user_id"`
	Kind      int            `json:"kind"`
	CreatedAt int            `json:"created_at"`
	UpdatedAt int            `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-"`
}

func CreateFolder(db *gorm.DB, userID int, name string, kind int) ([]Folder, error) {
	folder := &Folder{Name: name, UserID: userID, Kind: kind}
	if err := db.Save(&folder).Error; err != nil {
		return nil, err
	}
	folders := GetFolderList(db, userID, kind)
	return folders, nil
}

func GetFolderList(db *gorm.DB, userID int, kind int) []Folder {
	folders := []Folder{}
	db.Where("user_id = ? AND kind = ? AND deleted_at is null", userID, kind).Find(&folders)
	return folders
}

func DeletedFolder(db *gorm.DB, folderID int, userID int, kind int) ([]Folder, error) {
	// err := database.Transaction(db, func(tx *gorm.DB) error {
	// 	if err := db.Where("id = ?", folderID).Delete(&Folder{}).Error; err != nil {
	// 		return err
	// 	}
	// 	if err := db.Where("folder_id = ?", folderID).Delete(&Note{}).Error; err != nil {
	// 		return err
	// 	}
	// 	return nil
	// })
	if err := db.Where("id = ?", folderID).Delete(&Folder{}).Error; err != nil {
		return nil, err
	}
	folders := GetFolderList(db, userID, kind)
	return folders, nil
}
