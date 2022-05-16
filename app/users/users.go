package users

import (
	"fmt"
	"net/http"

	"github.com/Czcan/TimeLine/app/entries"
	"github.com/Czcan/TimeLine/app/helpers"
	"github.com/Czcan/TimeLine/models"
	"github.com/Czcan/TimeLine/utils/jwt"
	"github.com/jinzhu/gorm"
)

type Handler struct {
	DB         *gorm.DB
	JwtClient  jwt.JWTValidate
	UploadPath string
}

func New(db *gorm.DB, jwtClient jwt.JWTValidate, path string) Handler {
	return Handler{
		DB:         db,
		JwtClient:  jwtClient,
		UploadPath: path,
	}
}

func (h Handler) Auth(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	pwd := r.FormValue("password")
	if email == "" || pwd == "" {
		helpers.RenderFailureJSON(w, 400, "email or password is empty")
		return
	}
	user, err := models.FindUser(h.DB, email, pwd)
	if err != nil {
		helpers.RenderFailureJSON(w, 400, "email or password is error")
		return
	}
	token, err := h.JwtClient.GetToken(user.Uid)
	if err != nil {
		helpers.RenderFailureJSON(w, 400, err.Error())
		return
	}
	helpers.RenderSuccessJSON(w, 200, entries.Auth{
		Token:     token,
		Email:     user.Email,
		NickName:  user.NickName,
		Gender:    user.Gender,
		Avatar:    user.GetAvatarUrl(),
		Age:       user.Age,
		Signature: user.Signature,
	})
}

func (h Handler) Register(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")
	password1 := r.FormValue("password1")
	if email == "" {
		helpers.RenderFailureJSON(w, 400, "email is empty")
		return
	}
	if password != password1 {
		helpers.RenderFailureJSON(w, 400, "incorrect password")
		return
	}
	_, err := models.FindOrCreateUser(h.DB, email, password)
	if err != nil {
		helpers.RenderFailureJSON(w, 400, "用户已存在")
		return
	}
	helpers.RenderSuccessJSON(w, 200, "注册成功")
}

func (h Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	user, err := helpers.GetCurrentUser(r, h.DB)
	if err != nil {
		helpers.RenderFailureJSON(w, 400, "invalid user")
		return
	}
	key := r.FormValue("key")
	value := r.FormValue("value")
	if key == "" || value == "" {
		helpers.RenderFailureJSON(w, 400, "invalid params")
		return
	}
	updateSQL := fmt.Sprintf("UPDATE users SET %v = ? WHERE id = ?", key)
	h.DB.Exec(updateSQL, value, user.ID)
	u := models.User{}
	h.DB.Where("id = ?", user.ID).First(&u)
	helpers.RenderSuccessJSON(w, 200, entries.Auth{
		Token:     user.Uid,
		Email:     u.Email,
		NickName:  u.NickName,
		Gender:    u.Gender,
		Avatar:    u.Avatar,
		Age:       u.Age,
		Signature: u.Signature,
	})
}

func (h Handler) Collection(w http.ResponseWriter, r *http.Request) {
	user, err := helpers.GetCurrentUser(r, h.DB)
	if err != nil {
		helpers.RenderFailureJSON(w, 400, "invalid user")
		return
	}
	selectSQL := `
		SELECT accounts.id AS id, title, content, accounts.created_at, images, likers, follwers 
		FROM accounts 
		LEFT JOIN collections ON accounts.id = collections.account_id 
		WHERE collections.user_id = ?
	`
	collections := []models.Account{}
	h.DB.Raw(selectSQL, user.ID).Scan(&collections)
	for i := 0; i < len(collections); i++ {
		collections[i].ConCatImages()
	}
	helpers.RenderSuccessJSON(w, 200, collections)
}
