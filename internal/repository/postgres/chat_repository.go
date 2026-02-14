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

func (r *chatRepository) GetChats(userID uint, limit int) ([]*model.Chat, error) {
	var chats []*model.Chat

	query := r.db.
		Where("id IN (SELECT chat_id FROM chat_participants WHERE user_id = ?)", userID).
		Order("last_message_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&chats).Error
	if err != nil {
		return nil, fmt.Errorf("get chats: %w", err)
	}

	return chats, nil
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
