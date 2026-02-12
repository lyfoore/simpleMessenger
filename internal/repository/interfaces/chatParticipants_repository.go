package interfaces

import (
	"errors"
	"simpleMessenger/internal/model"
)

var (
	ErrChatParticipantsNotFound = errors.New("chat participants not found")
)

type ChatParticipantsRepo interface {
	Create(chatParticipants *model.ChatParticipants) error
	GetByID(id uint) (*model.ChatParticipants, error)
	IsChatExists(firstUserID, secondUserID uint) (bool, error)
	IsUserInChat(userID, chatID uint) (bool, error)
	Update(participants *model.ChatParticipants) error
	Delete(id uint) error
	DeleteChat(chatID uint) error
}
