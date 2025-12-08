package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/guided-traffic/lan-party-manager/backend/config"
	"github.com/guided-traffic/lan-party-manager/backend/middleware"
	"github.com/guided-traffic/lan-party-manager/backend/websocket"
)

// SettingsHandler handles admin settings endpoints
type SettingsHandler struct {
	cfg   *config.Config
	wsHub *websocket.Hub
}

// NewSettingsHandler creates a new settings handler
func NewSettingsHandler(cfg *config.Config, wsHub *websocket.Hub) *SettingsHandler {
	return &SettingsHandler{
		cfg:   cfg,
		wsHub: wsHub,
	}
}

// GetSettingsRequest represents the response for GET /settings
type GetSettingsResponse struct {
	CreditIntervalMinutes int `json:"credit_interval_minutes"`
	CreditMax             int `json:"credit_max"`
}

// UpdateSettingsRequest represents the request body for PUT /settings
type UpdateSettingsRequest struct {
	CreditIntervalMinutes *int `json:"credit_interval_minutes"`
	CreditMax             *int `json:"credit_max"`
}

// GetSettings returns the current settings
// GET /api/v1/admin/settings
func (h *SettingsHandler) GetSettings(c *gin.Context) {
	c.JSON(http.StatusOK, GetSettingsResponse{
		CreditIntervalMinutes: h.cfg.CreditIntervalMinutes,
		CreditMax:             h.cfg.CreditMax,
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

	// Broadcast settings change to all connected clients
	if updated {
		h.wsHub.BroadcastSettingsUpdate(&websocket.SettingsPayload{
			CreditIntervalMinutes: h.cfg.CreditIntervalMinutes,
			CreditMax:             h.cfg.CreditMax,
		})
	}

	c.JSON(http.StatusOK, GetSettingsResponse{
		CreditIntervalMinutes: h.cfg.CreditIntervalMinutes,
		CreditMax:             h.cfg.CreditMax,
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
