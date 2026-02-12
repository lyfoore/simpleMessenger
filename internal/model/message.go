package model

import "gorm.io/gorm"

type Message struct {
	gorm.Model
	Text   string `gorm:"column:text; not null" json:"text"`
	ChatID uint   `gorm:"column:chat_id; not null" json:"chatId"`
	UserID uint   `gorm:"column:user_id; not null" json:"userId"`
}
