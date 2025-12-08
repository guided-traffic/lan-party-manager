package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/guided-traffic/lan-party-manager/backend/database"
	"github.com/guided-traffic/lan-party-manager/backend/models"
)

// UserRepository handles user database operations
type UserRepository struct{}

// NewUserRepository creates a new user repository
func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

// Create creates a new user in the database
func (r *UserRepository) Create(user *models.User) error {
	result, err := database.DB.Exec(`
		INSERT INTO users (steam_id, username, avatar_url, avatar_small, profile_url, credits, last_credit_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		user.SteamID, user.Username, user.AvatarURL, user.AvatarSmall, user.ProfileURL, user.Credits, user.LastCreditAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	user.ID = uint64(id)
	return nil
}

// GetByID finds a user by ID
func (r *UserRepository) GetByID(id uint64) (*models.User, error) {
	user := &models.User{}
	err := database.DB.QueryRow(`
		SELECT id, steam_id, username, avatar_url, avatar_small, profile_url, credits, last_credit_at, created_at, updated_at
		FROM users WHERE id = ?`, id,
	).Scan(&user.ID, &user.SteamID, &user.Username, &user.AvatarURL, &user.AvatarSmall, &user.ProfileURL,
		&user.Credits, &user.LastCreditAt, &user.CreatedAt, &user.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	return user, nil
}

// GetBySteamID finds a user by Steam ID
func (r *UserRepository) GetBySteamID(steamID string) (*models.User, error) {
	user := &models.User{}
	err := database.DB.QueryRow(`
		SELECT id, steam_id, username, avatar_url, avatar_small, profile_url, credits, last_credit_at, created_at, updated_at
		FROM users WHERE steam_id = ?`, steamID,
	).Scan(&user.ID, &user.SteamID, &user.Username, &user.AvatarURL, &user.AvatarSmall, &user.ProfileURL,
		&user.Credits, &user.LastCreditAt, &user.CreatedAt, &user.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user by steam id: %w", err)
	}

	return user, nil
}

// GetAll returns all users
func (r *UserRepository) GetAll() ([]models.User, error) {
	rows, err := database.DB.Query(`
		SELECT id, steam_id, username, avatar_url, avatar_small, profile_url, credits, last_credit_at, created_at, updated_at
		FROM users ORDER BY username`)
	if err != nil {
		return nil, fmt.Errorf("failed to get all users: %w", err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.SteamID, &user.Username, &user.AvatarURL, &user.AvatarSmall, &user.ProfileURL,
			&user.Credits, &user.LastCreditAt, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user row: %w", err)
		}
		users = append(users, user)
	}

	return users, nil
}

// Update updates a user's profile information
func (r *UserRepository) Update(user *models.User) error {
	_, err := database.DB.Exec(`
		UPDATE users
		SET username = ?, avatar_url = ?, avatar_small = ?, profile_url = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?`,
		user.Username, user.AvatarURL, user.AvatarSmall, user.ProfileURL, user.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

// UpdateCredits updates a user's credits
func (r *UserRepository) UpdateCredits(userID uint64, credits int, lastCreditAt time.Time) error {
	_, err := database.DB.Exec(`
		UPDATE users
		SET credits = ?, last_credit_at = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?`,
		credits, lastCreditAt, userID,
	)
	if err != nil {
		return fmt.Errorf("failed to update credits: %w", err)
	}
	return nil
}

// DeductCredit deducts one credit from a user (atomic operation)
func (r *UserRepository) DeductCredit(userID uint64) error {
	result, err := database.DB.Exec(`
		UPDATE users
		SET credits = credits - 1, updated_at = CURRENT_TIMESTAMP
		WHERE id = ? AND credits > 0`,
		userID,
	)
	if err != nil {
		return fmt.Errorf("failed to deduct credit: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("insufficient credits")
	}

	return nil
}

// ResetAllCredits sets all users' credits to 0
func (r *UserRepository) ResetAllCredits() (int64, error) {
	result, err := database.DB.Exec(`
		UPDATE users
		SET credits = 0, updated_at = CURRENT_TIMESTAMP`)
	if err != nil {
		return 0, fmt.Errorf("failed to reset all credits: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return rowsAffected, nil
}

// GiveEveryoneCredit gives each user 1 credit (respecting max credits)
func (r *UserRepository) GiveEveryoneCredit(maxCredits int) (int64, error) {
	result, err := database.DB.Exec(`
		UPDATE users
		SET credits = MIN(credits + 1, ?), updated_at = CURRENT_TIMESTAMP
		WHERE credits < ?`,
		maxCredits, maxCredits)
	if err != nil {
		return 0, fmt.Errorf("failed to give everyone credit: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return rowsAffected, nil
}

// FindOrCreate finds a user by Steam ID or creates a new one
func (r *UserRepository) FindOrCreate(steamID, username, avatarURL, avatarSmall, profileURL string) (*models.User, bool, error) {
	// Try to find existing user
	user, err := r.GetBySteamID(steamID)
	if err != nil {
		return nil, false, err
	}

	if user != nil {
		// Update profile data if it changed
		if user.Username != username || user.AvatarURL != avatarURL {
			user.Username = username
			user.AvatarURL = avatarURL
			user.AvatarSmall = avatarSmall
			user.ProfileURL = profileURL
			if err := r.Update(user); err != nil {
				return nil, false, err
			}
		}
		return user, false, nil // false = existing user
	}

	// Create new user
	user = &models.User{
		SteamID:      steamID,
		Username:     username,
		AvatarURL:    avatarURL,
		AvatarSmall:  avatarSmall,
		ProfileURL:   profileURL,
		Credits:      0,
		LastCreditAt: time.Now(),
	}

	if err := r.Create(user); err != nil {
		return nil, false, err
	}

	return user, true, nil // true = new user created
}
