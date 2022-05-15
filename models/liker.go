package models

import (
	"time"

	"github.com/Czcan/TimeLine/utils/database"
	"github.com/jinzhu/gorm"
)

type Liker struct {
	ID        int
	UserID    int
	NoteID    int
	IsLiked   bool
	CreatedAt time.Time
}

func UpdateFollwerAndSyncCollection(db *gorm.DB, accountID, userID int, liker bool) (int, error) {
	var (
		account    = &Account{}
		collection = &Collection{}
	)
	err := database.Transaction(db, func(tx *gorm.DB) error {
		if err := tx.Where("id = ?", accountID).First(&account).Error; err != nil {
			return err
		}
		if liker {
			account.Follwers += 1
		} else {
			account.Follwers -= 1
		}
		if err := tx.Save(&account).Error; err != nil {
			return err
		}
		collection = &Collection{UserID: userID, AccountID: accountID}
		if err := tx.Save(&collection).Error; err != nil {
			return err
		}
		return nil
	})
	return account.Follwers, err
}
