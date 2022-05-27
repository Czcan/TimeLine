package models

import (
	"time"

	"gorm.io/gorm"
)

type Collection struct {
	ID        int
	UserID    int
	AccountID int
	CreatedAt time.Time
}

func FindCollection(db *gorm.DB, userID int) []Account {
	accounts := []Account{}
	db.
		Model(&Account{}).
		Select("accounts.id AS id, title, content, accounts.created_at, images, likers, follwers, users.nick_name, CONCAT('upload/avatar/images/', users.id, '.jpg') AS avatar_url").
		Joins("LEFT JOIN collections ON accounts.id = collections.account_id").
		Joins("LEFT JOIN users ON accounts.user_id = users.id").
		Where("collections.user_id = ?", userID).
		Find(&accounts)
	return accounts
}

func SaveCollection(db *gorm.DB, userID int, accountID int) error {
	collection := &Collection{UserID: userID, AccountID: accountID}
	if err := db.Save(&collection).Error; err != nil {
		return err
	}
	return nil
}
