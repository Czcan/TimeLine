package server

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Czcan/TimeLine/internal/api/accounts"
	"github.com/Czcan/TimeLine/internal/api/comments"
	"github.com/Czcan/TimeLine/internal/api/folder"
	"github.com/Czcan/TimeLine/internal/api/likers"
	"github.com/Czcan/TimeLine/internal/api/notes"
	"github.com/Czcan/TimeLine/internal/api/tasks"
	"github.com/Czcan/TimeLine/internal/api/upload"
	"github.com/Czcan/TimeLine/internal/api/users"
	"github.com/Czcan/TimeLine/internal/middlewares"
	"github.com/Czcan/TimeLine/pkg/jwt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"gorm.io/gorm"
)

func New(db *gorm.DB, jwtClient jwt.JWTValidate) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(middlewares.JwtAuthentication(jwtClient))
	r.Use(middlewares.Logger(middlewares.Option{
		ServiceName: "TimeLine",
		FormattedTime: func(t time.Time) string {
			return t.In(time.FixedZone("local", 8*60*60)).Format("2006-01-02 15:04:05")
		},
	}))

	userHandler := users.New(db, jwtClient)
	noteHandler := notes.New(db)
	uploadHandler := upload.New(db)
	accountHandler := accounts.New(db)
	likerHandler := likers.New(db)
	folderHandler := folder.New(db)
	commentHandler := comments.New(db)
	taskHandler := tasks.New(db)

	//user
	r.Post("/api/auth", userHandler.Auth)
	r.Post("/api/register", userHandler.Register)
	r.Post("/api/user/update", userHandler.UpdateUser)
	r.Get("/api/user/collection", userHandler.Collection)

	//account
	r.Get("/api/account/home", accountHandler.AccountList)
	r.Post("/api/account/create", accountHandler.CreateAccount)
	r.Get("/api/account/detail/{id}", accountHandler.AccoutDetail)
	r.Delete("/api/account/deleted", accountHandler.AccoutDelted)
	r.Get("/api/account/person/list", accountHandler.AcccountPersonal)

	//upload
	r.Post("/api/upload", uploadHandler.UploadImage)

	//note
	r.Get("/api/note/list", noteHandler.List)
	r.Post("/api/note/create", noteHandler.Create)
	r.Post("/api/note/update", noteHandler.Update)
	r.Delete("/api/note/deleted", noteHandler.Deleted)

	//task
	r.Get("/api/task/list", taskHandler.TaskList)
	r.Post("/api/task/create", taskHandler.Create)
	r.Delete("/api/task/deleted", taskHandler.Delete)

	//folder
	r.Get("/api/folder/list", folderHandler.List)
	r.Post("/api/folder/create", folderHandler.Create)
	r.Delete("/api/folder/deleted", folderHandler.Deleted)

	//liker
	r.Get("/api/liker", likerHandler.Liker)
	r.Get("/api/follwer", likerHandler.Follwer)

	//comment
	r.Post("/api/comment", commentHandler.Comment)

	//staticFS
	r.Get("/upload/*", StatisFS("/upload"))

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