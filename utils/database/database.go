package database

import "github.com/jinzhu/gorm"

func Transaction(db *gorm.DB, callback func(db *gorm.DB) error) error {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return err
	}
	if err := callback(tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}
