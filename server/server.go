package server

import (
	"net/http"
	"time"

	"github.com/Czcan/TimeLine/app/notes"
	"github.com/Czcan/TimeLine/app/users"
	email "github.com/Czcan/TimeLine/libs/emailch"
	middlewares "github.com/Czcan/TimeLine/middleware"
	"github.com/Czcan/TimeLine/utils/jwt"
	"github.com/Czcan/TimeLine/utils/logger"
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
	r.Use(logger.Logger(logger.Option{
		ServiceName: "TimeLine",
		FormattedTime: func(t time.Time) string {
			return t.In(time.FixedZone("local", 8*60*60)).Format("2006-01-02 15:04:05")
		},
	}))
	userHandler := users.New(db, jwtClient, cache)
	noteHandler := notes.New(db)
	
	r.Post("/api/auth", userHandler.Auth)
	r.Post("/api/register", userHandler.Register)

	r.Get("/api/note/list", noteHandler.NoteList)
	r.Post("/api/note/create", noteHandler.CreateNote)
	r.Post("/api/note/update", noteHandler.FinishNote)
	r.Get("/api/folder/list", noteHandler.FolderList)
	r.Post("/api/folder/create", noteHandler.CreateFolder)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	})

	return r
}
