package verify

import (
	"net/http"
	"strings"
	"time"

	"github.com/Czcan/TimeLine/app/helpers"
	email "github.com/Czcan/TimeLine/libs/emailch"
	"github.com/Czcan/TimeLine/utils"
	"github.com/patrickmn/go-cache"
)

type Handler struct {
	Cache       *cache.Cache
	EmailClient *email.EmailClient
}

func New(client *email.EmailClient, c *cache.Cache) Handler {
	return Handler{
		EmailClient: client,
		Cache:       c,
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
		helpers.RenderFailureJSON(w, 500, err.Error())
		return
	}
	h.Cache.Set(email, code, time.Minute*3)
	helpers.RenderSuccessJSON(w, 200, body)
}
