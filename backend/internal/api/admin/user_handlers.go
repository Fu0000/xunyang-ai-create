package admin

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"google-ai-proxy/internal/db"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type userResponse struct {
	ID              uint64 `json:"id"`
	Email           string `json:"email"`
	Nickname        string `json:"nickname"`
	Avatar          string `json:"avatar"`
	Credits         int    `json:"credits"`
	TotalRedeemed   int    `json:"total_redeemed"`
	UsageCount      int    `json:"usage_count"`
	Status          string `json:"status"`
	EmailVerified   bool   `json:"email_verified"`
	IsLinuxDo       bool   `json:"is_linuxdo"`
	InviteCount     int    `json:"invite_count"`
	CreatedAt       int64  `json:"created_at"`
	LastLoginAt     int64  `json:"last_login_at"`
}

func buildUserResponse(user db.User) userResponse {
	lastLogin := int64(0)
	if user.LastLoginAt != nil && !user.LastLoginAt.IsZero() {
		lastLogin = user.LastLoginAt.UnixMilli()
	}

	return userResponse{
		ID:            user.ID,
		Email:         user.Email,
		Nickname:      user.Nickname,
		Avatar:        user.Avatar,
		Credits:       user.Credits,
		TotalRedeemed: user.TotalRedeemed,
		UsageCount:    user.UsageCount,
		Status:        user.Status,
		EmailVerified: user.EmailVerified,
		IsLinuxDo:     user.LinuxDoID != nil && *user.LinuxDoID != "",
		InviteCount:   user.InviteCount,
		CreatedAt:     user.CreatedAt.UnixMilli(),
		LastLoginAt:   lastLogin,
	}
}

// ListUsers lists users with filtering and pagination.
func ListUsers(c *gin.Context) {
	limit, offset := parseListPagination(c)

	query := db.DB.Model(&db.User{})

	if uidStr := strings.TrimSpace(c.Query("user_id")); uidStr != "" {
		if uid, err := strconv.ParseUint(uidStr, 10, 64); err == nil {
			query = query.Where("id = ?", uid)
		}
	}

	if email := strings.TrimSpace(c.Query("email")); email != "" {
		query = query.Where("email LIKE ?", "%"+email+"%")
	}

	if status := strings.TrimSpace(c.Query("status")); status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to query users"})
		return
	}

	var users []db.User
	if err := query.Order("id DESC").Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to query users"})
		return
	}

	items := make([]userResponse, 0, len(users))
	for _, u := range users {
		items = append(items, buildUserResponse(u))
	}

	c.JSON(http.StatusOK, gin.H{
		"items":  items,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

type updateUserCreditsRequest struct {
	Delta int    `json:"delta" binding:"required"`
	Note  string `json:"note"`
}

// UpdateUserCredits adds or subtracts credits from a user.
func UpdateUserCredits(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	var req updateUserCreditsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if req.Delta == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "delta cannot be 0"})
		return
	}

	operatorSource := c.GetString("adminOperatorSource")
	if operatorSource == "" {
		operatorSource = "admin_console"
	}
	operatorID := c.GetString("adminOperatorID")

	err = db.DB.Transaction(func(tx *gorm.DB) error {
		var user db.User
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&user, userID).Error; err != nil {
			return err
		}

		newBalance := user.Credits + req.Delta
		if newBalance < 0 {
			newBalance = 0
			req.Delta = -user.Credits
		}

		if err := tx.Model(&db.User{}).Where("id = ?", userID).Updates(map[string]interface{}{
			"credits":    newBalance,
			"updated_at": time.Now(),
		}).Error; err != nil {
			return err
		}

		note := req.Note
		if note == "" {
			note = fmt.Sprintf("Admin manual adjustment by %s", operatorID)
		}

		txRecord := db.CreditTransaction{
			UserID:       userID,
			Delta:        req.Delta,
			BalanceAfter: newBalance,
			Type:         "admin_adjustment",
			Source:       operatorSource,
			SourceID:     operatorID,
			Note:         note,
			CreatedAt:    time.Now(),
		}
		if err := tx.Create(&txRecord).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user credits"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

type updateUserStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

// UpdateUserStatus updates user status (e.g. active, banned)
func UpdateUserStatus(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	var req updateUserStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	status := strings.ToLower(strings.TrimSpace(req.Status))
	if status != "active" && status != "banned" && status != "disabled" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status, expected active/banned/disabled"})
		return
	}

	if err := db.DB.Model(&db.User{}).Where("id = ?", userID).Update("status", status).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}
