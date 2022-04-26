package server

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
	"github.com/jinzhu/gorm"
)

var tokenAuth *jwtauth.JWTAuth

func New(db *gorm.DB) *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	// router.Route("/admin", func(r chi.Router) {
	// 	router.Use(jwtauth.Verifier(tokenAuth))
	// 	router.Get("/users")
	// })
	// router.Use(jwtauth.Verifier(tokenAuth))
	return router
}
