package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/guided-traffic/rate-your-mate/backend/auth"
	"github.com/guided-traffic/rate-your-mate/backend/websocket"
)

// WebSocketHandler handles WebSocket connections
type WebSocketHandler struct {
	hub        *websocket.Hub
	jwtService *auth.JWTService
}

// NewWebSocketHandler creates a new WebSocket handler
func NewWebSocketHandler(hub *websocket.Hub, jwtService *auth.JWTService) *WebSocketHandler {
	return &WebSocketHandler{
		hub:        hub,
		jwtService: jwtService,
	}
}

// HandleConnection handles WebSocket connection requests
// The token is passed as a query parameter since WebSocket doesn't support headers easily
// GET /api/v1/ws?token=xxx
func (h *WebSocketHandler) HandleConnection(c *gin.Context) {
	// Get token from query parameter
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Token required",
		})
		return
	}

	// Validate token
	claims, err := h.jwtService.ValidateToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid token",
		})
		return
	}

	// Upgrade to WebSocket
	websocket.ServeWs(h.hub, c.Writer, c.Request, claims.UserID, claims.SteamID, claims.Username)
}

// GetStatus returns WebSocket hub status
// GET /api/v1/ws/status
func (h *WebSocketHandler) GetStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"connected_users": h.hub.GetConnectedUserCount(),
	})
}
