package entity

import "time"

type RefreshToken struct {
	UserId    string    `db:"user_id"`
	TokenHash string    `db:"token_hash"`
	IpAddress string    `db:"ip_address"`
	CreatedAt time.Time `db:"created_at"`
	ExpiresAt time.Time `db:"expires_at"`
}
