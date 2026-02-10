package service

import (
	"errors"
	"simpleTodoList/internal/model"
	"simpleTodoList/internal/repository/interfaces"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
)

type LoginResponse struct {
	User      *model.User
	Token     string
	ExpiresIn int64
}

type AuthService struct {
	userRepo     interfaces.UserRepo
	tokenService TokenService
}

func NewAuthService(userRepo interfaces.UserRepo, tokenService TokenService) *AuthService {
	return &AuthService{userRepo: userRepo, tokenService: tokenService}
}

func (s *AuthService) Register(user *model.User) error {
	existing, _ := s.userRepo.GetByID(user.ID)
	if existing != nil {
		return ErrUserAlreadyExists
	}
	return s.userRepo.Create(user)
}

func (s *AuthService) Login(username string) (*LoginResponse, error) {
	user, err := s.userRepo.GetByLogin(username)
	if err != nil {
		return nil, err
	}

	token, err := s.tokenService.GenerateToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		User:      user,
		Token:     token,
		ExpiresIn: AccessTokenExpireDuration,
	}, nil

}
