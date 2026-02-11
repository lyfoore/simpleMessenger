package postgres

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"simpleTodoList/internal/model"
	repoInterfaces "simpleTodoList/internal/repository/interfaces"
)

type chatParticipantsRepository struct {
	db *gorm.DB
}

func NewChatParticipantsRepository(db *gorm.DB) repoInterfaces.ChatParticipantsRepo {
	return &chatParticipantsRepository{db: db}
}

func (c chatParticipantsRepository) Create(chatParticipants *model.ChatParticipants) error {
	result := c.db.Create(chatParticipants)
	if result.Error != nil {
		return fmt.Errorf("create chatParticipants: %w", result.Error)
	}
	return nil
}

func (c chatParticipantsRepository) GetByID(id uint) (*model.ChatParticipants, error) {
	chatParticipants := &model.ChatParticipants{}
	result := c.db.Where("id = ?", id).First(chatParticipants)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, repoInterfaces.ErrChatParticipantsNotFound
		}
		return nil, fmt.Errorf("get chatParticipants: %w", result.Error)
	}
	return chatParticipants, nil
}

func (c chatParticipantsRepository) Update(participants *model.ChatParticipants) error {
	result := c.db.Model(model.ChatParticipants{}).Updates(participants)
	if result.Error != nil {
		if result.RowsAffected == 0 {
			return repoInterfaces.ErrChatParticipantsNotFound
		}
		return fmt.Errorf("update chatParticipants: %w", result.Error)
	}
	return nil
}

func (c chatParticipantsRepository) Delete(id uint) error {
	result := c.db.Delete(model.ChatParticipants{}, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return repoInterfaces.ErrChatParticipantsNotFound
		}
		return fmt.Errorf("delete chatParticipants: %w", result.Error)
	}
	return nil
}
