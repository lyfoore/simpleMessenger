package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Login string `gorm:"column:login; not null; unique" json:"login"`
	Name  string `gorm:"column:name; not null" json:"name"`
}
