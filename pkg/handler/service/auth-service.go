package service

import (
	"github.com/bremcm/todo-app"
	"github.com/bremcm/todo-app/pkg/handler/service/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(user todo.User) (int, error) {
	hash, err := s.generatePasswordHash(user.Password)
	if err != nil {
		return 0, err
	}
	user.Password = hash
	return s.repo.CreateUser(user)
}

func (s *AuthService) generatePasswordHash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
