package admin

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware protects admin endpoints with a static token.
// Provide token via X-Admin-Token header (preferred) or Authorization: Bearer <token>.
func AuthMiddleware() gin.HandlerFunc {
	expected := strings.TrimSpace(os.Getenv("ADMIN_TOKEN"))
	if expected == "" {
		log.Fatal("ADMIN_TOKEN 环境变量未配置，无法启动管理接口")
	}
	// TASK-04: 强制最小长度，防止使用弱 Token
	if len(expected) < 32 {
		log.Fatal("ADMIN_TOKEN 至少需要 32 个字符，建议使用以下命令生成：\n  openssl rand -hex 32")
	}

	return func(c *gin.Context) {

		token := strings.TrimSpace(c.GetHeader("X-Admin-Token"))
		if token == "" {
			authHeader := strings.TrimSpace(c.GetHeader("Authorization"))
			if strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
				token = strings.TrimSpace(authHeader[7:])
			}
		}
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing admin token"})
			c.Abort()
			return
		}

		if subtle.ConstantTimeCompare([]byte(token), []byte(expected)) != 1 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid admin token"})
			c.Abort()
			return
		}

		sum := sha256.Sum256([]byte(token))
		c.Set("adminOperatorSource", "admin_console")
		c.Set("adminOperatorID", hex.EncodeToString(sum[:8]))
		c.Next()
	}
}
