package admin

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"google-ai-proxy/internal/db"

	"github.com/gin-gonic/gin"
)

type generationResponse struct {
	ID              uint64   `json:"id"`
	UserID          uint64   `json:"user_id"`
	Type            string   `json:"type"`
	Prompt          string   `json:"prompt"`
	ReferenceImages []string `json:"reference_images"`
	Params          string   `json:"params"`
	Images          []string `json:"images"`
	VideoURL        string   `json:"video_url"`
	Status          string   `json:"status"`
	CreditsCost     int      `json:"credits_cost"`
	ErrorMsg        string   `json:"error_msg"`
	TaskID          string   `json:"task_id"`
	IsFavorite      bool     `json:"is_favorite"`
	CreatedAt       int64    `json:"created_at"`
	UpdatedAt       int64    `json:"updated_at"`
	UserEmail       string   `json:"user_email"`
}

func buildGenerationResponse(g db.Generation, u db.User) generationResponse {
	var refImages []string
	if g.ReferenceImages != "" {
		_ = json.Unmarshal([]byte(g.ReferenceImages), &refImages)
	}
	var outImages []string
	if g.Images != "" {
		_ = json.Unmarshal([]byte(g.Images), &outImages)
	}
	taskID := ""
	if g.TaskID != nil {
		taskID = *g.TaskID
	}
	return generationResponse{
		ID:              g.ID,
		UserID:          g.UserID,
		Type:            g.Type,
		Prompt:          g.Prompt,
		ReferenceImages: refImages,
		Params:          g.Params,
		Images:          outImages,
		VideoURL:        g.VideoURL,
		Status:          g.Status,
		CreditsCost:     g.CreditsCost,
		ErrorMsg:        g.ErrorMsg,
		TaskID:          taskID,
		IsFavorite:      g.IsFavorite,
		CreatedAt:       g.CreatedAt.UnixMilli(),
		UpdatedAt:       g.UpdatedAt.UnixMilli(),
		UserEmail:       u.Email,
	}
}

// ListGenerations lists generations across all users with filtering.
func ListGenerations(c *gin.Context) {
	limit, offset := parseListPagination(c)

	query := db.DB.Model(&db.Generation{})

	if uidStr := strings.TrimSpace(c.Query("user_id")); uidStr != "" {
		if uid, err := strconv.ParseUint(uidStr, 10, 64); err == nil {
			query = query.Where("user_id = ?", uid)
		}
	}

	if gType := strings.TrimSpace(c.Query("type")); gType != "" {
		query = query.Where("type = ?", gType)
	}

	if status := strings.TrimSpace(c.Query("status")); status != "" {
		query = query.Where("status = ?", status)
	}

	if taskID := strings.TrimSpace(c.Query("task_id")); taskID != "" {
		query = query.Where("task_id LIKE ?", "%"+taskID+"%")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to query generations"})
		return
	}

	var generations []db.Generation
	if err := query.Order("id DESC").Limit(limit).Offset(offset).Find(&generations).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to query generations"})
		return
	}

	// Fetch users
	userIDs := make([]uint64, 0, len(generations))
	for _, g := range generations {
		userIDs = append(userIDs, g.UserID)
	}

	userMap := make(map[uint64]db.User)
	if len(userIDs) > 0 {
		var users []db.User
		if err := db.DB.Select("id", "email").Where("id IN ?", userIDs).Find(&users).Error; err == nil {
			for _, u := range users {
				userMap[u.ID] = u
			}
		}
	}

	items := make([]generationResponse, 0, len(generations))
	for _, g := range generations {
		u := userMap[g.UserID]
		if u.Email == "" {
			u.Email = "unknown"
		}
		items = append(items, buildGenerationResponse(g, u))
	}

	c.JSON(http.StatusOK, gin.H{
		"items":  items,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}
