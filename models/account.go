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
	NickName   string         `json:"nick_name"`
	AvatarUrl  string         `json:"avatar_url"`
}

//hook-find
func (a *Account) AfterFind(tx *gorm.DB) (err error) {
	for _, imageID := range strings.Split(a.Images, ",") {
		imageID = strings.TrimSpace(imageID)
		image := fmt.Sprintf("/upload/accountimg/%d/%s.jpg", a.ID, imageID)
		a.ImageSlice = append(a.ImageSlice, image)
	}
	return nil
}

func AcoountList(db *gorm.DB) []Account {
	accounts := []Account{}
	db.Model(&Account{}).
		Select("accounts.id, title, content, accounts.created_at, likers, follwers, images, users.nick_name, CONCAT('upload/avatar/images/', users.id, '.jpg') AS avatar_url").
		Joins("LEFT JOIN users ON accounts.user_id = users.id").
		Order("likers desc, follwers desc").
		Find(&accounts)
	return accounts
}

func FindAccountDetail(db *gorm.DB, accountID int, userID int) (*entries.Account, []entries.Comment, *entries.LikerFollwer, error) {
	var (
		account      = &Account{}
		LikerFollwer = &entries.LikerFollwer{}
	)
	if err := db.Model(&Account{}).
		Select("accounts.id, title, content, accounts.created_at, likers, follwers, images, users.nick_name, CONCAT('upload/avatar/images/', users.id, '.jpg') AS avatar_url").
		Joins("LEFT JOIN users ON accounts.user_id = users.id").
		Where("accounts.id = ?", accountID).
		First(&account).Error; err != nil {
		return nil, nil, nil, err
	}
	comments := FindComments(db, accountID)
	entryAccount := &entries.Account{
		Title:      account.Title,
		Content:    account.Content,
		ID:         account.ID,
		Likers:     account.Likers,
		Follwers:   account.Likers,
		CreatedAt:  account.CreatedAt,
		ImageSlice: account.ImageSlice,
		NickName:   account.NickName,
		AvatarUrl:  account.AvatarUrl,
	}
	if userID > 0 {
		db.Raw("SELECT is_liked, is_follwer FROM likers WHERE user_id = ? AND account_id = ?", userID, accountID).Scan(&LikerFollwer)
	}
	return entryAccount, comments, LikerFollwer, nil
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

func PersonAccount(db *gorm.DB, id int) []Account {
	accounts := []Account{}
	db.Where("user_id = ?", id).Find(&accounts)
	return accounts
}

func DeleteAccount(db *gorm.DB, id, userID int) ([]Account, error) {
	if err := db.Where("id = ?", id).Delete(&Account{}).Error; err != nil {
		return nil, err
	}
	return PersonAccount(db, userID), nil
}
