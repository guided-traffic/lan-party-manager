package services

import (
	"time"

	"github.com/guided-traffic/rate-your-mate/backend/config"
	"github.com/guided-traffic/rate-your-mate/backend/models"
	"github.com/guided-traffic/rate-your-mate/backend/repository"
)

// CreditService handles credit calculation and management
type CreditService struct {
	cfg      *config.Config
	userRepo *repository.UserRepository
}

// NewCreditService creates a new credit service
func NewCreditService(cfg *config.Config, userRepo *repository.UserRepository) *CreditService {
	return &CreditService{
		cfg:      cfg,
		userRepo: userRepo,
	}
}

// CalculateAndUpdateCredits calculates new credits based on time elapsed and updates the user
// Returns the updated credit count
// Note: When voting is paused, no new credits are generated
func (s *CreditService) CalculateAndUpdateCredits(user *models.User) (int, error) {
	// If voting is paused, don't generate new credits
	if s.cfg.VotingPaused {
		return user.Credits, nil
	}

	now := time.Now()

	// Calculate time elapsed since last credit was given
	elapsed := now.Sub(user.LastCreditAt)
	intervalDuration := time.Duration(s.cfg.CreditIntervalMinutes) * time.Minute

	// Calculate how many new credits should be added
	newCredits := int(elapsed / intervalDuration)

	if newCredits <= 0 {
		// No new credits earned yet
		return user.Credits, nil
	}

	// Calculate total credits (capped at max)
	totalCredits := user.Credits + newCredits
	if totalCredits > s.cfg.CreditMax {
		totalCredits = s.cfg.CreditMax
	}

	// Calculate new last_credit_at time
	// We move it forward by the number of intervals used
	creditsActuallyAdded := totalCredits - user.Credits
	if creditsActuallyAdded > 0 || totalCredits == s.cfg.CreditMax {
		// Move last_credit_at forward
		newLastCreditAt := user.LastCreditAt.Add(time.Duration(newCredits) * intervalDuration)

		// Don't set it to the future
		if newLastCreditAt.After(now) {
			newLastCreditAt = now
		}

		// Update in database
		if err := s.userRepo.UpdateCredits(user.ID, totalCredits, newLastCreditAt); err != nil {
			return user.Credits, err
		}

		user.Credits = totalCredits
		user.LastCreditAt = newLastCreditAt
	}

	return totalCredits, nil
}

// GetTimeUntilNextCredit returns the duration until the user earns their next credit
// Returns 0 if the user is at max credits
// Returns -1 if voting is paused (credit generation is disabled)
func (s *CreditService) GetTimeUntilNextCredit(user *models.User) time.Duration {
	// If voting is paused, credit generation is disabled
	if s.cfg.VotingPaused {
		return -1
	}

	if user.Credits >= s.cfg.CreditMax {
		return 0
	}

	intervalDuration := time.Duration(s.cfg.CreditIntervalMinutes) * time.Minute
	nextCreditAt := user.LastCreditAt.Add(intervalDuration)

	remaining := time.Until(nextCreditAt)
	if remaining < 0 {
		return 0
	}

	return remaining
}

// CanAffordVote checks if a user has enough credits to vote
func (s *CreditService) CanAffordVote(user *models.User) bool {
	return user.Credits >= 1
}

// CanAffordVoteWithPoints checks if a user has enough credits for a vote with specific points
func (s *CreditService) CanAffordVoteWithPoints(user *models.User, points int) bool {
	return user.Credits >= points
}

// DeductVoteCost deducts the cost of a vote from the user's credits
func (s *CreditService) DeductVoteCost(userID uint64) error {
	return s.userRepo.DeductCredit(userID)
}

// DeductVoteCostWithPoints deducts multiple credits for a vote with points
func (s *CreditService) DeductVoteCostWithPoints(userID uint64, points int) error {
	return s.userRepo.DeductCredits(userID, points)
}
