package service

import (
	"evelinqua/internal/repository"
)

type AuthService struct {
	repo *repository.AuthRepository
}

func NewAuthService(repo *repository.AuthRepository) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) Login(auth repository.Auth) error {
	return s.repo.Login(auth)
}
