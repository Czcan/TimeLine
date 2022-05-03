package models

type Account struct {
	ID         int
	UserID     int
	NotePadsID int
	Title      string
	Content    string `gorm:"type:text;size:65532"`
	IsPrivate  bool
	Images     string
	Favorite   int
	BgRgb      string
	TextSize   int
	TextRgb    string
}
