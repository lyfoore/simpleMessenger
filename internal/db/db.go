package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"simpleTodoList/internal/model"
)

func InitDB(dsn string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	log.Println("Connected to database")

	err = db.AutoMigrate(&model.User{}, &model.Chat{}, &model.Message{}, &model.ChatParticipants{})
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}
	return db
}
