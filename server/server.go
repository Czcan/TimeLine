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
	"github.com/jinzhu/gorm"
	"github.com/patrickmn/go-cache"
)

func New(db *gorm.DB, cache *cache.Cache, jwtClient jwt.JWTValidate, emailClient *email.EmailClient) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(middlewares.JwtAuthentication(jwtClient))

	userHandler := users.New(db, jwtClient, cache)
	captchaHandler := verify.New(emailClient, cache)

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
