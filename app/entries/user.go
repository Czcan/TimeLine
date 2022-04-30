package entries

type Auth struct {
	Token    string `json:"token"`
	Email    string `json:"email"`
	NickName string `json:"nick_name"`
	Avatar   string `json:"avatar"`
	Gender   int    `json:"gender"`
	Age      int    `json:"age"`
}
