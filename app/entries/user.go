package entries

type User struct {
	ID       int    `json:"id"`
	Avatar   string `json:"avatar"`
	NickName string `json:"nick_name"`
	Uid      string `json:"Uid"`
	Gender   bool   `json:"gender"`
	Age      int    `json:"age"`
	Email    string `json:"email"`
}
