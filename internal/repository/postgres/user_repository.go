package postgres

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"simpleMessenger/internal/model"
	repoInterfaces "simpleMessenger/internal/repository/interfaces"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repoInterfaces.UserRepo {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *model.User) error {
	result := r.db.Create(user)
	if result.Error != nil {
		return fmt.Errorf("create user: %w", result.Error)
	}
	return nil
}

func (r *userRepository) GetByID(id uint) (*model.User, error) {
	user := &model.User{}
	err := r.db.Where("id = ?", id).First(user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repoInterfaces.ErrUserNotFound
		}
		return nil, fmt.Errorf("get user by id: %w", err)
	}
	return user, nil
}

func (r *userRepository) GetByLogin(login string) (*model.User, error) {
	user := &model.User{}
	err := r.db.Where("login = ?", login).First(user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repoInterfaces.ErrUserNotFound
		}
		return nil, fmt.Errorf("get user by login: %w", err)
	}
	return user, nil
}

func (r *userRepository) Update(user *model.User) error {
	result := r.db.Model(model.User{}).Updates(user)
	if result.Error != nil {
		if result.RowsAffected == 0 {
			return repoInterfaces.ErrUserNotFound
		}
		return fmt.Errorf("update user: %w", result.Error)
	}
	return nil
}

func (r *userRepository) Delete(id uint) error {
	result := r.db.Delete(model.User{}, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return repoInterfaces.ErrUserNotFound
		}
		return fmt.Errorf("delete user: %w", result.Error)
	}
	return nil
}
