package modules

type User struct {
	ID       int    `json:"id"`
	Avatar   string `json:"avatar"`
	NickName string `json:"nick_name"`
	Token    string `json:"token"`
	Gender   bool   `json:"gender"`
	Age      int    `json:"age"`
}
