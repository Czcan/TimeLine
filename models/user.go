package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/Czcan/TimeLine/utils"
	"github.com/jinzhu/gorm"
)

type User struct {
	ID        int    `gorm:"primary_key"`
	Email     string `gorm:"unique"`
	Password  string
	Uid       string `gorm:"unique"`
	Avatar    string
	NickName  string
	Gender    int
	Age       int
	CreatedAt time.Time
	UpdatedAt time.Time
}

func FindUser(db *gorm.DB, email string, pwd string) (*User, error) {
	user := &User{}
	if err := db.Where("email = ? AND password = ?", email, pwd).First(&user).Error; err != nil {
		return nil, err
	}
	if user.NickName == "" {
		user.NickName = fmt.Sprintf("用户%s", utils.GenerateNumber(11, "0123456789"))
	}
	db.Save(&user)
	return user, nil
}

func FindOrCreateUser(db *gorm.DB, email string, pwd string) (*User, error) {
	user := &User{}
	if err := db.Where("email = ? AND password = ?", email, pwd).First(&user).Error; err == nil {
		return nil, errors.New("账号已存在")
	}
	name := fmt.Sprintf("用户%s", utils.GenerateNumber(11, "0123456789"))
	user = &User{
		Email:    email,
		Password: pwd,
		Uid:      utils.GenerateToken(email),
		NickName: name,
	}
	if err := db.Save(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}
