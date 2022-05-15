package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Liker struct {
	ID        int
	UserID    int
	NoteID    int
	IsLiked   bool
	CreatedAt time.Time
}

func UpdateFollwerAndSyncCollection(db *gorm.DB, accountID, userID int, liker bool) int {
	var (
		account    = &Account{}
		collection = &Collection{}
	)
	tx := db.Begin()
	tx.Where("id = ?", accountID).First(&account)
	if liker {
		account.Follwers += 1
	} else {
		account.Follwers -= 1
	}
	tx.Save(&account)
	collection = &Collection{UserID: userID, AccountID: accountID}
	tx.Save(&collection)
	if tx.Error != nil {
		tx.Rollback()
	}
	tx.Commit()
	return account.Follwers
}
