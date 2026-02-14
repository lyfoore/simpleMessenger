package model

import (
	"gorm.io/gorm"
	"time"
)

type Chat struct {
	gorm.Model
	Name          string `gorm:"column:name; not null" json:"name"`
	LastMessageAt time.Time
}
