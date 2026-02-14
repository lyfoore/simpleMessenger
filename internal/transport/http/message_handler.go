package http

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"simpleMessenger/internal/service"
	"strconv"
)

type GetMessagesRequest struct {
	ChatID uint `json:"chatId"`
	UserID uint `json:"userId"`
}

type MessageHandler struct {
	messageService *service.MessageService
}

func NewMessageHandler(messageService *service.MessageService) *MessageHandler {
	return &MessageHandler{messageService: messageService}
}

func (h *MessageHandler) SendMessage(c *gin.Context) {
	var req service.SendMessageRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err := h.messageService.SendMessage(&req)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Message sent successfully"})
}

func (h *MessageHandler) GetMessages(c *gin.Context) {
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

	var req GetMessagesRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	limitStr := c.DefaultQuery("limit", "20")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit parameter"})
		return
	}

	messages, err := h.messageService.GetMessages(req.ChatID, limit, userID)

	if err != nil {
		log.Printf("failed to get messages for chat %d: %v", req.ChatID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve messages"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"messages": messages})
}

func (h *MessageHandler) DeleteMessage(c *gin.Context) {
	var req service.DeleteMessageRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err := h.messageService.DeleteMessage(req.MessageID, req.UserID)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Message deleted successfully"})
}
