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
	Update(participants *model.ChatParticipants) error
	Delete(id uint) error
}
