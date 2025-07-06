package models

import (
	"time"
)

// User represents the structure of a user entity
type User struct {
	ID            int64     `json:"id"`
	Phone         string    `json:"phone"`
	Password      string    `json:"-"`
	Salt          []byte    `json:"-"`
	RoleID        RoleId    `json:"role_id"`
	CreatedAt     time.Time `json:"-"`
	UpdatedAt     time.Time `json:"-"`
	DeletedAt     time.Time `json:"-"`
	LastVisitedAt time.Time `json:"-"`
}
