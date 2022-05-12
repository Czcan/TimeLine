package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/jinzhu/configor"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type AppConfig struct {
	DB          string
	Port        string `default:":9091"`
	SecretKey   string
	EmailUser   string
	EmailSecret string
	EmailHost   string
	EmailName   string
	AvatarPath  string
}

var (
	_AppConfig *AppConfig
	once       sync.Once
)

func MustGetAppConfig() AppConfig {
	root := inferRootDir()
	if _AppConfig == nil {
		once.Do(
			func() {
				appConfig := &AppConfig{}
				err := configor.Load(appConfig, root+"/config.yml")
				if err != nil {
					panic(err)
				}
				_AppConfig = appConfig
			})
	}
	return *_AppConfig
}

func MustGetDB() *gorm.DB {
	c := MustGetAppConfig()
	DB, err := gorm.Open("mysql", c.DB)
	if err != nil {
		panic(err)
	}
	fmt.Println("Success to connect!")
	fmt.Println("DB:", c.DB)
	fmt.Println("HOST Port:", c.Port)
	DB.LogMode(true)
	return DB
}

func inferRootDir() string {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	var infer func(d string) string
	infer = func(d string) string {
		if exists(d + "/config.yml") {
			return d
		}
		return infer(filepath.Dir(d))
	}

	return infer(cwd)
}
func exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}
