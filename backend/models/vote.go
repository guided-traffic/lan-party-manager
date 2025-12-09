package models

import "time"

// Vote represents a vote from one user to another
type Vote struct {
	ID            uint64    `json:"id"`
	FromUserID    uint64    `json:"from_user_id"`
	ToUserID      uint64    `json:"to_user_id"`
	AchievementID string    `json:"achievement_id"`
	Points        int       `json:"points"`
	CreatedAt     time.Time `json:"created_at"`
}

// VoteWithDetails includes user information for display
type VoteWithDetails struct {
	ID            uint64      `json:"id"`
	FromUser      PublicUser  `json:"from_user"`
	ToUser        PublicUser  `json:"to_user"`
	AchievementID string      `json:"achievement_id"`
	Achievement   Achievement `json:"achievement"`
	Points        int         `json:"points"`
	CreatedAt     time.Time   `json:"created_at"`
}

// CreateVoteRequest is the request body for creating a vote
type CreateVoteRequest struct {
	ToUserID      uint64 `json:"to_user_id" binding:"required"`
	AchievementID string `json:"achievement_id" binding:"required"`
	Points        int    `json:"points"` // 1-3 points, defaults to 1 if not provided
}
