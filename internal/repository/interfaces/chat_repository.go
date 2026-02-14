package interfaces

import (
	"errors"
	"simpleMessenger/internal/model"
)

var (
	ErrChatNotFound      = errors.New("chat not found")
	ErrChatAlreadyExists = errors.New("chat already exists")
)

type ChatRepo interface {
	Create(chat *model.Chat) error
	GetByID(id uint) (*model.Chat, error)
	Update(chat *model.Chat) error
	Delete(id uint) error
}
