package reports

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetPlatformStats(c *gin.Context) {
	startDate := c.Query("start")
	endDate := c.Query("end")

	// Validate date parameters
	if startDate == "" || endDate == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "start and end date parameters are required (format: YYYY-MM-DD)",
		})
		return
	}

	// Validate date format
	if _, err := time.Parse("2006-01-02", startDate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid start date format. Use YYYY-MM-DD",
		})
		return
	}

	if _, err := time.Parse("2006-01-02", endDate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid end date format. Use YYYY-MM-DD",
		})
		return
	}

	stats, err := h.service.GetPlatformStats(startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve platform statistics",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  stats,
		"count": len(stats),
		"query": gin.H{
			"startDate": startDate,
			"endDate":   endDate,
		},
	})
}

func (h *Handler) GetContentHealth(c *gin.Context) {
	platform := c.Query("platform")
	startDate := c.Query("start")
	endDate := c.Query("end")

	// Validate required parameters
	if platform == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "platform parameter is required (CTV or Audio)",
		})
		return
	}

	if startDate == "" || endDate == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "start and end date parameters are required (format: YYYY-MM-DD)",
		})
		return
	}

	// Validate platform
	if platform != "CTV" && platform != "Audio" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "platform must be either 'CTV' or 'Audio'",
		})
		return
	}

	// Validate date format
	if _, err := time.Parse("2006-01-02", startDate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid start date format. Use YYYY-MM-DD",
		})
		return
	}

	if _, err := time.Parse("2006-01-02", endDate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid end date format. Use YYYY-MM-DD",
		})
		return
	}

	health, err := h.service.GetContentHealth(platform, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve content health data",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  health,
		"count": len(health),
		"query": gin.H{
			"platform":  platform,
			"startDate": startDate,
			"endDate":   endDate,
		},
	})
}

func (h *Handler) GetVideoHealth(c *gin.Context) {
	platform := c.Query("platform")
	startDate := c.Query("start")
	endDate := c.Query("end")

	// Validate required parameters
	if platform == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "platform parameter is required (CTV, Display, or App)",
		})
		return
	}

	if startDate == "" || endDate == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "start and end date parameters are required (format: YYYY-MM-DD)",
		})
		return
	}

	// Validate platform
	validPlatforms := map[string]bool{
		"CTV":     true,
		"Display": true,
		"App":     true,
	}

	if !validPlatforms[platform] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "platform must be one of: CTV, Display, App",
		})
		return
	}

	// Validate date format
	if _, err := time.Parse("2006-01-02", startDate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid start date format. Use YYYY-MM-DD",
		})
		return
	}

	if _, err := time.Parse("2006-01-02", endDate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid end date format. Use YYYY-MM-DD",
		})
		return
	}

	health, err := h.service.GetVideoHealth(platform, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve video health data",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  health,
		"count": len(health),
		"query": gin.H{
			"platform":  platform,
			"startDate": startDate,
			"endDate":   endDate,
		},
	})
}

func (h *Handler) GetDashboard(c *gin.Context) {
	summary, err := h.service.GetDashboardSummary()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve dashboard data",
		})
		return
	}

	c.JSON(http.StatusOK, summary)
}