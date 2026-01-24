package models

import "time"

type User struct {
	Id           int64  `json:"id"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	PasswordHash string
}

type RefreshToken struct {
	Id                int64
	UserId            int64
	Hash              string
	CreatedAt         time.Time
	LastUsedAt        *time.Time
	ExpiresAt         time.Time
	RevokedAt         *time.Time
	ReplacedByTokenId *int64
}
