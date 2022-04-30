package server

import (
	"net/http"

	"github.com/Czcan/TimeLine/app/users"
	"github.com/Czcan/TimeLine/app/verify"
	email "github.com/Czcan/TimeLine/libs/emailch"
	middlewares "github.com/Czcan/TimeLine/middleware"
	"github.com/Czcan/TimeLine/utils/jwt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gorilla/sessions"
	"github.com/jinzhu/gorm"
)

func New(db *gorm.DB, sessionStore *sessions.CookieStore, jwtClient jwt.JWTValidate, emailClient *email.EmailClient) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(middlewares.JwtAuthentication(jwtClient))

	userHandler := users.New(db, jwtClient, sessionStore)
	captchaHandler := verify.New(emailClient, sessionStore)

	r.Post("/api/auth", userHandler.Auth)
	r.Post("/api/register", userHandler.Register)
	r.Get("/api/email", captchaHandler.Emailcaptcha)

	r.Get("/api", func(w http.ResponseWriter, r *http.Request) {
		claim, ok := r.Context().Value("token").(*jwt.Token)
		if ok {
			w.Write([]byte(claim.Uid))
		}
	})

	return r
}
