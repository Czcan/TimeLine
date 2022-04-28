package main

import (
	"fmt"
	"log"
	"net/http"
	"timeline/config"
	"timeline/models"
	"timeline/server"
)

func main() {
	c := config.MustGetAppConfig()
	fmt.Println(c)
	db := config.MustGetDB()
	db.AutoMigrate(&models.User{})
	defer db.Close()

	router := server.New(db)
	log.Fatal(http.ListenAndServe(c.Port, router))
}
