package model

import "gorm.io/gorm"

type Message struct {
	gorm.Model
	ChatID uint `gorm:"column:chat_id; not null" json:"chatId"`
	UserID uint `gorm:"column:user_id; not null" json:"userId"`
}
