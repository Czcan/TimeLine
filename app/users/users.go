package users

import (
	"net/http"

	"github.com/Czcan/TimeLine/app/entries"
	"github.com/Czcan/TimeLine/app/helpers"
	"github.com/Czcan/TimeLine/models"
	"github.com/Czcan/TimeLine/utils/jwt"
	"github.com/jinzhu/gorm"
	"github.com/patrickmn/go-cache"
)

type Handler struct {
	DB        *gorm.DB
	JwtClient jwt.JWTValidate
	Cache     *cache.Cache
}

func New(db *gorm.DB, jwtClient jwt.JWTValidate, cache *cache.Cache) Handler {
	return Handler{
		DB:        db,
		JwtClient: jwtClient,
		Cache:     cache,
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
	val, ok := h.Cache.Get(email)
	if !ok {
		helpers.RenderFailureJSON(w, 500, "验证码过期")
		return
	}
	captcha, ok := val.(string)
	if !ok {
		helpers.RenderFailureJSON(w, 500, "注册失败")
		return
	}
	if captcha != code {
		helpers.RenderFailureJSON(w, 400, "验证码错误")
		return
	}
	_, err := models.FindOrCreateUser(h.DB, email, pwd)
	if err != nil {
		helpers.RenderFailureJSON(w, 400, "用户已存在")
		return
	}
	defer h.Cache.Delete(email)
	helpers.RenderSuccessJSON(w, 200, "注册成功")
}
