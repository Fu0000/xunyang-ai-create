package api

import (
	"encoding/base64"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"google-ai-proxy/internal/storage"
)

// allowedImageMagicBytes TASK-05: 允许的图片类型 Magic Bytes
// JPEG: FF D8 FF | PNG: 89 50 4E 47 | GIF: 47 49 46 38 | WebP: 52 49 46 46 xx xx xx xx 57 45 42 50
var allowedImageMagicBytes = []struct {
	offset int
	magic  []byte
}{
	{0, []byte{0xFF, 0xD8, 0xFF}},                               // JPEG
	{0, []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}}, // PNG
	{0, []byte{0x47, 0x49, 0x46, 0x38}},                         // GIF
	{8, []byte{0x57, 0x45, 0x42, 0x50}},                         // WebP (RIFF....WEBP)
}

// isAllowedImageType TASK-05: 检查字节数组是否为允许的图片类型
func isAllowedImageType(data []byte) bool {
	for _, m := range allowedImageMagicBytes {
		if len(data) >= m.offset+len(m.magic) {
			match := true
			for i, b := range m.magic {
				if data[m.offset+i] != b {
					match = false
					break
				}
			}
			if match {
				return true
			}
		}
	}
	return false
}

// validateBase64Image TASK-05: 对 Base64 图片做大小和文件类型校验
// 返回解码后的字节和错误（会直接写 HTTP 响应）
func validateBase64Image(c *gin.Context, b64 string) ([]byte, bool) {
	if b64 == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "图片不能为空"})
		return nil, false
	}
	// 剥除 data URL 前缀（如 "data:image/png;base64,"）
	if idx := strings.Index(b64, ","); idx != -1 {
		b64 = b64[idx+1:]
	}

	// Base64 长度上限：10MB 对应约 13.6M Base64 字符
	const maxBase64Len = 14 * 1024 * 1024
	if len(b64) > maxBase64Len {
		c.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "图片过大，最大支持 10MB"})
		return nil, false
	}

	data, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		// 尝试 URL-safe base64
		data, err = base64.RawStdEncoding.DecodeString(b64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "图片格式无效"})
			return nil, false
		}
	}

	// 大小上限：10MB
	const maxSize = 10 * 1024 * 1024
	if len(data) > maxSize {
		c.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "图片过大，最大支持 10MB"})
		return nil, false
	}

	// Magic Bytes 类型校验
	if !isAllowedImageType(data) {
		c.JSON(http.StatusUnsupportedMediaType, gin.H{"error": "不支持的图片格式，请上传 JPEG/PNG/WebP/GIF"})
		return nil, false
	}

	return data, true
}

// UploadImage 上传图片到 OSS，返回 URL（TASK-05：已加大小和 MIME 校验）
func UploadImage(c *gin.Context) {
	userID := c.GetUint64("userID")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "请先登录"})
		return
	}

	var req UploadImageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求格式无效"})
		return
	}

	// TASK-05: 校验图片大小和文件类型
	if _, ok := validateBase64Image(c, req.Image); !ok {
		return
	}

	userIDStr := strconv.FormatUint(userID, 10)
	url, err := storage.UploadBase64Image(req.Image, userIDStr, "useredit")
	if err != nil {
		log.Printf("上传图片失败 [用户:%d]: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "上传图片失败"})
		return
	}

	c.JSON(http.StatusOK, UploadImageResponse{URL: url})
}

// UploadVideo uploads a user provided video file to OSS and returns public URL.
func UploadVideo(c *gin.Context) {
	userID := c.GetUint64("userID")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "please login first"})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing file"})
		return
	}
	if file.Size <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "empty file"})
		return
	}
	if file.Size > 100*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "video file too large"})
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	switch ext {
	case ".mp4", ".mov", ".webm", ".m4v":
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported video format"})
		return
	}

	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read file"})
		return
	}
	defer src.Close()

	videoData, err := io.ReadAll(src)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read video data"})
		return
	}

	userIDStr := strconv.FormatUint(userID, 10)
	url, err := storage.UploadVideoData(videoData, userIDStr, ext)
	if err != nil {
		log.Printf("upload video failed [user:%d]: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upload video"})
		return
	}

	c.JSON(http.StatusOK, UploadImageResponse{URL: url})
}
