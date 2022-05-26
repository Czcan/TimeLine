package entries

type Account struct {
	ID         int      `json:"id"`
	Title      string   `json:"title"`
	Content    string   `json:"content"`
	Likers     int      `json:"likers"`
	Follwers   int      `json:"follwers"`
	CreatedAt  int      `json:"created_at"`
	ImageSlice []string `json:"images"`
}

type Comment struct {
	NickName  string `json:"nick_name"`
	Content   string `json:"content"`
	AvatarUrl string `json:"avatar_url"`
	Date      int    `json:"date"`
}

type AccountDetail struct {
	Account      *Account      `json:"account"`
	Comments     []Comment     `json:"comments"`
	LikerFollwer *LikerFollwer `json:"liker_follwer"`
}
