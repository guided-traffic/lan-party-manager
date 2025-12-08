package repository

import (
	"database/sql"
	"fmt"

	"github.com/guided-traffic/lan-party-manager/backend/database"
	"github.com/guided-traffic/lan-party-manager/backend/models"
)

// VoteRepository handles vote database operations
type VoteRepository struct{}

// NewVoteRepository creates a new vote repository
func NewVoteRepository() *VoteRepository {
	return &VoteRepository{}
}

// Create creates a new vote
func (r *VoteRepository) Create(vote *models.Vote) error {
	result, err := database.DB.Exec(`
		INSERT INTO votes (from_user_id, to_user_id, achievement_id)
		VALUES (?, ?, ?)`,
		vote.FromUserID, vote.ToUserID, vote.AchievementID,
	)
	if err != nil {
		return fmt.Errorf("failed to create vote: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	vote.ID = uint64(id)
	return nil
}

// GetRecent returns the most recent votes for the timeline
func (r *VoteRepository) GetRecent(limit int) ([]models.VoteWithDetails, error) {
	rows, err := database.DB.Query(`
		SELECT
			v.id, v.achievement_id, v.created_at,
			fu.id, fu.steam_id, fu.username, fu.avatar_url, fu.avatar_small, fu.profile_url,
			tu.id, tu.steam_id, tu.username, tu.avatar_url, tu.avatar_small, tu.profile_url
		FROM votes v
		JOIN users fu ON v.from_user_id = fu.id
		JOIN users tu ON v.to_user_id = tu.id
		ORDER BY v.created_at DESC
		LIMIT ?`, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent votes: %w", err)
	}
	defer rows.Close()

	var votes []models.VoteWithDetails
	for rows.Next() {
		var v models.VoteWithDetails
		err := rows.Scan(
			&v.ID, &v.AchievementID, &v.CreatedAt,
			&v.FromUser.ID, &v.FromUser.SteamID, &v.FromUser.Username, &v.FromUser.AvatarURL, &v.FromUser.AvatarSmall, &v.FromUser.ProfileURL,
			&v.ToUser.ID, &v.ToUser.SteamID, &v.ToUser.Username, &v.ToUser.AvatarURL, &v.ToUser.AvatarSmall, &v.ToUser.ProfileURL,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan vote row: %w", err)
		}

		// Add achievement details
		if achievement, ok := models.GetAchievement(v.AchievementID); ok {
			v.Achievement = achievement
		}

		votes = append(votes, v)
	}

	return votes, nil
}

// GetByID returns a vote by ID with full details
func (r *VoteRepository) GetByID(id uint64) (*models.VoteWithDetails, error) {
	var v models.VoteWithDetails
	err := database.DB.QueryRow(`
		SELECT
			v.id, v.achievement_id, v.created_at,
			fu.id, fu.steam_id, fu.username, fu.avatar_url, fu.avatar_small, fu.profile_url,
			tu.id, tu.steam_id, tu.username, tu.avatar_url, tu.avatar_small, tu.profile_url
		FROM votes v
		JOIN users fu ON v.from_user_id = fu.id
		JOIN users tu ON v.to_user_id = tu.id
		WHERE v.id = ?`, id,
	).Scan(
		&v.ID, &v.AchievementID, &v.CreatedAt,
		&v.FromUser.ID, &v.FromUser.SteamID, &v.FromUser.Username, &v.FromUser.AvatarURL, &v.FromUser.AvatarSmall, &v.FromUser.ProfileURL,
		&v.ToUser.ID, &v.ToUser.SteamID, &v.ToUser.Username, &v.ToUser.AvatarURL, &v.ToUser.AvatarSmall, &v.ToUser.ProfileURL,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get vote by id: %w", err)
	}

	// Add achievement details
	if achievement, ok := models.GetAchievement(v.AchievementID); ok {
		v.Achievement = achievement
	}

	return &v, nil
}

// LeaderboardEntry represents a user's position on the leaderboard for an achievement
type LeaderboardEntry struct {
	User       models.PublicUser `json:"user"`
	VoteCount  int               `json:"vote_count"`
	Rank       int               `json:"rank"`
}

// AchievementLeaderboard represents the leaderboard for a single achievement
type AchievementLeaderboard struct {
	Achievement models.Achievement `json:"achievement"`
	Leaders     []LeaderboardEntry `json:"leaders"`
}

// GetLeaderboard returns the top N users per achievement
func (r *VoteRepository) GetLeaderboard(topN int) ([]AchievementLeaderboard, error) {
	// Get all achievements and their top voters
	rows, err := database.DB.Query(`
		SELECT
			v.achievement_id,
			u.id, u.steam_id, u.username, u.avatar_url, u.avatar_small, u.profile_url,
			COUNT(*) as vote_count
		FROM votes v
		JOIN users u ON v.to_user_id = u.id
		GROUP BY v.achievement_id, v.to_user_id
		ORDER BY v.achievement_id, vote_count DESC`)
	if err != nil {
		return nil, fmt.Errorf("failed to get leaderboard: %w", err)
	}
	defer rows.Close()

	// Group by achievement
	achievementMap := make(map[string][]LeaderboardEntry)
	for rows.Next() {
		var achievementID string
		var user models.PublicUser
		var voteCount int

		err := rows.Scan(
			&achievementID,
			&user.ID, &user.SteamID, &user.Username, &user.AvatarURL, &user.AvatarSmall, &user.ProfileURL,
			&voteCount,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan leaderboard row: %w", err)
		}

		// Only keep top N per achievement
		if len(achievementMap[achievementID]) < topN {
			entry := LeaderboardEntry{
				User:      user,
				VoteCount: voteCount,
				Rank:      len(achievementMap[achievementID]) + 1,
			}
			achievementMap[achievementID] = append(achievementMap[achievementID], entry)
		}
	}

	// Build result with all achievements (even those with no votes)
	var result []AchievementLeaderboard
	for _, achievement := range models.GetAllAchievements() {
		lb := AchievementLeaderboard{
			Achievement: achievement,
			Leaders:     achievementMap[achievement.ID],
		}
		if lb.Leaders == nil {
			lb.Leaders = []LeaderboardEntry{}
		}
		result = append(result, lb)
	}

	return result, nil
}

// GetVotesForUser returns all votes received by a user
func (r *VoteRepository) GetVotesForUser(userID uint64) ([]models.VoteWithDetails, error) {
	rows, err := database.DB.Query(`
		SELECT
			v.id, v.achievement_id, v.created_at,
			fu.id, fu.steam_id, fu.username, fu.avatar_url, fu.avatar_small, fu.profile_url,
			tu.id, tu.steam_id, tu.username, tu.avatar_url, tu.avatar_small, tu.profile_url
		FROM votes v
		JOIN users fu ON v.from_user_id = fu.id
		JOIN users tu ON v.to_user_id = tu.id
		WHERE v.to_user_id = ?
		ORDER BY v.created_at DESC`, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get votes for user: %w", err)
	}
	defer rows.Close()

	var votes []models.VoteWithDetails
	for rows.Next() {
		var v models.VoteWithDetails
		err := rows.Scan(
			&v.ID, &v.AchievementID, &v.CreatedAt,
			&v.FromUser.ID, &v.FromUser.SteamID, &v.FromUser.Username, &v.FromUser.AvatarURL, &v.FromUser.AvatarSmall, &v.FromUser.ProfileURL,
			&v.ToUser.ID, &v.ToUser.SteamID, &v.ToUser.Username, &v.ToUser.AvatarURL, &v.ToUser.AvatarSmall, &v.ToUser.ProfileURL,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan vote row: %w", err)
		}

		if achievement, ok := models.GetAchievement(v.AchievementID); ok {
			v.Achievement = achievement
		}

		votes = append(votes, v)
	}

	return votes, nil
}
