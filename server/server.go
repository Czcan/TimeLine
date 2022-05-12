package server

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Czcan/TimeLine/app/notes"
	"github.com/Czcan/TimeLine/app/upload"
	"github.com/Czcan/TimeLine/app/users"
	"github.com/Czcan/TimeLine/config"
	middlewares "github.com/Czcan/TimeLine/middleware"
	"github.com/Czcan/TimeLine/utils/jsonwt"
	"github.com/Czcan/TimeLine/utils/logger"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/jinzhu/gorm"
)

func New(db *gorm.DB, jwtClient jsonwt.JWTValidate, c config.AppConfig) *chi.Mux {
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
	userHandler := users.New(db, jwtClient, c.AvatarPath)
	noteHandler := notes.New(db)
	uploadHandler := upload.New(db, c.AvatarPath)

	r.Post("/api/auth", userHandler.Auth)
	r.Post("/api/register", userHandler.Register)
	r.Post("/api/user/update", userHandler.UpdateUser)

	r.Post("/api/upload", uploadHandler.UploadImage)

	r.Get("/api/note/list", noteHandler.NoteList)
	r.Post("/api/note/create", noteHandler.CreateNote)
	r.Post("/api/note/update", noteHandler.FinishNote)
	r.Get("/api/folder/list", noteHandler.FolderList)
	r.Post("/api/folder/create", noteHandler.CreateFolder)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	})
	workDir, _ := os.Getwd()
	fmt.Println(workDir)
	r.Get("/images/*", StatisFS(filepath.Join(workDir, c.AvatarPath)))

	return r
}

func StatisFS(filepath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := chi.RouteContext(r.Context())
		path := strings.TrimSuffix(ctx.RoutePattern(), "/*")
		fs := http.StripPrefix(path, http.FileServer(http.Dir(filepath)))
		fs.ServeHTTP(w, r)
	}
}
