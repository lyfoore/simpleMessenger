package main

import (
	"fmt"
	"os"
	"simpleMessenger/internal/db"
	"simpleMessenger/internal/repository/postgres"
	"simpleMessenger/internal/service"
	"simpleMessenger/internal/transport/http"
	"simpleMessenger/internal/transport/websocket"
)

func main() {
	dsn := createDSN()

	database := db.InitDB(dsn)

	userRepo := postgres.NewUserRepository(database)
	chatRepo := postgres.NewChatRepository(database)
	messageRepo := postgres.NewMessageRepository(database)
	chatParticipantsRepo := postgres.NewChatParticipantsRepository(database)

	secret := getEnv("JWT_SECRET_KEY", "")

	tokenService := service.NewJwtService(secret)
	authService := service.NewAuthService(userRepo, tokenService)
	chatService := service.NewChatService(chatRepo, chatParticipantsRepo, messageRepo)
	messageService := service.NewMessageService(messageRepo, chatRepo, chatParticipantsRepo)

	authHandler := http.NewAuthHandler(authService)
	chatHandler := http.NewChatHandler(chatService)
	messageHandler := http.NewMessageHandler(messageService)

	wsHub := websocket.NewHub(messageService, chatService)
	go wsHub.Run()

	r := http.NewRouter()
	r.SetupRouter(authHandler, chatHandler, messageHandler, tokenService, wsHub)
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
