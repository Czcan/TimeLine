package users

import (
	"net/http"

	"github.com/Czcan/TimeLine/app/entries"
	"github.com/Czcan/TimeLine/app/helpers"
	"github.com/Czcan/TimeLine/models"
	"github.com/Czcan/TimeLine/utils/errcode"
	"github.com/Czcan/TimeLine/utils/jwt"
	"github.com/Czcan/TimeLine/utils/validate"
	"gorm.io/gorm"
)

type Handler struct {
	DB        *gorm.DB
	JWTClient jwt.JWTValidate
}

func New(db *gorm.DB, client jwt.JWTValidate) Handler {
	return Handler{
		DB:        db,
		JWTClient: client,
	}
}

func (h Handler) Auth(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	pwd := r.FormValue("password")

	if !validate.ValidateEmail(email) || !validate.ValidateStringEmpty(pwd) {
		helpers.RenderFailureJSON(w, 400, "email or password is empty")
		return
	}

	user, err := models.FindUser(h.DB, email, pwd)
	if err != nil {
		helpers.RenderFailureJSON(w, 400, "email or password is error")
		return
	}

	token, err := h.JWTClient.GetToken(user.Uid)
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
	if !validate.ValidateEmail(email) {
		helpers.RenderFailureJSON(w, 400, "email is empty")
		return
	}
	if password != password1 || !validate.ValidateStringEmpty(password, password1) {
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
		helpers.RenderFailureJSON(w, 400, errcode.GetMsg(errcode.ERROR_TOKEN))
		return
	}
	key := r.FormValue("key")
	value := r.FormValue("value")
	if !validate.ValidateStringEmpty(key, value) {
		helpers.RenderFailureJSON(w, 400, errcode.GetMsg(errcode.ERROR_PARAMS))
		return
	}
	user = models.UpdateAndFindUser(h.DB, user.ID, key, value)
	helpers.RenderSuccessJSON(w, 200, entries.Auth{
		Token:     user.Uid,
		Email:     user.Email,
		NickName:  user.NickName,
		Gender:    user.Gender,
		Avatar:    user.Avatar,
		Age:       user.Age,
		Signature: user.Signature,
	})
}

func (h Handler) Collection(w http.ResponseWriter, r *http.Request) {
	user, err := helpers.GetCurrentUser(r, h.DB)
	if err != nil {
		helpers.RenderFailureJSON(w, 400, errcode.GetMsg(errcode.ERROR_TOKEN))
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
