package server

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Czcan/TimeLine/app/accounts"
	"github.com/Czcan/TimeLine/app/folder"
	"github.com/Czcan/TimeLine/app/likers"
	"github.com/Czcan/TimeLine/app/notes"
	"github.com/Czcan/TimeLine/app/upload"
	"github.com/Czcan/TimeLine/app/users"
	"github.com/Czcan/TimeLine/config"
	"github.com/Czcan/TimeLine/middlewares"
	"github.com/Czcan/TimeLine/utils/jwt"
	"github.com/Czcan/TimeLine/utils/logger"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"gorm.io/gorm"
)

func New(db *gorm.DB, jwtClient jwt.JWTValidate, c config.AppConfig) *chi.Mux {
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
	accountHandler := accounts.New(db)
	likerHandler := likers.New(db)
	folderHandler := folder.New(db)
	//user
	r.Post("/api/auth", userHandler.Auth)
	r.Post("/api/register", userHandler.Register)
	r.Post("/api/user/update", userHandler.UpdateUser)
	r.Get("/api/user/collection", userHandler.Collection)

	//account
	r.Get("/api/account/home", accountHandler.AccountList)
	r.Post("/api/account/create", accountHandler.CreateAccount)
	r.Get("/api/account/detail/{id}", accountHandler.AccoutDetail)

	//upload
	r.Post("/api/upload", uploadHandler.UploadImage)

	//note
	r.Get("/api/note/list", noteHandler.List)
	r.Post("/api/note/create", noteHandler.Create)
	r.Post("/api/note/update", noteHandler.Update)
	r.Delete("/api/note/deleted", noteHandler.Deleted)

	//folder
	r.Get("/api/folder/list", folderHandler.List)
	r.Post("/api/folder/create", folderHandler.Create)
	r.Delete("/api/folder/deleted", folderHandler.Deleted)

	//liker
	r.Get("/api/liker", likerHandler.Liker)
	r.Get("/api/follwer", likerHandler.Follwer)

	//staticFS
	r.Get("/images/*", StatisFS(c.AvatarPath))
	r.Get("/accountimg/*", StatisFS(c.AccountImgPath))

	//test
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	})

	return r
}

func StatisFS(name string) http.HandlerFunc {
	workDir, _ := os.Getwd()
	filepath := filepath.Join(workDir, name)
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := chi.RouteContext(r.Context())
		path := strings.TrimSuffix(ctx.RoutePattern(), "/*")
		fs := http.StripPrefix(path, http.FileServer(http.Dir(filepath)))
		fs.ServeHTTP(w, r)
	}
}
