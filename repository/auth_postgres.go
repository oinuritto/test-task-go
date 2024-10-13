package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"testTask/entity"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) CreateUser(user entity.User) (string, error) {
	var id string
	query := fmt.Sprintf("INSERT INTO %s (id, email, password) values ($1, $2, $3) RETURNING id", usersTable)

	row := r.db.QueryRow(query, user.Id, user.Email, user.Password)
	if err := row.Scan(&id); err != nil {
		return "", err
	}

	return id, nil
}

func (r *AuthPostgres) GetUser(id string) (entity.User, error) {
	var user entity.User
	query := fmt.Sprintf("SELECT * FROM %s WHERE id=$1", usersTable)
	err := r.db.Get(&user, query, id)

	return user, err
}

func (r *AuthPostgres) CreateRefreshToken(token entity.RefreshToken) error {
	query := fmt.Sprintf("INSERT INTO %s (user_id, token_hash, ip_address, created_at, expires_at) "+
		"values ($1, $2, $3, $4, $5)", refreshTokensTable)
	_, err := r.db.Exec(query, token.UserId, token.TokenHash, token.IpAddress, token.CreatedAt, token.ExpiresAt)

	if err != nil {
		return err
	}

	return nil
}

func (r *AuthPostgres) GetRefreshTokenById(id string) (entity.RefreshToken, error) {
	var refreshToken entity.RefreshToken
	query := fmt.Sprintf("SELECT * FROM %s WHERE user_id=$1", refreshTokensTable)
	err := r.db.Get(&refreshToken, query, id)

	return refreshToken, err
}

func (r *AuthPostgres) DeleteRefreshTokenById(id string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE user_id=$1", refreshTokensTable)
	_, err := r.db.Exec(query, id)

	return err
}
