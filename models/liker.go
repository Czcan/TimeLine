package models

import (
	"errors"
	"time"

	"github.com/Czcan/TimeLine/utils/database"
	"gorm.io/gorm"
)

type Liker struct {
	ID        int
	UserID    int  `gorm:"index:idx_liker"`
	AccountID int  `gorm:"index:idx_liker"`
	IsLiked   bool `gorm:"null"`
	IsFollwer bool `gorm:"null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func UpdateLiker(db *gorm.DB, userID int, accountID int, status bool) (int, error) {
	var (
		account = &Account{}
		liker   = &Liker{}
		likers  = 0
	)
	if err := db.Where("id = ?", accountID).First(&account).Error; err != nil {
		return 0, err
	}
	if err := db.Where("user_id = ? AND account_id = ?", userID, accountID).First(&liker).Error; err != nil {
		tx := db.Begin()
		tx.Exec("INSERT INTO likers (user_id, account_id, is_liked) VALUES (?, ?, ?)", userID, accountID, 1)
		tx.Exec("UPDATE accounts SET likers = likers + 1 WHERE id = ?", accountID)
		if tx.Error != nil {
			tx.Rollback()
			return 0, tx.Error
		}
		tx.Commit()
		return account.Likers + 1, nil
	} else {
		if liker.IsLiked == status {
			return account.Likers, errors.New("invalid operation")
		}
		if status {
			likers = account.Likers + 1
		} else {
			likers = account.Likers - 1
		}
		db.Exec("UPDATE likers SET is_liked = ? WHERE user_id = ? AND account_id = ?", status, userID, accountID)
		db.Exec("UPDATE accounts SET likers = ? WHERE id = ?", likers, accountID)
	}
	return likers, nil
}

func UpdateFollwerAndSyncCollection(db *gorm.DB, accountID, userID int, status bool) (int, error) {
	var (
		account  = &Account{}
		liker    = &Liker{}
		follwers = 0
	)
	if err := db.Where("id = ?", accountID).First(&account).Error; err != nil {
		return 0, err
	}
	if err := db.Where("user_id = ? AND account_id = ?", userID, accountID).First(&liker).Error; err != nil {
		err := database.Transaction(db, func(tx *gorm.DB) error {
			tx.Exec("INSERT INTO likers (user_id, account_id, is_follwer) VALUES (?, ?, ?)", userID, accountID, 1)
			tx.Exec("UPDATE accounts SET follwers = follwers + 1 WHERE id = ?", accountID)
			if err := SaveCollection(tx, userID, accountID); err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return 0, err
		}
		return account.Follwers + 1, nil
	}
	if liker.IsFollwer == status {
		return account.Likers, errors.New("invalid operation")
	}
	if status {
		follwers = account.Follwers + 1
	} else {
		follwers = account.Follwers - 1
	}
	err := database.Transaction(db, func(tx *gorm.DB) error {
		tx.Exec("UPDATE likers SET is_follwer = ? WHERE user_id = ? AND account_id = ?", status, userID, accountID)
		tx.Exec("UPDATE accounts SET follwers = ? WHERE id = ?", follwers, accountID)
		if status {
			if err := SaveCollection(tx, userID, accountID); err != nil {
				return err
			}
		} else {
			if err := DeleteCollection(tx, userID, accountID); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return follwers, nil
}
