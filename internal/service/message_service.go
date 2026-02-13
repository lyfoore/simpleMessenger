package service

import (
	"fmt"
	"simpleMessenger/internal/model"
	repoInterfaces "simpleMessenger/internal/repository/interfaces"
)

type MessageService struct {
	messageRepo          repoInterfaces.MessageRepo
	chatParticipantsRepo repoInterfaces.ChatParticipantsRepo
}

func NewMessageService(messageRepo repoInterfaces.MessageRepo, chatParticipantsRepo repoInterfaces.ChatParticipantsRepo) *MessageService {
	return &MessageService{messageRepo: messageRepo, chatParticipantsRepo: chatParticipantsRepo}
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
	isUserInChat, err := s.chatParticipantsRepo.IsUserInChat(req.UserID, req.ChatID)
	if err != nil {
		return fmt.Errorf("cant check if user is in chat: %v", err)
	}

	if isUserInChat != true {
		return fmt.Errorf("cant send message. user is not in chat")
	}

	msg := &model.Message{
		Text:   req.Text,
		ChatID: req.ChatID,
		UserID: req.UserID,
	}

	err = s.messageRepo.Create(msg)
	if err != nil {
		return fmt.Errorf("cant send message: %v", err)
	}

	return nil
}

func (s *MessageService) DeleteMessage(messageID, userID uint) error {
	msg, err := s.messageRepo.GetByID(messageID)
	if err != nil {
		return fmt.Errorf("cant get message: %v", err)
	}

	if msg.UserID != userID {
		return fmt.Errorf("cant delete another user's message")
	}

	err = s.messageRepo.Delete(messageID)
	if err != nil {
		return fmt.Errorf("cant delete message: %v", err)
	}

	return nil
}
