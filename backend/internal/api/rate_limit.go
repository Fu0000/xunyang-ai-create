package api

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// ipLimiter holds a rate.Limiter per IP and a last-seen timestamp for cleanup.
type ipLimiter struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// ipRateLimiter manages per-IP rate limiters with automatic cleanup.
type ipRateLimiter struct {
	mu       sync.Mutex
	limiters map[string]*ipLimiter
	r        rate.Limit
	b        int
}

func newIPRateLimiter(r rate.Limit, b int) *ipRateLimiter {
	rl := &ipRateLimiter{
		limiters: make(map[string]*ipLimiter),
		r:        r,
		b:        b,
	}
	// 定时清理超过 10 分钟未活跃的 IP 记录，防止内存泄漏
	go rl.cleanup()
	return rl
}

func (rl *ipRateLimiter) getLimiter(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if v, ok := rl.limiters[ip]; ok {
		v.lastSeen = time.Now()
		return v.limiter
	}

	l := &ipLimiter{
		limiter:  rate.NewLimiter(rl.r, rl.b),
		lastSeen: time.Now(),
	}
	rl.limiters[ip] = l
	return l.limiter
}

func (rl *ipRateLimiter) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		rl.mu.Lock()
		for ip, v := range rl.limiters {
			if time.Since(v.lastSeen) > 10*time.Minute {
				delete(rl.limiters, ip)
			}
		}
		rl.mu.Unlock()
	}
}

// 全局速率限制器：所有 IP 默认 100 req/min（约 1.67 req/s，burst 20）
var globalLimiter = newIPRateLimiter(rate.Every(600*time.Millisecond), 20)

// 认证接口（发验证码、登录等）严格限速：10 req/min（约 1 req/6s，burst 5）
var authLimiter = newIPRateLimiter(rate.Every(6*time.Second), 5)

// getClientIP 获取真实客户端 IP，兼容 X-Forwarded-For 和 X-Real-IP
func getClientIP(c *gin.Context) string {
	if xff := c.GetHeader("X-Forwarded-For"); xff != "" {
		// X-Forwarded-For 可能包含多个 IP，取第一个（最原始的客户端 IP）
		for i := 0; i < len(xff); i++ {
			if xff[i] == ',' {
				return xff[:i]
			}
		}
		return xff
	}
	if xri := c.GetHeader("X-Real-IP"); xri != "" {
		return xri
	}
	return c.ClientIP()
}

// GlobalRateLimitMiddleware 全局 IP 限速中间件（TASK-11）
func GlobalRateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := getClientIP(c)
		limiter := globalLimiter.getLimiter(ip)
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "请求过于频繁，请稍后再试",
				"code":  429,
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// AuthRateLimitMiddleware 认证接口严格限速中间件（TASK-11）
func AuthRateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := getClientIP(c)
		limiter := authLimiter.getLimiter(ip)
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "操作过于频繁，请等待后再试",
				"code":  429,
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
