package config

import (
	"github.com/jinzhu/configor"
	"github.com/jinzhu/gorm"
)

type APPConfig struct {
	Port string `default:"8099"`
	DB   string
}

var _APPConfig *APPConfig

func MustGetAPPConfig() APPConfig {
	if _APPConfig != nil {
		return *_APPConfig
	}

	appconfig := &APPConfig{}
	err := configor.New(&configor.Config{ENVPrefix: "APP"}).Load(appconfig)
	if err != nil {
		panic(err)
	}

	_APPConfig = appconfig
	return *_APPConfig
}

func MustGetAPPDB() *gorm.DB {
	c := MustGetAPPConfig()
	DB, err := gorm.Open("mysql", c.DB)
	if err != nil {
		panic(err)
	}

	DB.LogMode(true)
	return DB
}
