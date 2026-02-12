package postgres

import (
	"fmt"
	"gorm.io/gorm"
	"simpleMessenger/internal/model"
	repoInterfaces "simpleMessenger/internal/repository/interfaces"
)

type chatRepository struct {
	db *gorm.DB
}

func NewChatRepository(db *gorm.DB) repoInterfaces.ChatRepo {
	return &chatRepository{db: db}
}

func (r *chatRepository) Create(chat *model.Chat) error {
	result := r.db.Create(chat)
	if result.Error != nil {
		return fmt.Errorf("create chat: %w", result.Error)
	}
	return nil
}

func (r *chatRepository) GetByID(id uint) (*model.Chat, error) {
	chat := &model.Chat{}
	err := r.db.Where("id = ?", id).First(chat).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, repoInterfaces.ErrChatNotFound
		}
		return nil, fmt.Errorf("get chat by id: %w", err)
	}
	return chat, nil
}

func (r *chatRepository) Update(chat *model.Chat) error {
	result := r.db.Model(model.Chat{}).Updates(chat)
	if result.Error != nil {
		if result.RowsAffected == 0 {
			return repoInterfaces.ErrChatNotFound
		}
		return fmt.Errorf("update chat: %w", result.Error)
	}
	return nil
}

func (r *chatRepository) Delete(id uint) error {
	result := r.db.Delete(model.Chat{}, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return repoInterfaces.ErrChatNotFound
		}
		return fmt.Errorf("delete chat: %w", result.Error)
	}
	return nil
}
