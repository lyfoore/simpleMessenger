package model

import "gorm.io/gorm"

type Chat struct {
	gorm.Model
	Name string `gorm:"column:name; not null" json:"name"`
}
