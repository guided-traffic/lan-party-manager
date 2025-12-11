package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/guided-traffic/rate-your-mate/backend/models"
)

// AchievementHandler handles achievement-related endpoints
type AchievementHandler struct{}

// NewAchievementHandler creates a new achievement handler
func NewAchievementHandler() *AchievementHandler {
	return &AchievementHandler{}
}

// GetAll returns all available achievements
// GET /api/v1/achievements
func (h *AchievementHandler) GetAll(c *gin.Context) {
	achievements := models.GetAllAchievements()

	// Separate positive and negative achievements
	positive := make([]models.Achievement, 0)
	negative := make([]models.Achievement, 0)

	for _, a := range achievements {
		if a.IsPositive {
			positive = append(positive, a)
		} else {
			negative = append(negative, a)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"achievements": achievements,
		"positive":     positive,
		"negative":     negative,
	})
}

// GetByID returns a single achievement by ID
// GET /api/v1/achievements/:id
func (h *AchievementHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	achievement, ok := models.GetAchievement(id)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Achievement not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"achievement": achievement,
	})
}
