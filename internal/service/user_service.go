package service

import (
	"errors"
	"simpleMessenger/internal/model"
	"simpleMessenger/internal/repository/interfaces"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type UserService struct {
	userRepo interfaces.UserRepo
}

func NewUserService(userRepo interfaces.UserRepo) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) GetUserByID(id uint) (*model.User, error) {
	user, err := s.userRepo.GetByID(id)

	if err != nil {
		if errors.Is(err, interfaces.ErrUserNotFound) {
			return nil, ErrUserNotFound
		}

		return nil, err
	}

	return user, nil
}

func (s *UserService) GetUserByLogin(login string) (*model.User, error) {
	user, err := s.userRepo.GetByLogin(login)

	if err != nil {
		if errors.Is(err, interfaces.ErrUserNotFound) {
			return nil, ErrUserNotFound
		}

		return nil, err
	}

	return user, nil
}

func (s *UserService) SearchUsersByLogin(login string, limit int) ([]*model.User, error) {
	if limit <= 0 {
		return nil, errors.New("invalid limit value")
	}

	if limit > 20 {
		limit = 20
	}

	users, err := s.userRepo.SearchByLogin(login, limit)

	if err != nil {
		if errors.Is(err, interfaces.ErrUserNotFound) {
			return nil, ErrUserNotFound
		}

		return nil, err
	}

	return users, nil
}

func (s *UserService) UpdateUser(user *model.User) error {
	err := s.userRepo.Update(user)

	if err != nil {
		if errors.Is(err, interfaces.ErrUserNotFound) {
			return ErrUserNotFound
		}
		return err
	}

	return nil
}
