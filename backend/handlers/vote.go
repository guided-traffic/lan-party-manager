package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/guided-traffic/lan-party-manager/backend/middleware"
	"github.com/guided-traffic/lan-party-manager/backend/models"
	"github.com/guided-traffic/lan-party-manager/backend/repository"
	"github.com/guided-traffic/lan-party-manager/backend/services"
	"github.com/guided-traffic/lan-party-manager/backend/websocket"
)

// VoteHandler handles vote-related endpoints
type VoteHandler struct {
	voteRepo      *repository.VoteRepository
	userRepo      *repository.UserRepository
	creditService *services.CreditService
	wsHub         *websocket.Hub
}

// NewVoteHandler creates a new vote handler
func NewVoteHandler(voteRepo *repository.VoteRepository, userRepo *repository.UserRepository, creditService *services.CreditService, wsHub *websocket.Hub) *VoteHandler {
	return &VoteHandler{
		voteRepo:      voteRepo,
		userRepo:      userRepo,
		creditService: creditService,
		wsHub:         wsHub,
	}
}

// Create creates a new vote
// POST /api/v1/votes
func (h *VoteHandler) Create(c *gin.Context) {
	// Get current user
	fromUserID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Not authenticated",
		})
		return
	}

	// Parse request body
	var req models.CreateVoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	// Validate achievement
	if !models.IsValidAchievement(req.AchievementID) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid achievement ID",
		})
		return
	}

	// Can't vote for yourself
	if fromUserID == req.ToUserID {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Cannot vote for yourself",
		})
		return
	}

	// Check if target user exists
	toUser, err := h.userRepo.GetByID(req.ToUserID)
	if err != nil {
		log.Printf("Failed to check target user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to process vote",
		})
		return
	}
	if toUser == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Target user not found",
		})
		return
	}

	// Check and update credits for current user
	fromUser, err := h.userRepo.GetByID(fromUserID)
	if err != nil {
		log.Printf("Failed to load current user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to process vote",
		})
		return
	}

	// Calculate current credits
	_, err = h.creditService.CalculateAndUpdateCredits(fromUser)
	if err != nil {
		log.Printf("Failed to calculate credits: %v", err)
	}

	// Reload user to get updated credits
	fromUser, _ = h.userRepo.GetByID(fromUserID)

	// Check if user has enough credits
	if !h.creditService.CanAffordVote(fromUser) {
		c.JSON(http.StatusPaymentRequired, gin.H{
			"error":   "Insufficient credits",
			"credits": fromUser.Credits,
		})
		return
	}

	// Deduct credit
	if err := h.creditService.DeductVoteCost(fromUserID); err != nil {
		log.Printf("Failed to deduct credit: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to process vote",
		})
		return
	}

	// Create the vote
	vote := &models.Vote{
		FromUserID:    fromUserID,
		ToUserID:      req.ToUserID,
		AchievementID: req.AchievementID,
	}

	if err := h.voteRepo.Create(vote); err != nil {
		log.Printf("Failed to create vote: %v", err)
		// Try to refund the credit
		// (In a real app, this should be a transaction)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create vote",
		})
		return
	}

	// Get full vote details for response
	voteDetails, err := h.voteRepo.GetByID(vote.ID)
	if err != nil {
		log.Printf("Failed to get vote details: %v", err)
	}

	// Broadcast vote to WebSocket clients
	if voteDetails != nil && h.wsHub != nil {
		achievement, _ := models.GetAchievement(voteDetails.AchievementID)
		payload := &websocket.VotePayload{
			VoteID:        voteDetails.ID,
			FromUserID:    voteDetails.FromUser.ID,
			FromUsername:  voteDetails.FromUser.Username,
			FromAvatar:    voteDetails.FromUser.AvatarSmall,
			ToUserID:      voteDetails.ToUser.ID,
			ToUsername:    voteDetails.ToUser.Username,
			ToAvatar:      voteDetails.ToUser.AvatarSmall,
			AchievementID: voteDetails.AchievementID,
			Achievement:   achievement.Name,
			IsPositive:    achievement.IsPositive,
			CreatedAt:     voteDetails.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}

		// Broadcast to all clients (for timeline updates)
		h.wsHub.BroadcastVote(payload)

		// Send personal notification to the recipient
		h.wsHub.NotifyVoteReceived(toUser.ID, payload)
	}

	// Return updated credits
	fromUser, _ = h.userRepo.GetByID(fromUserID)

	c.JSON(http.StatusCreated, gin.H{
		"vote":    voteDetails,
		"credits": fromUser.Credits,
	})
}

// GetTimeline returns recent votes for the timeline
// GET /api/v1/votes
func (h *VoteHandler) GetTimeline(c *gin.Context) {
	votes, err := h.voteRepo.GetRecent(100)
	if err != nil {
		log.Printf("Failed to get timeline: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to load timeline",
		})
		return
	}

	if votes == nil {
		votes = []models.VoteWithDetails{}
	}

	c.JSON(http.StatusOK, gin.H{
		"votes": votes,
	})
}

// GetLeaderboard returns the leaderboard (top 3 per achievement)
// GET /api/v1/leaderboard
func (h *VoteHandler) GetLeaderboard(c *gin.Context) {
	leaderboard, err := h.voteRepo.GetLeaderboard(3)
	if err != nil {
		log.Printf("Failed to get leaderboard: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to load leaderboard",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"leaderboard": leaderboard,
	})
}
