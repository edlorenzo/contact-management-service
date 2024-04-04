package config

import (
	"fmt"
	"time"
)

type AppConfig struct {
	Name         string
	Addr         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	Environment  string
}

func getAppConfig() AppConfig {
	addr := fmt.Sprintf("%s:%v", getString("APP_HOST"), getInt("APP_PORT"))
	return AppConfig{
		Addr:         addr,
		Name:         getString("APP_NAME"),
		Environment:  getString("APP_ENVIRONMENT"),
		ReadTimeout:  time.Duration(getInt("APP_READ_TIMEOUT")) * time.Second, // time.Duration(getInt("APP_READ_TIMEOUT")),
		WriteTimeout: time.Duration(getInt("APP_WRITE_TIMEOUT")) * time.Second,
	}
}
