package verify

import (
	"net/http"
	"strings"
	"time"

	"github.com/Czcan/TimeLine/app/helpers"
	email "github.com/Czcan/TimeLine/libs/emailch"
	"github.com/Czcan/TimeLine/utils"
	"github.com/gorilla/sessions"
)

type Captcha struct {
	Code      string
	ExpiresAt time.Time
}

type Handler struct {
	SessionStore *sessions.CookieStore
	EmailClient  *email.EmailClient
}

func New(client *email.EmailClient, sessionStore *sessions.CookieStore) Handler {
	return Handler{
		EmailClient:  client,
		SessionStore: sessionStore,
	}
}

func (h Handler) Emailcaptcha(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	code := utils.GenerateNumber(6, "1234567890")
	body := strings.Replace(`
		<html>
		<body>
			<h3>TimeLine时间轴手账</h3>
			<p>$，请在三分钟内输入验证码。为了你的账号安全，请勿将验证码告知他人</p>
		</body>
		</html>
	`, "$", code, -1)
	err := h.EmailClient.Send(email, body, "[TimeLine]时间轴手账验证码", 1)
	if err != nil {
		helpers.RenderFailureJSON(w, 500, "验证码发送失败")
		return
	}
	captcha := Captcha{Code: code, ExpiresAt: time.Now().Add(time.Minute * 3)}
	session, err := h.SessionStore.Get(r, "timeLine")
	if err != nil {
		helpers.RenderFailureJSON(w, 500, "验证码发送失败")
		return
	}
	session.Values[email] = captcha
	session.Save(r, w)
	helpers.RenderSuccessJSON(w, 200, "发送成功")
}
