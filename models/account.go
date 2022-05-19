package models

import (
	"fmt"
	"strings"

	"github.com/Czcan/TimeLine/app/entries"
	"gorm.io/gorm"
)

type Account struct {
	ID         int            `json:"id"`
	UserID     int            `json:"user_id"`
	Images     string         `json:"-"`
	Title      string         `json:"title"`
	Content    string         `json:"content"`
	Likers     int            `json:"likers"`
	Follwers   int            `json:"follwers"`
	CreatedAt  int            `json:"created_at"`
	UpdatedAt  int            `json:"-"`
	DeletedAt  gorm.DeletedAt `json:"-"`
	ImageSlice []string       `json:"images" gorm:"-"`
}

func (a *Account) ConCatImages() {
	for _, imageID := range strings.Split(a.Images, ",") {
		imageID = strings.TrimSpace(imageID)
		image := fmt.Sprintf("/accountimg/%d/%s.jpg", a.ID, imageID)
		a.ImageSlice = append(a.ImageSlice, image)
	}
}

func FindAccountDetail(db *gorm.DB, accountID int) (*entries.Account, []entries.Comment, error) {
	var (
		account  = &Account{}
		comments = []entries.Comment{}
	)
	if err := db.Where("id = ?", accountID).First(&account).Error; err != nil {
		return nil, nil, err
	}
	account.ConCatImages()
	selectSQL := `
		SELECT comments.content, comments.created_at AS date, users.nick_name, CONCAT('/images/', users.id, '.jpg') AS avatar_url
		FROM comments
		LEFT JOIN users ON comments.user_id = users.id
		WHERE comments.account_id = ?
	`
	db.Raw(selectSQL, accountID).Scan(&comments)
	entryAccount := &entries.Account{
		Title:      account.Title,
		Content:    account.Content,
		ID:         account.ID,
		Likers:     account.Likers,
		Follwers:   account.Likers,
		CreatedAt:  account.CreatedAt,
		ImageSlice: account.ImageSlice,
	}
	return entryAccount, comments, nil
}
