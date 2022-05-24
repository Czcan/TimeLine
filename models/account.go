package models

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Czcan/TimeLine/app/entries"
	"github.com/Czcan/TimeLine/config"
	"github.com/Czcan/TimeLine/utils/database"
	"github.com/Czcan/TimeLine/utils/file"
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
		image := fmt.Sprintf("/upload/accountimg/%d/%s.jpg", a.ID, imageID)
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

func CreatedAccount(db *gorm.DB, userID int, content, title string, files []*multipart.FileHeader) error {
	var (
		count  = 0
		images = []string{}
	)
	for i := range files {
		images = append(images, strconv.Itoa(i+1))
	}
	account := &Account{
		UserID:  userID,
		Content: content,
		Title:   title,
		Images:  strings.Join(images, ","),
	}
	err := database.Transaction(db, func(tx *gorm.DB) error {
		if err := tx.Create(&account).Error; err != nil {
			return err
		}
		if err := tx.Raw("SELECT LAST_INSERT_ID() AS count").Scan(&count).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	c := config.MustGetAppConfig()
	path := filepath.Join(c.AccountImgPath, strconv.Itoa(count))
	for i, f := range files {
		_, err := file.SaveUploadFile(f, path, strconv.Itoa(i+1))
		if err != nil {
			return err
		}
		images = append(images, strconv.Itoa(i+1))
	}
	return nil
}
