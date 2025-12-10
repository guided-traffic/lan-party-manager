package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/guided-traffic/lan-party-manager/backend/services"
)

// GameHandler handles game-related HTTP requests
type GameHandler struct {
	gameService *services.GameService
}

// NewGameHandler creates a new game handler
func NewGameHandler(gameService *services.GameService) *GameHandler {
	return &GameHandler{
		gameService: gameService,
	}
}

// GetMultiplayerGames returns all multiplayer games owned by players
// GET /api/v1/games
func (h *GameHandler) GetMultiplayerGames(c *gin.Context) {
	games, err := h.gameService.GetMultiplayerGames()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch games",
		})
		return
	}

	c.JSON(http.StatusOK, games)
}

// RefreshGames invalidates the cache and returns fresh game data
// POST /api/v1/games/refresh
func (h *GameHandler) RefreshGames(c *gin.Context) {
	h.gameService.InvalidateCache()

	games, err := h.gameService.GetMultiplayerGames()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to refresh games",
		})
		return
	}

	c.JSON(http.StatusOK, games)
}
