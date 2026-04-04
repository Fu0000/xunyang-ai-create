package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorResponse TASK-13: 统一错误响应格式
// 所有 API 错误均通过此结构体返回，确保前端可以一致地处理错误
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// ErrBadRequest 400 - 请求参数无效
func ErrBadRequest(c *gin.Context, msg string) {
	c.JSON(http.StatusBadRequest, ErrorResponse{
		Code:    http.StatusBadRequest,
		Message: msg,
	})
}

// ErrUnauthorized 401 - 未登录或 Token 无效
func ErrUnauthorized(c *gin.Context, msg string) {
	c.JSON(http.StatusUnauthorized, ErrorResponse{
		Code:    http.StatusUnauthorized,
		Message: msg,
	})
}

// ErrForbidden 403 - 无权限
func ErrForbidden(c *gin.Context, msg string) {
	c.JSON(http.StatusForbidden, ErrorResponse{
		Code:    http.StatusForbidden,
		Message: msg,
	})
}

// ErrNotFound 404 - 资源不存在
func ErrNotFound(c *gin.Context, msg string) {
	c.JSON(http.StatusNotFound, ErrorResponse{
		Code:    http.StatusNotFound,
		Message: msg,
	})
}

// ErrConflict 409 - 资源冲突
func ErrConflict(c *gin.Context, msg string) {
	c.JSON(http.StatusConflict, ErrorResponse{
		Code:    http.StatusConflict,
		Message: msg,
	})
}

// ErrTooManyRequests 429 - 请求过于频繁
func ErrTooManyRequests(c *gin.Context, msg string) {
	c.JSON(http.StatusTooManyRequests, ErrorResponse{
		Code:    http.StatusTooManyRequests,
		Message: msg,
	})
}

// ErrInternal 500 - 服务器内部错误
func ErrInternal(c *gin.Context, msg string) {
	c.JSON(http.StatusInternalServerError, ErrorResponse{
		Code:    http.StatusInternalServerError,
		Message: msg,
	})
}

// ErrPaymentRequired 402 - 钻石不足，附带额外信息
func ErrPaymentRequired(c *gin.Context, msg string, extra gin.H) {
	resp := gin.H{
		"code":    http.StatusPaymentRequired,
		"message": msg,
	}
	for k, v := range extra {
		resp[k] = v
	}
	c.JSON(http.StatusPaymentRequired, resp)
}

// ErrUnsupportedMedia 415 - 不支持的文件类型
func ErrUnsupportedMedia(c *gin.Context, msg string) {
	c.JSON(http.StatusUnsupportedMediaType, ErrorResponse{
		Code:    http.StatusUnsupportedMediaType,
		Message: msg,
	})
}

// ErrServiceUnavailable 503 - 服务暂时不可用
func ErrServiceUnavailable(c *gin.Context, msg string) {
	c.JSON(http.StatusServiceUnavailable, ErrorResponse{
		Code:    http.StatusServiceUnavailable,
		Message: msg,
	})
}
