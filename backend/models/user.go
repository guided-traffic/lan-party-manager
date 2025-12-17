package models

import "time"

// User represents a registered player
type User struct {
	ID                 uint64     `json:"id"`
	SteamID            string     `json:"steam_id"`
	Username           string     `json:"username"`
	AvatarURL          string     `json:"avatar_url"`
	AvatarSmall        string     `json:"avatar_small"`
	ProfileURL         string     `json:"profile_url"`
	Credits            int        `json:"credits"`
	LastCreditAt       time.Time  `json:"last_credit_at"`
	LastGamesRefreshAt *time.Time `json:"last_games_refresh_at"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

// PublicUser represents the public-facing user data (no sensitive info)
type PublicUser struct {
	ID          uint64 `json:"id"`
	SteamID     string `json:"steam_id"`
	Username    string `json:"username"`
	AvatarURL   string `json:"avatar_url"`
	AvatarSmall string `json:"avatar_small"`
	ProfileURL  string `json:"profile_url"`
}

// ToPublic converts a User to PublicUser
func (u *User) ToPublic() PublicUser {
	return PublicUser{
		ID:          u.ID,
		SteamID:     u.SteamID,
		Username:    u.Username,
		AvatarURL:   u.AvatarURL,
		AvatarSmall: u.AvatarSmall,
		ProfileURL:  u.ProfileURL,
	}
}

// BannedUser represents a banned player
type BannedUser struct {
	ID        uint64    `json:"id"`
	SteamID   string    `json:"steam_id"`
	Username  string    `json:"username"`
	Reason    string    `json:"reason"`
	BannedBy  string    `json:"banned_by"`
	BannedAt  time.Time `json:"banned_at"`
}

// AdminUserInfo represents user info for admin view
type AdminUserInfo struct {
	ID          uint64    `json:"id"`
	SteamID     string    `json:"steam_id"`
	Username    string    `json:"username"`
	AvatarSmall string    `json:"avatar_small"`
	CreatedAt   time.Time `json:"created_at"`
}
