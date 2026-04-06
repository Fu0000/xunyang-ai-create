package admin

import (
	"log"
	"net/http"
	"strings"

	"google-ai-proxy/internal/db"

	"github.com/gin-gonic/gin"
)

type SettingResponse struct {
	Key         string `json:"key"`
	Value       string `json:"value"`
	Description string `json:"description"`
}

func GetSettings(c *gin.Context) {
	var settings []db.SystemSetting
	if err := db.DB.Find(&settings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to query settings"})
		return
	}

	res := make([]SettingResponse, 0, len(settings))
	for _, s := range settings {
		res = append(res, SettingResponse{
			Key:         s.Key,
			Value:       s.Value,
			Description: s.Description,
		})
	}

	c.JSON(http.StatusOK, gin.H{"settings": res})
}

type UpdateSettingRequest struct {
	Key   string `json:"key" binding:"required"`
	Value string `json:"value" binding:"required"`
}

func UpdateSetting(c *gin.Context) {
	var req UpdateSettingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request format"})
		return
	}

	key := strings.TrimSpace(req.Key)
	val := strings.TrimSpace(req.Value)

	if key == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "key cannot be empty"})
		return
	}

	var setting db.SystemSetting
	if err := db.DB.Where("`key` = ?", key).First(&setting).Error; err != nil {
		// Does not exist, we can optionally insert it.
		setting = db.SystemSetting{
			Key:         key,
			Value:       val,
			Description: "Added via Admin API",
		}
		if err := db.DB.Create(&setting).Error; err != nil {
			log.Printf("Create setting failed: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create setting"})
			return
		}
	} else {
		// Exists, update it
		setting.Value = val
		if err := db.DB.Save(&setting).Error; err != nil {
			log.Printf("Update setting failed: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update setting"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Setting updated successfully"})
}
