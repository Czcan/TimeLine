package models

import (
	"time"

	"github.com/Czcan/TimeLine/app/entries"
	"gorm.io/gorm"
)

type Comment struct {
	ID        int
	AccountID int
	UserID    int
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func SaveComment(db *gorm.DB, userID int, accountID int, content string) error {
	comment := &Comment{
		AccountID: accountID,
		UserID:    userID,
		Content:   content,
	}
	if err := db.Save(&comment).Error; err != nil {
		return err
	}
	return nil
}

func FindComments(db *gorm.DB, accountID int) []entries.Comment {
	comments := []entries.Comment{}
	selectSQL := `
		SELECT comments.content, comments.created_at AS date, users.nick_name, CONCAT('upload/avatar/images/', users.id, '.jpg') AS avatar_url
		FROM comments
		LEFT JOIN users ON comments.user_id = users.id
		WHERE comments.account_id = ?
	`
	db.Raw(selectSQL, accountID).Scan(&comments)
	return comments
}
