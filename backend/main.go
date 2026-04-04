package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"google-ai-proxy/internal/api"
	adminapi "google-ai-proxy/internal/api/admin"
	"google-ai-proxy/internal/auth"
	"google-ai-proxy/internal/config"
	"google-ai-proxy/internal/db"
	"google-ai-proxy/internal/email"
	_ "google-ai-proxy/internal/provider"
	"google-ai-proxy/internal/storage"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// If APP_ENV is set (e.g. "dev"), load .env.{APP_ENV} first.
	var envPaths []string
	if appEnv := os.Getenv("APP_ENV"); appEnv != "" {
		envFile := ".env." + appEnv
		envPaths = append(envPaths,
			filepath.Join(".", envFile),
			filepath.Join("backend", envFile),
		)
	}
	envPaths = append(envPaths,
		filepath.Join(".", ".env"),
		filepath.Join("backend", ".env"),
		"/opt/nanobanana/.env",
		"/etc/google-ai-proxy/.env",
	)

	envLoaded := false
	for _, envPath := range envPaths {
		if err := godotenv.Load(envPath); err == nil {
			log.Printf("loaded .env: %s", envPath)
			envLoaded = true
			break
		}
	}
	if !envLoaded {
		log.Printf("warning: .env not found in standard locations, using process env")
		log.Printf("debug cwd: %s", getCurrentDir())
	}

	auth.InitSecretKey()
	db.InitDB()
	email.InitEmail()

	if err := storage.InitOSS(); err != nil {
		log.Printf("warning: OSS init failed: %v", err)
	}

	// Background workers.
	api.StartVideoTaskPoller()
	api.StartVerificationCleanup()
	api.StartGenerationCleanup()
	api.StartAPILogCleanup() // TASK-24: 定期清理 30 天前的 API 日志

	r := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsOrigins := config.GetCORSOrigins()
	if corsOrigins != "" {
		corsConfig.AllowOrigins = strings.Split(corsOrigins, ",")
	} else {
		corsConfig.AllowOrigins = []string{"http://localhost:5173", "http://localhost:5174"}
	}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization", "X-Admin-Token"}
	r.Use(cors.New(corsConfig))

	// TASK-11: 全局 IP 限速（100 req/min/IP）
	r.Use(api.GlobalRateLimitMiddleware())

	// TASK-09: 健康检查端点（不需要认证），供负载均衡器和 K8s 存活探针使用
	r.GET("/health", func(c *gin.Context) {
		dbStatus := "ok"
		sqlDB, err := db.DB.DB()
		if err != nil || sqlDB.Ping() != nil {
			dbStatus = "error"
		}
		status := "ok"
		if dbStatus != "ok" {
			status = "degraded"
		}
		c.JSON(http.StatusOK, gin.H{
			"status":    status,
			"db":        dbStatus,
			"timestamp": time.Now().Unix(),
		})
	})

	apiGroup := r.Group("/api")
	{
		// Public
		apiGroup.GET("/pricing", api.GetPricing)
		apiGroup.GET("/models", api.GetModels)


		// Auth - 严格限速（10 req/min/IP），防止验证码爆破 (TASK-11)
		authRoutes := apiGroup.Group("/auth")
		authRoutes.Use(api.AuthRateLimitMiddleware())
		{
			authRoutes.POST("/send-code", api.SendVerificationCode)
			authRoutes.POST("/register", api.Register)
			authRoutes.POST("/login", api.LoginWithEmail)
			authRoutes.POST("/reset-password", api.ResetPassword)
		}

		// OAuth
		apiGroup.GET("/auth/oauth/linuxdo", api.LinuxDoOAuthURL)
		apiGroup.POST("/auth/oauth/linuxdo/callback", api.LinuxDoOAuthCallback)

		// Payment callback (public, no auth required)
		apiGroup.GET("/payment/notify/linuxdo", api.LinuxDoCreditNotify)

		// User
		userGroup := apiGroup.Group("/user")
		userGroup.Use(api.UserAuthMiddleware())
		{
			userGroup.GET("/me", api.GetUserMe)
			userGroup.POST("/redeem", api.RedeemKey)
			userGroup.POST("/daily-checkin", api.DailyCheckin)
			userGroup.GET("/invitations", api.GetInvitationRecords)
			userGroup.GET("/credits/transactions", api.GetCreditTransactions)
			userGroup.GET("/notifications", api.ListUserNotifications)
			userGroup.POST("/notifications/read-all", api.MarkAllNotificationsRead)
			userGroup.POST("/notifications/:id/read", api.MarkNotificationRead)
			userGroup.POST("/bind-email", api.BindEmail)
			userGroup.POST("/upload/image", api.UploadImage)
			userGroup.POST("/upload/video", api.UploadVideo)
			userGroup.PUT("/profile", api.UpdateUserProfile)          // TASK-19: 更新用户资料
			userGroup.POST("/logout", api.Logout)                     // TASK-23: 主动退出登录

			// Payment (authenticated)
			userGroup.POST("/payment/create", api.CreatePaymentOrder)
			userGroup.GET("/payment/status/:orderNo", api.GetPaymentStatus)
			userGroup.GET("/payment/orders", api.GetPaymentOrders)
		}

		// Unified generation
		apiGroup.POST("/generate", api.UserAuthMiddleware(), api.UnifiedGenerate)
		apiGroup.POST("/prompt/optimize", api.UserAuthMiddleware(), api.OptimizePrompt)
		apiGroup.POST("/tools/reverse-prompt", api.UserAuthMiddleware(), api.ReversePrompt)

		// Public inspirations
		apiGroup.GET("/inspirations", api.ListPublicInspirations)
		apiGroup.GET("/inspirations/liked", api.UserAuthMiddleware(), api.ListLikedInspirations)
		apiGroup.GET("/inspirations/mine", api.UserAuthMiddleware(), api.ListMyInspirations)
		apiGroup.GET("/inspirations/:shareID", api.GetPublicInspiration)
		apiGroup.GET("/inspirations/:shareID/liked", api.UserAuthMiddleware(), api.GetInspirationLikeStatus)
		apiGroup.POST("/inspirations/:shareID/like", api.UserAuthMiddleware(), api.LikeInspiration)
		apiGroup.DELETE("/inspirations/:shareID/like", api.UserAuthMiddleware(), api.UnlikeInspiration)
		apiGroup.POST("/inspirations/:shareID/remix", api.UserAuthMiddleware(), api.MarkInspirationRemix)
		apiGroup.DELETE("/inspirations/:shareID", api.UserAuthMiddleware(), api.UnshareInspirationByShareID)
		apiGroup.POST("/inspirations/publish", api.UserAuthMiddleware(), api.PublishInspiration)
		apiGroup.GET("/inspiration-tags", api.ListInspirationTags)

		// Generation history
		generationsGroup := apiGroup.Group("/generations")
		generationsGroup.Use(api.UserAuthMiddleware())
		{
			generationsGroup.GET("", api.ListGenerations)
			generationsGroup.GET("/:id", api.GetGeneration)
			generationsGroup.PUT("/:id", api.UpdateGeneration)
			generationsGroup.POST("/:id/share", api.ShareGeneration)
			generationsGroup.DELETE("/:id/share", api.UnshareGeneration)
			generationsGroup.DELETE("/:id", api.DeleteGeneration)
		}

		// Admin moderation
		adminGroup := apiGroup.Group("/admin")
		adminGroup.Use(adminapi.AuthMiddleware())
		{
			adminGroup.GET("/inspirations", adminapi.ListInspirations)
			adminGroup.POST("/inspirations/:id/review", adminapi.ReviewInspiration)
		}
	}

	port := ":" + config.GetPort()
	log.Printf("server listening on %s", port)
	if err := r.Run(port); err != nil {
		log.Printf("server start failed: %v", err)
	}
}

func getCurrentDir() string {
	dir, err := os.Getwd()
	if err != nil {
		return "unknown"
	}
	return dir
}
