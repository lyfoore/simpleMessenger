package main

import (
	"fmt"
	"log"
	"os"
	"simpleTodoList/internal/db"
	"simpleTodoList/internal/model"
	"simpleTodoList/internal/router"
)

func main() {
	dsn := createDSN()

	database := db.InitDB(dsn)

	//userRepo := postgres.NewUserRepository(database)
	//chatRepo := postgres.NewChatRepository(database)
	//messageRepo := postgres.NewMessageRepository(database)
	//chatParticipantsRepo := postgres.NewChatParticipantsRepository(database)

	r := router.NewRouter()
	r.SetupRouter()
	r.Run()
}

func getEnv(key, defaultValue string) string {
	result := os.Getenv(key)
	if result != "" {
		return result
	} else {
		return defaultValue
	}
}

func createDSN() string {
	host := getEnv("DB_HOST", "localhost")
	user := getEnv("DB_USER", "postgres")
	pass := getEnv("DB_PASS", "postgres")
	port := getEnv("DB_PORT", "5432")
	dbname := getEnv("DB_NAME", "database")
	sslmode := getEnv("DB_SSLMODE", "disable")

	dsn := fmt.Sprintf("host=%s user=%s password=%s port=%s dbname=%s sslmode=%s", host, user, pass, port, dbname, sslmode)
	return dsn
}
