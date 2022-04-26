package main

import (
	"log"
	"net/http"
	"timeline/config"
	"timeline/models"
	"timeline/server"
)

func main() {
	// appconfig := config.Configuration.APPConfig
	db := config.MustGetAPPDB()
	db.AutoMigrate(&models.User{})
	defer db.Close()

	router := server.New(db)

	log.Fatal(http.ListenAndServe(config.Configuration.APPConfig.Port, router))
}
