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
	Create(user *model.Chat) error
	GetByID(id uint) (*model.Chat, error)
	Update(user *model.Chat) error
	Delete(id uint) error
}
