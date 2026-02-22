package service

import (
	"fmt"
	"log"
	"simpleMessenger/internal/model"
	repoInterfaces "simpleMessenger/internal/repository/interfaces"
	"time"
)

type ChatService struct {
	chatRepo             repoInterfaces.ChatRepo
	chatParticipantsRepo repoInterfaces.ChatParticipantsRepo
	messageRepo          repoInterfaces.MessageRepo
	userRepo             repoInterfaces.UserRepo
}

func NewChatService(chatRepo repoInterfaces.ChatRepo, chatParticipantsRepo repoInterfaces.ChatParticipantsRepo, messageRepo repoInterfaces.MessageRepo, userRepo repoInterfaces.UserRepo) *ChatService {
	return &ChatService{chatRepo, chatParticipantsRepo, messageRepo, userRepo}
}

func (s *ChatService) CreateChat(firstUserID, secondUserID uint) error {
	isChatExists, err := s.chatParticipantsRepo.IsChatExists(firstUserID, secondUserID)
	if err != nil {
		return fmt.Errorf("cant check if the chat is already exists: %w", err)
	}

	if isChatExists {
		return fmt.Errorf("chat is already created")
	}

	firstUser, err := s.userRepo.GetByID(firstUserID)
	if err != nil {
		return fmt.Errorf("cant get first user by id: %w", err)
	}

	secondUser, err := s.userRepo.GetByID(secondUserID)
	if err != nil {
		return fmt.Errorf("cant get second user by id: %w", err)
	}

	chatName := formatChatName(firstUser.Login, secondUser.Login)

	chat := &model.Chat{
		LastMessageAt: time.Now(),
		Name:          chatName,
	}

	err = s.chatRepo.Create(chat)
	if err != nil {
		return fmt.Errorf("cant create chat: %w", err)
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
		return fmt.Errorf("cant create chat participants: %w", err)
	}

	err = s.chatParticipantsRepo.Create(chatParticipantsSecond)
	if err != nil {
		_ = s.chatRepo.Delete(chat.ID)
		_ = s.chatParticipantsRepo.Delete(chatParticipantsFirst.ID)
		return fmt.Errorf("cant create chat participants: %w", err)
	}

	return nil
}

func formatChatName(login1, login2 string) string {
	if login1 < login2 {
		return login1 + " • " + login2
	}
	return login2 + " • " + login1
}

func (s *ChatService) GetChats(userID uint, limit int) ([]*model.Chat, error) {
	chats, err := s.chatRepo.GetChats(userID, limit)
	if err != nil {
		return nil, fmt.Errorf("cant get chats by userID: %w", err)
	}
	return chats, nil
}

func (s *ChatService) GetUsersInChat(chatId uint) ([]uint, error) {
	userIDs, err := s.chatParticipantsRepo.GetChatParticipantsByChatID(chatId)

	if err != nil {
		log.Printf("cant get users by chat id: %v", chatId)
		return nil, fmt.Errorf("cant get users by chat id: %w", err)
	}

	return userIDs, nil
}

func (s *ChatService) DeleteChat(chatID, userID uint) error {
	isUserInChat, err := s.chatParticipantsRepo.IsUserInChat(userID, chatID)
	if err != nil {
		return fmt.Errorf("cant check if user is in chat: %w", err)
	}

	if !isUserInChat {
		return fmt.Errorf("cant delete other user's chat")
	}

	err = s.messageRepo.DeleteAllMessagesInChat(chatID)
	if err != nil {
		return fmt.Errorf("cant delete all messages in chat: %w", err)
	}

	err = s.chatParticipantsRepo.DeleteChat(chatID)
	if err != nil {
		return fmt.Errorf("cant delete chat from chat participants: %w", err)
	}

	err = s.chatRepo.Delete(chatID)
	if err != nil {
		return fmt.Errorf("cant delete chat: %w", err)
	}

	return nil
}
