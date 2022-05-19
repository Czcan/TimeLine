package helpers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/Czcan/TimeLine/app/entries"
	"github.com/Czcan/TimeLine/models"
	"github.com/Czcan/TimeLine/utils/jwt"
	"github.com/iancoleman/strcase"
	"gorm.io/gorm"
)

func GetCurrentUser(r *http.Request, db *gorm.DB) (*models.User, error) {
	claim, ok := r.Context().Value("token").(*jwt.Token)
	if !ok {
		return nil, errors.New("invalid token")
	}
	user := &models.User{}
	if err := db.Where("Uid = ?", claim.Uid).First(user).Error; err != nil {
		return nil, errors.New("invalid user")
	}
	return user, nil
}

func RenderSuccessJSON(w http.ResponseWriter, code int, data interface{}) {
	result, _ := json.Marshal(entries.Response{
		Code: code,
		Data: data,
	})
	status(w, code)
	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
}

func RenderFailureJSON(w http.ResponseWriter, code int, message interface{}) {
	result, _ := json.Marshal(entries.Response{
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

func GetParams(r *http.Request, params string) string {
	url := r.FormValue(strcase.ToCamel(params))
	if url == "" {
		url = r.FormValue(strings.ToLower(strcase.ToScreamingSnake(params)))
	}
	return url
}

func GetParamsInt(r *http.Request, params string) int {
	value := GetParams(r, params)
	i, _ := strconv.Atoi(value)
	return i
}

func GetParamsBool(r *http.Request, params string) bool {
	value := GetParams(r, params)
	if value == "1" || value == "true" {
		return true
	}
	if value == "0" || value == "false" {
		return false
	}
	return false
}
