package http

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"simpleMessenger/internal/service"
	"strconv"
)

type ChatHandler struct {
	chatService *service.ChatService
}

func NewChatHandler(chatService *service.ChatService) *ChatHandler {
	return &ChatHandler{chatService: chatService}
}

func (h *ChatHandler) CreateChat(c *gin.Context) {
	userID, err := GetUserIdFromContext(c)

	if err != nil {
		log.Printf("failed to get userID from context: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "failed to get a userID"})
		return
	}

	var req struct {
		CompanionID uint `json:"companionId"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("failed to parse body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to unmarshal companion id"})
		return
	}

	err = h.chatService.CreateChat(userID, req.CompanionID)
	if err != nil {
		log.Printf("failed to create chat: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create chat"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Chat created successfully"})
}

func (h *ChatHandler) GetChats(c *gin.Context) {
	userID, err := GetUserIdFromContext(c)

	if err != nil {
		log.Printf("failed to get userID from context: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "failed to get a userID"})
		return
	}

	limitStr := c.DefaultQuery("limit", "20")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit parameter"})
		return
	}

	if limit > 100 {
		limit = 100
	}

	chats, err := h.chatService.GetChats(userID, limit)
	if err != nil {
		log.Printf("failed to get chats for user %d: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve chats"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"chats": chats})
}

func (h *ChatHandler) DeleteChat(c *gin.Context) {
	userID, err := GetUserIdFromContext(c)

	if err != nil {
		log.Printf("failed to get userID from context: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "failed to get a userID"})
		return
	}

	chatIDStr := c.Param("chatId")
	chatID, err := strconv.ParseUint(chatIDStr, 10, 64)

	if err != nil {
		log.Printf("failed to parse chat id param: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid chatID"})
		return
	}

	err = h.chatService.DeleteChat(uint(chatID), userID)

	if err != nil {
		log.Printf("failed to delete chat %d: %v", chatID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete chat"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Chat deleted successfully"})
}
