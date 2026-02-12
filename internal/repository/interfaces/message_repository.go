package interfaces

import (
	"errors"
	"simpleMessenger/internal/model"
)

var (
	ErrMessageNotFound = errors.New("message not found")
)

type MessageRepo interface {
	Create(message *model.Message) error
	GetByID(id uint) (*model.Message, error)
	Update(message *model.Message) error
	Delete(id uint) error
	DeleteAllMessagesInChat(chatID uint) error
}
