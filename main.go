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
	"github.com/patrickmn/go-cache"
)

func main() {
	c := config.MustGetAppConfig()
	db := config.MustGetDB()
	db.AutoMigrate(&models.User{})
	defer db.Close()

	cache := cache.New(time.Minute*3, time.Minute*5)
	jwtClient := jwt.New([]byte("123"), time.Hour*2, "TimeLine")
	emailClient, err := email.New(c.EmailHost, c.EmailUser, c.EmailName, c.EmailSecret)
	if err != nil {
		panic(err)
	}
	router := server.New(db, cache, jwtClient, emailClient)

	log.Println("localhost:8099")
	log.Fatal(http.ListenAndServe(c.Port, router))
}
