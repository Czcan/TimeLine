package models

import (
	"fmt"
	"strings"
)

type Account struct {
	ID         int      `json:"id"`
	UserID     int      `json:"user_id"`
	Images     string   `json:"-"`
	Title      string   `json:"title"`
	Content    string   `json:"content"`
	Likers     int      `json:"likers"`
	Follwers   int      `json:"follwers"`
	CreatedAt  int      `json:"created_at"`
	UpdatedAt  int      `json:"-"`
	ImageSlice []string `json:"images" gorm:"-"`
}

func (a *Account) ConCatImages() {
	for _, imageID := range strings.Split(a.Images, ",") {
		imageID = strings.TrimSpace(imageID)
		image := fmt.Sprintf("/accountimg/%d/%s.jpg", a.ID, imageID)
		a.ImageSlice = append(a.ImageSlice, image)
	}
}
