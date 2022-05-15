package entries

import (
	"time"

	"github.com/Czcan/TimeLine/models"
)

type Account struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type Comment struct {
	NickName  string `json:"nick_name"`
	Content   string `json:"content"`
	AvatarUrl string `json:"avatar_url"`
	Date      int    `json:"date"`
}

type AccountDetail struct {
	Account  models.Account `json:"account"`
	Comments []Comment      `json:"comments"`
}
