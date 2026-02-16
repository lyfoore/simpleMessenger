package postgres

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"simpleMessenger/internal/model"
	repoInterfaces "simpleMessenger/internal/repository/interfaces"
)

type chatParticipantsRepository struct {
	db *gorm.DB
}

func NewChatParticipantsRepository(db *gorm.DB) repoInterfaces.ChatParticipantsRepo {
	return &chatParticipantsRepository{db: db}
}

func (c *chatParticipantsRepository) Create(chatParticipants *model.ChatParticipants) error {
	result := c.db.Create(chatParticipants)
	if result.Error != nil {
		return fmt.Errorf("create chatParticipants: %w", result.Error)
	}
	return nil
}

func (c *chatParticipantsRepository) GetByID(id uint) (*model.ChatParticipants, error) {
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

func (c *chatParticipantsRepository) GetChatParticipantsByChatID(chatId uint) ([]uint, error) {
	var chatParticipants []uint
	result := c.db.Model(&model.ChatParticipants{}).Select("user_id").Where("chat_id = ?", chatId).Find(&chatParticipants)

	if result.Error != nil {
		return nil, fmt.Errorf("get chat participants by chat id: %v", result.Error)
	}

	return chatParticipants, nil
}

func (c *chatParticipantsRepository) IsChatExists(firstUserID, secondUserID uint) (bool, error) {
	var chatID uint
	err := c.db.Model(&model.ChatParticipants{}).Select("chat_id").Where("user_id IN (?, ?)", firstUserID, secondUserID).Group("chat_id").Having("COUNT(DISTINCT user_id) = ?", 2).Having("COUNT(*) = ?", 2).Having("(SELECT COUNT(*) FROM chat_participants WHERE chat_id = chat_participants.chat_id) = ?", 2).Take(&chatID).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (c *chatParticipantsRepository) IsUserInChat(userID, chatID uint) (bool, error) {
	err := c.db.Model(&model.ChatParticipants{}).Where("user_id = ? AND chat_id = ?", userID, chatID).First(&model.ChatParticipants{}).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (c *chatParticipantsRepository) Update(participants *model.ChatParticipants) error {
	result := c.db.Model(&model.ChatParticipants{}).Updates(participants)
	if result.Error != nil {
		if result.RowsAffected == 0 {
			return repoInterfaces.ErrChatParticipantsNotFound
		}
		return fmt.Errorf("update chatParticipants: %w", result.Error)
	}
	return nil
}

func (c *chatParticipantsRepository) Delete(id uint) error {
	result := c.db.Delete(&model.ChatParticipants{}, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return repoInterfaces.ErrChatParticipantsNotFound
		}
		return fmt.Errorf("delete chatParticipants: %w", result.Error)
	}
	return nil
}

func (c *chatParticipantsRepository) DeleteChat(chatID uint) error {
	err := c.db.Model(&model.ChatParticipants{}).Delete(&model.ChatParticipants{}, "chat_id = ?", chatID).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("delete chatParticipants for whole chat: %w", repoInterfaces.ErrChatParticipantsNotFound)
		}
		return fmt.Errorf("delete chatParticipants for whole chat: %w", err)
	}

	return nil
}
