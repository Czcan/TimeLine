package helpers

import (
	"encoding/json"
	"net/http"

	"github.com/Czcan/TimeLine/app/entries"
)

// func GetCurrentUser(db *gorm.DB, jwtToken string) (*models.User, error) {
// 	// 验证 jwt_token

// 	// 验证成功，通过uid 查询用户信息并返回
// }

func RenderSuccessJSON(w http.ResponseWriter, code int, data interface{}) {
	result, _ := json.Marshal(entries.Success{
		Code: code,
		Data: data,
	})
	status(w, code)
	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
}

func RenderFailureJSON(w http.ResponseWriter, code int, message interface{}) {
	result, _ := json.Marshal(entries.Error{
		Code:    code,
		Message: message,
	})
	status(w, code)
	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
}

func status(w http.ResponseWriter, code int) {
	w.WriteHeader(code)
}
