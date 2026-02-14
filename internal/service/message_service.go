package service

import (
	"errors"
	"fmt"
	"simpleMessenger/internal/model"
	repoInterfaces "simpleMessenger/internal/repository/interfaces"
	"strings"
	"time"
)

type MessageService struct {
	messageRepo          repoInterfaces.MessageRepo
	chatRepo             repoInterfaces.ChatRepo
	chatParticipantsRepo repoInterfaces.ChatParticipantsRepo
}

func NewMessageService(messageRepo repoInterfaces.MessageRepo, chatRepo repoInterfaces.ChatRepo, chatParticipantsRepo repoInterfaces.ChatParticipantsRepo) *MessageService {
	return &MessageService{
		messageRepo:          messageRepo,
		chatRepo:             chatRepo,
		chatParticipantsRepo: chatParticipantsRepo,
	}
}

type SendMessageRequest struct {
	Text   string `json:"text"`
	UserID uint   `json:"user_id"`
	ChatID uint   `json:"chat_id"`
}

type DeleteMessageRequest struct {
	UserID    uint `json:"user_id"`
	MessageID uint `json:"message_id"`
}

func (s *MessageService) SendMessage(req *SendMessageRequest) error {
	if req == nil {
		return errors.New("nil request")
	}

	if strings.TrimSpace(req.Text) == "" {
		return errors.New("message text cannot be empty")
	}

	isUserInChat, err := s.chatParticipantsRepo.IsUserInChat(req.UserID, req.ChatID)
	if err != nil {
		return fmt.Errorf("cant check if user is in chat: %w", err)
	}

	if !isUserInChat {
		return fmt.Errorf("cant send message. user is not in chat")
	}

	msg := &model.Message{
		Text:   req.Text,
		ChatID: req.ChatID,
		UserID: req.UserID,
	}

	err = s.messageRepo.Create(msg)
	if err != nil {
		return fmt.Errorf("cant send message: %w", err)
	}

	chat, err := s.chatRepo.GetByID(req.ChatID)

	if err != nil {
		_ = s.messageRepo.Delete(msg.ID)
		return fmt.Errorf("cant get chat for updating: %w", err)
	}

	chat.LastMessageAt = time.Now()

	err = s.chatRepo.Update(chat)

	if err != nil {
		_ = s.messageRepo.Delete(msg.ID)
		return fmt.Errorf("cant update last_message_at in chat: %w", err)
	}

	return nil
}

func (s *MessageService) DeleteMessage(messageID, userID uint) error {
	msg, err := s.messageRepo.GetByID(messageID)
	if err != nil {
		return fmt.Errorf("cant get message: %w", err)
	}

	if msg.UserID != userID {
		return fmt.Errorf("cant delete another user's message")
	}

	err = s.messageRepo.Delete(messageID)
	if err != nil {
		return fmt.Errorf("cant delete message: %w", err)
	}

	return nil
}
