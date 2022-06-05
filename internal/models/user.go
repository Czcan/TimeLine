package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/Czcan/TimeLine/internal/api/entries"
	"github.com/Czcan/TimeLine/pkg/utils"
	"github.com/iancoleman/strcase"
	"gorm.io/gorm"
)

type User struct {
	ID        int    `gorm:"primary_key"`
	Email     string `gorm:"unique"`
	Password  string
	Uid       string `gorm:"unique"`
	Signature string
	Avatar    string
	NickName  string
	Gender    int
	Age       int
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (u *User) GetAvatarUrl() string {
	return fmt.Sprintf("upload/avatar/images/%d.jpg", u.ID)
}

func FindUser(db *gorm.DB, email string, pwd string) (*entries.Auth, string, error) {
	user := &User{}
	if err := db.Where("email = ? AND password = ?", email, pwd).First(&user).Error; err != nil {
		return nil, "", err
	}
	if user.NickName == "" {
		user.NickName = fmt.Sprintf("用户%s", utils.GenerateNumber(11, "0123456789"))
	}
	db.Save(&user)
	auth := NewEntryAuth(user)
	return auth, user.Uid, nil
}

func FindOrCreateUser(db *gorm.DB, email string, pwd string) error {
	user := &User{}
	if err := db.Where("email = ? AND password = ?", email, pwd).First(&user).Error; err == nil {
		return errors.New("账号已存在")
	}
	name := fmt.Sprintf("用户%s", utils.GenerateNumber(11, "0123456789"))
	user = &User{
		Email:    email,
		Password: pwd,
		Uid:      utils.GenerateToken(email),
		NickName: name,
	}
	if err := db.Save(&user).Error; err != nil {
		return err
	}
	return nil
}

func UpdateAndFindUser(db *gorm.DB, id int, key, value string) *entries.Auth {
	updateSQL := fmt.Sprintf("UPDATE users SET %v = ? WHERE id = ?", strcase.ToSnake(key))
	db.Exec(updateSQL, value, id)
	user := &User{}
	db.Where("id = ?", id).First(&user)
	auth := NewEntryAuth(user)
	return auth
}

func NewEntryAuth(user *User) *entries.Auth {
	return &entries.Auth{
		Email:     user.Email,
		NickName:  user.NickName,
		Gender:    user.Gender,
		Avatar:    user.GetAvatarUrl(),
		Age:       user.Age,
		Signature: user.Signature,
	}
}
