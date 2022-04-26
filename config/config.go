package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type configuration struct {
	APPConfig appconfig
}

type appconfig struct {
	Port string `default:"8099"`
	DB   string
}

var Configuration = configuration{}

func init() {
	filePath := path.Join(os.Getenv("GOPATH"), "src/github.com/Czcan/Timeline/config/config.json")
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Open file error : %v\n", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&Configuration)
	if err != nil {
		fmt.Printf("Init config error: %v\n", err)
	}
}

func MustGetAPPDB() *gorm.DB {
	DB, err := gorm.Open("mysql", Configuration.APPConfig.DB)
	if err != nil {
		panic(err)
	}
	fmt.Println("Success to connect!")
	fmt.Println("DB:", Configuration.APPConfig.DB)
	fmt.Println("HOST Port:", Configuration.APPConfig.Port)
	DB.LogMode(true)
	return DB
}
