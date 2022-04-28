package server

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/jinzhu/gorm"
)

func New(db *gorm.DB) *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	// router.Route("/users", func(r chi.Router) {
	// 	router.Use(jwtauth.Verifier(tokenAuth))
	// 	router.Get("/v1/my")
	// })

	return router
}
