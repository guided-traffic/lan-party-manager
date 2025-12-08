package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/guided-traffic/lan-party-manager/backend/config"
	"github.com/guided-traffic/lan-party-manager/backend/middleware"
	"github.com/guided-traffic/lan-party-manager/backend/repository"
	"github.com/guided-traffic/lan-party-manager/backend/websocket"
)

// SettingsHandler handles admin settings endpoints
type SettingsHandler struct {
	cfg      *config.Config
	wsHub    *websocket.Hub
	userRepo *repository.UserRepository
}

// NewSettingsHandler creates a new settings handler
func NewSettingsHandler(cfg *config.Config, wsHub *websocket.Hub, userRepo *repository.UserRepository) *SettingsHandler {
	return &SettingsHandler{
		cfg:      cfg,
		wsHub:    wsHub,
		userRepo: userRepo,
	}
}

// GetSettingsRequest represents the response for GET /settings
type GetSettingsResponse struct {
	CreditIntervalMinutes int  `json:"credit_interval_minutes"`
	CreditMax             int  `json:"credit_max"`
	VotingPaused          bool `json:"voting_paused"`
}

// UpdateSettingsRequest represents the request body for PUT /settings
type UpdateSettingsRequest struct {
	CreditIntervalMinutes *int  `json:"credit_interval_minutes"`
	CreditMax             *int  `json:"credit_max"`
	VotingPaused          *bool `json:"voting_paused"`
}

// VotingStatusResponse represents the response for GET /voting-status
type VotingStatusResponse struct {
	VotingPaused bool `json:"voting_paused"`
}

// GetVotingStatus returns only the voting paused status (for non-admin users)
// GET /api/v1/voting-status
func (h *SettingsHandler) GetVotingStatus(c *gin.Context) {
	c.JSON(http.StatusOK, VotingStatusResponse{
		VotingPaused: h.cfg.VotingPaused,
	})
}

// GetSettings returns the current settings
// GET /api/v1/admin/settings
func (h *SettingsHandler) GetSettings(c *gin.Context) {
	c.JSON(http.StatusOK, GetSettingsResponse{
		CreditIntervalMinutes: h.cfg.CreditIntervalMinutes,
		CreditMax:             h.cfg.CreditMax,
		VotingPaused:          h.cfg.VotingPaused,
	})
}

// UpdateSettings updates the settings (admin only)
// PUT /api/v1/admin/settings
func (h *SettingsHandler) UpdateSettings(c *gin.Context) {
	var req UpdateSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	// Validate and update settings
	updated := false

	if req.CreditIntervalMinutes != nil {
		if *req.CreditIntervalMinutes < 1 || *req.CreditIntervalMinutes > 60 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "credit_interval_minutes must be between 1 and 60",
			})
			return
		}
		h.cfg.CreditIntervalMinutes = *req.CreditIntervalMinutes
		updated = true
		log.Printf("Admin updated credit_interval_minutes to %d", *req.CreditIntervalMinutes)
	}

	if req.CreditMax != nil {
		if *req.CreditMax < 1 || *req.CreditMax > 100 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "credit_max must be between 1 and 100",
			})
			return
		}
		h.cfg.CreditMax = *req.CreditMax
		updated = true
		log.Printf("Admin updated credit_max to %d", *req.CreditMax)
	}

	if req.VotingPaused != nil {
		h.cfg.VotingPaused = *req.VotingPaused
		updated = true
		if *req.VotingPaused {
			log.Printf("Admin paused voting")
		} else {
			log.Printf("Admin resumed voting")
		}
	}

	// Broadcast settings change to all connected clients
	if updated {
		h.wsHub.BroadcastSettingsUpdate(&websocket.SettingsPayload{
			CreditIntervalMinutes: h.cfg.CreditIntervalMinutes,
			CreditMax:             h.cfg.CreditMax,
			VotingPaused:          h.cfg.VotingPaused,
		})
	}

	c.JSON(http.StatusOK, GetSettingsResponse{
		CreditIntervalMinutes: h.cfg.CreditIntervalMinutes,
		CreditMax:             h.cfg.CreditMax,
		VotingPaused:          h.cfg.VotingPaused,
	})
}

// ResetAllCreditsResponse represents the response for POST /admin/credits/reset
type ResetAllCreditsResponse struct {
	Message       string `json:"message"`
	UsersAffected int64  `json:"users_affected"`
}

// GiveEveryoneCreditResponse represents the response for POST /admin/credits/give
type GiveEveryoneCreditResponse struct {
	Message       string `json:"message"`
	UsersAffected int64  `json:"users_affected"`
}

// ResetAllCredits sets all users' credits to 0
// POST /api/v1/admin/credits/reset
func (h *SettingsHandler) ResetAllCredits(c *gin.Context) {
	usersAffected, err := h.userRepo.ResetAllCredits()
	if err != nil {
		log.Printf("Error resetting all credits: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to reset credits",
		})
		return
	}

	log.Printf("Admin reset all credits - %d users affected", usersAffected)

	// Broadcast credit reset to all connected clients
	h.wsHub.BroadcastCreditsReset()

	c.JSON(http.StatusOK, ResetAllCreditsResponse{
		Message:       "Alle Credits wurden auf 0 gesetzt",
		UsersAffected: usersAffected,
	})
}

// GiveEveryoneCredit gives each user 1 credit
// POST /api/v1/admin/credits/give
func (h *SettingsHandler) GiveEveryoneCredit(c *gin.Context) {
	usersAffected, err := h.userRepo.GiveEveryoneCredit(h.cfg.CreditMax)
	if err != nil {
		log.Printf("Error giving everyone credit: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to give credits",
		})
		return
	}

	log.Printf("Admin gave everyone a credit - %d users affected", usersAffected)

	// Broadcast credit update to all connected clients
	h.wsHub.BroadcastCreditsGiven()

	c.JSON(http.StatusOK, GiveEveryoneCreditResponse{
		Message:       "Jedem Spieler wurde 1 Credit gegeben",
		UsersAffected: usersAffected,
	})
}

// AdminMiddleware checks if the current user is an admin
func (h *SettingsHandler) AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, ok := middleware.GetClaims(c)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Not authenticated",
			})
			c.Abort()
			return
		}

		if !h.cfg.IsAdmin(claims.SteamID) {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Admin access required",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
