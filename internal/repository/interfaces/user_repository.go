package interfaces

import (
	"errors"
	"simpleMessenger/internal/model"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrDuplicateEmail    = errors.New("email already registered")
)

type UserRepo interface {
	Create(user *model.User) error
	GetByID(id uint) (*model.User, error)
	GetByLogin(login string) (*model.User, error)
	Update(user *model.User) error
	Delete(id uint) error
}
