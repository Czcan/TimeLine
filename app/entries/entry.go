package entries

type Success struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}

type Error struct {
	Code    int         `json:"code"`
	Message interface{} `json:"message"`
}
