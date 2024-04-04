package model

import (
	"time"
)

type (
	Contacts struct {
		ID         uint      `gorm:"primary_key"`
		Timestamp  time.Time `gorm:"index"`
		Name       string
		Email      string
		Phone      string
		ExternalID int
	}
)
