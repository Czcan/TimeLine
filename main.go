package main

import (
	"log"
	"net/http"
	"time"

	"github.com/Czcan/TimeLine/config"
	"github.com/Czcan/TimeLine/models"
	"github.com/Czcan/TimeLine/server"
	"github.com/Czcan/TimeLine/utils/jwt"
)

func main() {
	c := config.MustGetAppConfig()
	db := config.MustGetDB()
	db.AutoMigrate(&models.User{}, &models.Folder{}, &models.Note{}, &models.Collection{}, &models.Account{}, &models.Comment{}, &models.Liker{})
	jwtClient := jwt.New([]byte("123"), time.Hour*2, "TimeLine")
	router := server.New(db, jwtClient)

	log.Printf("localhost%s\n", c.Port)
	log.Fatal(http.ListenAndServe(c.Port, router))
}
