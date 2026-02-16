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

	var req struct {
		Text string `json:"text"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("failed to unmarshal json with text")
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to unmarshal text message"})
		return
	}

	sendMessageRequest := &service.SendMessageRequest{
		Text:   req.Text,
		UserID: userID,
		ChatID: uint(chatID),
	}

	_, err = h.messageService.SendMessage(sendMessageRequest)

	if err != nil {
		log.Printf("failed to send message: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to send message"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Message sent"})
}

func (h *MessageHandler) GetMessages(c *gin.Context) {
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

	limitStr := c.DefaultQuery("limit", "20")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		log.Printf("failed to parse limit param: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit parameter"})
		return
	}

	if limit > 100 {
		limit = 100
	}

	messages, err := h.messageService.GetMessages(uint(chatID), limit, userID)

	if err != nil {
		log.Printf("failed to get messages for chat %d: %v", chatID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve messages"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"messages": messages})
}

func (h *MessageHandler) DeleteMessage(c *gin.Context) {
	userID, err := GetUserIdFromContext(c)

	if err != nil {
		log.Printf("failed to get userID from context: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "failed to get a userID"})
		return
	}

	messageIDStr := c.Param("messageId")
	messageID, err := strconv.ParseUint(messageIDStr, 10, 64)

	if err != nil {
		log.Printf("failed to parse message id param: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid messageID"})
		return
	}

	err = h.messageService.DeleteMessage(uint(messageID), userID)

	if err != nil {
		log.Printf("failed to delete message: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete message"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Message deleted successfully"})
}
