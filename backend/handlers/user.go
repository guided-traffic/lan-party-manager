package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/guided-traffic/rate-your-mate/backend/middleware"
	"github.com/guided-traffic/rate-your-mate/backend/repository"
)

// UserHandler handles user-related endpoints
type UserHandler struct {
	userRepo *repository.UserRepository
}

// NewUserHandler creates a new user handler
func NewUserHandler(userRepo *repository.UserRepository) *UserHandler {
	return &UserHandler{
		userRepo: userRepo,
	}
}

// GetAll returns all registered users
// GET /api/v1/users
func (h *UserHandler) GetAll(c *gin.Context) {
	users, err := h.userRepo.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to load users",
		})
		return
	}

	// Convert to public user data
	publicUsers := make([]gin.H, len(users))
	for i, user := range users {
		publicUsers[i] = gin.H{
			"id":           user.ID,
			"steam_id":     user.SteamID,
			"username":     user.Username,
			"avatar_url":   user.AvatarURL,
			"avatar_small": user.AvatarSmall,
			"profile_url":  user.ProfileURL,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"users": publicUsers,
	})
}

// GetByID returns a single user by ID
// GET /api/v1/users/:id
func (h *UserHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	user, err := h.userRepo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to load user",
		})
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":           user.ID,
			"steam_id":     user.SteamID,
			"username":     user.Username,
			"avatar_url":   user.AvatarURL,
			"avatar_small": user.AvatarSmall,
			"profile_url":  user.ProfileURL,
		},
	})
}

// GetOthers returns all users except the current user (for voting)
// GET /api/v1/users/others
func (h *UserHandler) GetOthers(c *gin.Context) {
	currentUserID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Not authenticated",
		})
		return
	}

	users, err := h.userRepo.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to load users",
		})
		return
	}

	// Filter out current user
	publicUsers := make([]gin.H, 0)
	for _, user := range users {
		if user.ID != currentUserID {
			publicUsers = append(publicUsers, gin.H{
				"id":           user.ID,
				"steam_id":     user.SteamID,
				"username":     user.Username,
				"avatar_url":   user.AvatarURL,
				"avatar_small": user.AvatarSmall,
				"profile_url":  user.ProfileURL,
			})
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"users": publicUsers,
	})
}
