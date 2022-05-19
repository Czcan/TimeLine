package config

import (
	"sync"

	"github.com/Czcan/TimeLine/utils"
	"github.com/jinzhu/configor"
)

type AppConfig struct {
	DB             string
	Port           string `default:":9091"`
	SecretKey      string
	EmailUser      string
	EmailSecret    string
	EmailHost      string
	EmailName      string
	AvatarPath     string
	AccountImgPath string
}

var (
	_AppConfig *AppConfig
	once       sync.Once
)

func MustGetAppConfig() AppConfig {
	root := utils.InferRootDir()
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
