package users

import (
	"net/http"

	"github.com/Czcan/TimeLine/internal/api/helpers"
	"github.com/Czcan/TimeLine/internal/models"
	"github.com/Czcan/TimeLine/pkg/errcode"
	"github.com/Czcan/TimeLine/pkg/jwt"
	"github.com/Czcan/TimeLine/pkg/validate"
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
	user, uid, err := models.FindUser(h.DB, email, pwd)
	if err != nil {
		helpers.RenderFailureJSON(w, 400, "email or password is error")
		return
	}
	token, err := h.JWTClient.GetToken(uid)
	if err != nil {
		helpers.RenderFailureJSON(w, 400, err.Error())
		return
	}
	user.Token = token
	helpers.RenderSuccessJSON(w, 200, user)
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
	err := models.FindOrCreateUser(h.DB, email, password)
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
	auth := models.UpdateAndFindUser(h.DB, user.ID, key, value)
	helpers.RenderSuccessJSON(w, 200, auth)
}

func (h Handler) Collection(w http.ResponseWriter, r *http.Request) {
	user, err := helpers.GetCurrentUser(r, h.DB)
	if err != nil {
		helpers.RenderFailureJSON(w, 400, errcode.GetMsg(errcode.ERROR_TOKEN))
		return
	}
	accounts := models.FindCollection(h.DB, user.ID)
	helpers.RenderSuccessJSON(w, 200, accounts)
}

func (h Handler) DeleteCollection(w http.ResponseWriter, r *http.Request) {

}