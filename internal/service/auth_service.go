package service

import (
	"simpleTodoList/internal/repository/interfaces"
)

type AuthService struct {
	userRepo interfaces.UserRepo
}

func NewAuthService(userRepo interfaces.UserRepo) *AuthService {
	return &AuthService{userRepo: userRepo}
}
