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

type CreateChatRequest struct {
	FirstUserID  uint `json:"first_user_id"`
	SecondUserID uint `json:"second_user_id"`
}

type DeleteChatRequest struct {
	ChatID uint `json:"chat_id"`
}

func NewChatHandler(chatService *service.ChatService) *ChatHandler {
	return &ChatHandler{chatService: chatService}
}

func (h *ChatHandler) CreateChat(c *gin.Context) {
	var req CreateChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err := h.chatService.CreateChat(req.FirstUserID, req.SecondUserID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Chat created successfully"})
}

func (h *ChatHandler) GetChats(c *gin.Context) {
	value, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}
	userID, ok := value.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user id type"})
		return
	}

	limitStr := c.DefaultQuery("limit", "20")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit parameter"})
		return
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
	var req DeleteChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Chat deleted successfully"})
}
