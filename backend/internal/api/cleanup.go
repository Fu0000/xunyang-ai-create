package api

// TASK-24: api_logs 表自动清理

import (
	"log"
	"time"

	"google-ai-proxy/internal/db"
)

// StartAPILogCleanup TASK-24: 每日清理 30 天前的 api_logs 记录，防止表体积无限增长。
// 采用批量删除（每次最多 1000 条）降低锁表风险。
func StartAPILogCleanup() {
	go func() {
		// 服务启动后 10 分钟首次执行，之后每 24 小时执行一次
		time.Sleep(10 * time.Minute)

		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()

		cleanupAPILogs()
		for range ticker.C {
			cleanupAPILogs()
		}
	}()
	log.Println("[APILog清理] 后台清理 Worker 已启动 (每24小时执行，保留最近30天)")
}

// cleanupAPILogs 删除 30 天前的 api_logs，批量删除最多 1000 条/次
func cleanupAPILogs() {
	cutoff := time.Now().AddDate(0, 0, -30)
	totalDeleted := int64(0)

	for {
		result := db.DB.Where("created_at < ?", cutoff).
			Limit(1000).
			Delete(&db.APILog{})

		if result.Error != nil {
			log.Printf("[APILog清理] 删除出错: %v", result.Error)
			break
		}
		totalDeleted += result.RowsAffected
		if result.RowsAffected < 1000 {
			// 已删完本批
			break
		}
		// 短暂休眠，避免持续重锁
		time.Sleep(200 * time.Millisecond)
	}

	if totalDeleted > 0 {
		log.Printf("[APILog清理] 已清理 %d 条超过30天的 API 日志", totalDeleted)
	}
}
