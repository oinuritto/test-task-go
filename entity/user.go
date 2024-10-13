package entity

import "github.com/google/uuid"

type User struct {
	Id       uuid.UUID `json:"id" db:"id"`
	Email    string    `json:"email" binding:"required"`
	Password string    `json:"password" binding:"required"`
}
