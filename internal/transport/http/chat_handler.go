package http

import (
	"github.com/gin-gonic/gin"
	"simpleMessenger/internal/service"
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

func (h *ChatHandler) DeleteChat(c *gin.Context) {
	var req DeleteChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Chat deleted successfully"})
}
