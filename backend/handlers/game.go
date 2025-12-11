package handlers

import (
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/guided-traffic/rate-your-mate/backend/auth"
	"github.com/guided-traffic/rate-your-mate/backend/config"
	"github.com/guided-traffic/rate-your-mate/backend/repository"
	"github.com/guided-traffic/rate-your-mate/backend/services"
	"github.com/guided-traffic/rate-your-mate/backend/websocket"
)

// GameHandler handles game-related HTTP requests
type GameHandler struct {
	gameService       *services.GameService
	imageCacheService *services.ImageCacheService
	gameCacheRepo     *repository.GameCacheRepository
	cfg               *config.Config
	wsHub             *websocket.Hub
}

// NewGameHandler creates a new game handler
func NewGameHandler(gameService *services.GameService, imageCacheService *services.ImageCacheService, gameCacheRepo *repository.GameCacheRepository, cfg *config.Config, wsHub *websocket.Hub) *GameHandler {
	return &GameHandler{
		gameService:       gameService,
		imageCacheService: imageCacheService,
		gameCacheRepo:     gameCacheRepo,
		cfg:               cfg,
		wsHub:             wsHub,
	}
}

// GetMultiplayerGames returns all multiplayer games owned by players
// GET /api/v1/games
func (h *GameHandler) GetMultiplayerGames(c *gin.Context) {
	// First, return cached data immediately
	games, needsSync, err := h.gameService.GetMultiplayerGamesCached()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch games",
		})
		return
	}

	// Check current sync status
	isSyncing, phase, currentGame, processed, total := h.gameService.GetSyncStatus()

	// Return response with sync status
	c.JSON(http.StatusOK, gin.H{
		"pinned_games": games.PinnedGames,
		"all_games":    games.AllGames,
		"sync_status": gin.H{
			"needs_sync":   needsSync && !isSyncing,
			"is_syncing":   isSyncing,
			"phase":        phase,
			"current_game": currentGame,
			"processed":    processed,
			"total":        total,
		},
	})
}

// StartBackgroundSync triggers a background sync for game data
// POST /api/v1/games/sync
func (h *GameHandler) StartBackgroundSync(c *gin.Context) {
	if h.gameService.IsSyncing() {
		c.JSON(http.StatusConflict, gin.H{
			"message": "Sync already in progress",
		})
		return
	}

	// Start sync with WebSocket progress updates
	h.gameService.SyncGames(func(phase string, currentGame string, processed, total int) {
		percentage := 0
		if total > 0 {
			percentage = (processed * 100) / total
		}

		if phase == "complete" {
			h.wsHub.BroadcastGamesSyncComplete(processed)
		} else {
			h.wsHub.BroadcastGamesSyncProgress(&websocket.GamesSyncProgressPayload{
				Phase:          phase,
				CurrentGame:    currentGame,
				ProcessedCount: processed,
				TotalCount:     total,
				Percentage:     percentage,
			})
		}
	})

	c.JSON(http.StatusAccepted, gin.H{
		"message": "Background sync started",
	})
}

// GetSyncStatus returns the current sync status
// GET /api/v1/games/sync/status
func (h *GameHandler) GetSyncStatus(c *gin.Context) {
	isSyncing, phase, currentGame, processed, total := h.gameService.GetSyncStatus()

	percentage := 0
	if total > 0 {
		percentage = (processed * 100) / total
	}

	c.JSON(http.StatusOK, gin.H{
		"is_syncing":   isSyncing,
		"phase":        phase,
		"current_game": currentGame,
		"processed":    processed,
		"total":        total,
		"percentage":   percentage,
	})
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

// InvalidateDBCache invalidates the database cache, forcing a re-fetch from Steam
// POST /api/v1/admin/games/invalidate-cache
func (h *GameHandler) InvalidateDBCache(c *gin.Context) {
	// Check admin permission
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	jwtClaims := claims.(*auth.Claims)
	if !h.cfg.IsAdmin(jwtClaims.SteamID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		return
	}

	// Invalidate DB cache
	if err := h.gameCacheRepo.InvalidateAll(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to invalidate cache",
		})
		return
	}

	// Also invalidate in-memory cache
	h.gameService.InvalidateCache()

	c.JSON(http.StatusOK, gin.H{
		"message": "Game cache invalidated. Games will be re-fetched from Steam on next request.",
	})
}

// ServeGameImage serves a cached game image
// GET /api/v1/games/images/:filename
func (h *GameHandler) ServeGameImage(c *gin.Context) {
	filename := c.Param("filename")

	// Validate filename format (must be <appid>.jpg)
	if !strings.HasSuffix(filename, ".jpg") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid image format"})
		return
	}

	// Extract app ID from filename
	appIDStr := strings.TrimSuffix(filename, ".jpg")
	appID, err := strconv.Atoi(appIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid app ID"})
		return
	}

	// Check if image exists locally
	imagePath := h.imageCacheService.GetImagePath(appID)

	// If not cached, try to cache it now
	if !h.imageCacheService.HasImage(appID) {
		if !h.imageCacheService.CacheImage(appID) {
			// Redirect to Steam CDN as fallback
			c.Redirect(http.StatusTemporaryRedirect, h.imageCacheService.GetSteamImageURL(appID))
			return
		}
	}

	// Serve the cached image
	c.Header("Cache-Control", "public, max-age=86400") // Cache for 24 hours
	c.File(filepath.Clean(imagePath))
}
