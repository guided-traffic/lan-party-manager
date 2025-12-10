package handlers

import (
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/guided-traffic/lan-party-manager/backend/auth"
	"github.com/guided-traffic/lan-party-manager/backend/config"
	"github.com/guided-traffic/lan-party-manager/backend/repository"
	"github.com/guided-traffic/lan-party-manager/backend/services"
)

// GameHandler handles game-related HTTP requests
type GameHandler struct {
	gameService       *services.GameService
	imageCacheService *services.ImageCacheService
	gameCacheRepo     *repository.GameCacheRepository
	cfg               *config.Config
}

// NewGameHandler creates a new game handler
func NewGameHandler(gameService *services.GameService, imageCacheService *services.ImageCacheService, gameCacheRepo *repository.GameCacheRepository, cfg *config.Config) *GameHandler {
	return &GameHandler{
		gameService:       gameService,
		imageCacheService: imageCacheService,
		gameCacheRepo:     gameCacheRepo,
		cfg:               cfg,
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
