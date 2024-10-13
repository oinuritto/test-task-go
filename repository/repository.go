package repository

import (
	"github.com/jmoiron/sqlx"
	"testTask/entity"
)

type Authorization interface {
	CreateUser(user entity.User) (string, error)
	GetUser(id string) (entity.User, error)
	CreateRefreshToken(token entity.RefreshToken) error
	GetRefreshTokenById(id string) (entity.RefreshToken, error)
	DeleteRefreshTokenById(id string) error
}

type Repository struct {
	Authorization
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
	}
}
