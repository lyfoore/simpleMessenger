package service

import (
	"fmt"
	"simpleMessenger/internal/model"
	repoInterfaces "simpleMessenger/internal/repository/interfaces"
)

type ChatService struct {
	chatRepo             repoInterfaces.ChatRepo
	chatParticipantsRepo repoInterfaces.ChatParticipantsRepo
	messageRepo          repoInterfaces.MessageRepo
}

func NewChatService(chatRepo repoInterfaces.ChatRepo, chatParticipantsRepo repoInterfaces.ChatParticipantsRepo, messageRepo repoInterfaces.MessageRepo) *ChatService {
	return &ChatService{chatRepo, chatParticipantsRepo, messageRepo}
}

func (s *ChatService) CreateChat(firstUserID, secondUserID uint) error {
	isChatExists, err := s.chatParticipantsRepo.IsChatExists(firstUserID, secondUserID)
	if err != nil {
		return fmt.Errorf("cant check if the chat is already exists: %v", err)
	}

	if isChatExists {
		return fmt.Errorf("chat is already created")
	}

	chat := &model.Chat{}

	err = s.chatRepo.Create(chat)
	if err != nil {
		return fmt.Errorf("cant create chat: %v", err)
	}

	chatParticipantsFirst := &model.ChatParticipants{
		ChatID: chat.ID,
		UserID: firstUserID,
	}

	chatParticipantsSecond := &model.ChatParticipants{
		ChatID: chat.ID,
		UserID: secondUserID,
	}

	err = s.chatParticipantsRepo.Create(chatParticipantsFirst)
	if err != nil {
		_ = s.chatRepo.Delete(chat.ID)
		return fmt.Errorf("cant create chat participants: %v", err)
	}

	err = s.chatParticipantsRepo.Create(chatParticipantsSecond)
	if err != nil {
		_ = s.chatRepo.Delete(chat.ID)
		_ = s.chatParticipantsRepo.Delete(chatParticipantsFirst.ID)
		return fmt.Errorf("cant create chat participants: %v", err)
	}

	return nil
}

func (s *ChatService) DeleteChat(chatID, userID uint) error {
	isUserInChat, err := s.chatParticipantsRepo.IsUserInChat(userID, chatID)
	if err != nil {
		return fmt.Errorf("cant check if user is in chat: %v", err)
	}

	if !isUserInChat {
		return fmt.Errorf("cant delete other user's chat")
	}

	err = s.messageRepo.DeleteAllMessagesInChat(chatID)
	if err != nil {
		return fmt.Errorf("cant delete all messages in chat: %v", err)
	}

	err = s.chatParticipantsRepo.DeleteChat(chatID)
	if err != nil {
		return fmt.Errorf("cant delete chat from chat participants: %v", err)
	}

	err = s.chatRepo.Delete(chatID)
	if err != nil {
		return fmt.Errorf("cant delete chat: %v", err)
	}

	return nil
}
