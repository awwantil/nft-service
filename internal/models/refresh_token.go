package models

import (
	"time"
)

// UserToken represents the structure of a user token
type UserToken struct {
	ID               int64     `json:"id"`
	UserID           int64     `json:"user_id"`
	Token            string    `json:"token"`
	RefreshToken     string    `json:"refresh_token"`
	ExpiredAt        time.Time `json:"expired_at"`
	RefreshExpiredAt time.Time `json:"refresh_expired_at"`
}
