package postgres

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"simpleMessenger/internal/model"
	repoInterfaces "simpleMessenger/internal/repository/interfaces"
)

type messageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) repoInterfaces.MessageRepo {
	return &messageRepository{db: db}
}

func (r *messageRepository) Create(message *model.Message) error {
	result := r.db.Create(message)
	if result.Error != nil {
		return fmt.Errorf("create message: %w", result.Error)
	}
	return nil
}

func (r *messageRepository) GetByID(id uint) (*model.Message, error) {
	message := &model.Message{}
	err := r.db.Where("id = ?", id).First(message).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repoInterfaces.ErrMessageNotFound
		}
		return nil, fmt.Errorf("get message by id: %w", err)
	}
	return message, nil
}

func (r *messageRepository) GetMessagesByChatID(chatID uint, limit int) ([]*model.Message, error) {
	var messages []*model.Message

	query := r.db.Model(&model.Message{}).
		Where("chat_id = ?", chatID).
		Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&messages).Error
	if err != nil {
		return nil, fmt.Errorf("get messages in chat: %w", err)
	}

	return messages, nil
}

func (r *messageRepository) Update(message *model.Message) error {
	result := r.db.Model(model.Message{}).Updates(message)
	if result.Error != nil {
		if result.RowsAffected == 0 {
			return repoInterfaces.ErrMessageNotFound
		}
		return fmt.Errorf("update message: %w", result.Error)
	}
	return nil
}

func (r *messageRepository) Delete(id uint) error {
	result := r.db.Delete(model.Message{}, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return repoInterfaces.ErrMessageNotFound
		}
		return fmt.Errorf("delete message: %w", result.Error)
	}
	return nil
}

func (r *messageRepository) DeleteAllMessagesInChat(chatID uint) error {
	err := r.db.Delete(&model.Message{}, "chat_id = ?", chatID).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return repoInterfaces.ErrMessageNotFound
		}
		return fmt.Errorf("delete all messages in chat: %w", err)
	}
	return nil
}
