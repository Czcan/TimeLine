package main

import (
	"log"
	"net/http"
	"time"

	"github.com/Czcan/TimeLine/config"
	email "github.com/Czcan/TimeLine/libs/emailch"
	"github.com/Czcan/TimeLine/models"
	"github.com/Czcan/TimeLine/server"
	"github.com/Czcan/TimeLine/utils/jwt"
	"github.com/gorilla/sessions"
)

func main() {
	c := config.MustGetAppConfig()
	db := config.MustGetDB()
	db.AutoMigrate(&models.User{})
	defer db.Close()

	sessionStore := sessions.NewCookieStore([]byte("TimLine"))
	jwtClient := jwt.New([]byte("123"), time.Hour*2, "TimeLine")
	emailClient, err := email.New(c.EmailHost, c.EmailUser, c.EmailName, c.EmailSecret)
	if err != nil {
		panic(err)
	}
	router := server.New(db, sessionStore, jwtClient, emailClient)

	log.Println("localhost:8099")
	log.Fatal(http.ListenAndServe(c.Port, router))
}
