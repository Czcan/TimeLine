package users

import (
	"net/http"
	"time"

	"github.com/Czcan/TimeLine/app/entries"
	"github.com/Czcan/TimeLine/app/helpers"
	"github.com/Czcan/TimeLine/app/verify"
	"github.com/Czcan/TimeLine/models"
	"github.com/Czcan/TimeLine/utils/jwt"
	"github.com/gorilla/sessions"
	"github.com/jinzhu/gorm"
)

type Handler struct {
	DB           *gorm.DB
	JwtClient    jwt.JWTValidate
	SessionStore *sessions.CookieStore
}

func New(db *gorm.DB, jwtClient jwt.JWTValidate, sessionStore *sessions.CookieStore) Handler {
	return Handler{
		DB:           db,
		JwtClient:    jwtClient,
		SessionStore: sessionStore,
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
		Token:    token,
		Email:    user.Email,
		NickName: user.NickName,
		Gender:   user.Gender,
		Avatar:   user.Avatar,
		Age:      user.Age,
	})
}

func (h Handler) Register(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	pwd := r.FormValue("password")
	code := r.FormValue("code")
	if email == "" || pwd == "" {
		helpers.RenderFailureJSON(w, 400, "email or password is empty")
		return
	}
	sessions, err := h.SessionStore.Get(r, "timeLine")
	if err != nil {
		helpers.RenderFailureJSON(w, 500, "注册失败")
		return
	}
	captcha, ok := sessions.Values[email].(verify.Captcha)
	if !ok {
		helpers.RenderFailureJSON(w, 500, "注册失败")
		return
	}
	if !time.Now().Before(captcha.ExpiresAt) {
		helpers.RenderFailureJSON(w, 400, "验证码过期")
		return
	}
	if captcha.Code != code {
		helpers.RenderFailureJSON(w, 400, "验证码错误")
		return
	}
	_, err = models.FindOrCreateUser(h.DB, email, pwd)
	if err != nil {
		helpers.RenderFailureJSON(w, 400, err.Error())
		return
	}
	helpers.RenderSuccessJSON(w, 200, "注册成功")
}
