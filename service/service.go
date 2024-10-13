package service

import (
	"testTask/entity"
	"testTask/repository"
)

type Authorization interface {
	CreateUser(user entity.User) (string, error)
	GenerateAccessToken(id, ip string) (string, error)
	ParseToken(token string) (string, error)
	GenerateRefreshToken(id, ip string) (string, error)
	RefreshTokens(token string, ip string) (string, string, error)
}

type Service struct {
	Authorization
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
	}
}
