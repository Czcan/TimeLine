package config

import (
	"path/filepath"
	"sync"

	"github.com/Czcan/TimeLine/pkg/utils"
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
	root := utils.InferRootDir("config/config.yml")
	if _AppConfig == nil {
		once.Do(
			func() {
				appConfig := &AppConfig{}
				err := configor.Load(appConfig, filepath.Join(root, "config/config.yml"))
				if err != nil {
					panic(err)
				}
				_AppConfig = appConfig
			})
	}
	return *_AppConfig
}
