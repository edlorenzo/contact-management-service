package config

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DbConfig struct {
	Host        string
	Database    string
	User        string
	Password    string
	Port        string
	ConnTimeout int
}

func getDbConfig() DbConfig {
	return DbConfig{
		Host:        getString("DB_HOST"),
		Port:        getString("DB_PORT"),
		Database:    getString("DB_DATABASE"),
		User:        getString("DB_USERNAME"),
		Password:    getString("DB_PASSWORD"),
		ConnTimeout: getIntOrDefault("DB_CONN_TIMEOUT", 5),
	}
}

func (m DbConfig) GetDB() (*gorm.DB, error) {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?connect_timeout=%d", m.User, m.Password, m.Host, m.Port, m.Database, m.ConnTimeout)
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{},
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: newLogger})

	if err != nil {
		return nil, fmt.Errorf("db Conn: %w", err)
	}

	return db, nil
}
