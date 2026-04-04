package api

import (
	"strings"
	"testing"
)

// TestIsValidEmail 验证邮箱格式校验函数
func TestIsValidEmail(t *testing.T) {
	valid := []string{
		"user@example.com",
		"user.name+tag@domain.co",
		"user123@sub.domain.org",
		"test@163.com",
		"hello@qq.com",
	}
	invalid := []string{
		"",
		"notanemail",
		"missing@",
		"@nodomain",
		"no spaces@domain.com",
		"double@@domain.com",
	}

	for _, email := range valid {
		if !isValidEmail(email) {
			t.Errorf("isValidEmail(%q) 应返回 true", email)
		}
	}
	for _, email := range invalid {
		if isValidEmail(email) {
			t.Errorf("isValidEmail(%q) 应返回 false", email)
		}
	}
}

// TestEscapeLIKE 验证 LIKE 特殊字符转义
func TestEscapeLIKE(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello", "hello"},
		{"100%", `100\%`},
		{"user_name", `user\_name`},
		{`back\slash`, `back\\slash`},
		{"%_%", `\%\_\%`},
		{"", ""},
	}

	for _, tt := range tests {
		result := escapeLIKE(tt.input)
		if result != tt.expected {
			t.Errorf("escapeLIKE(%q) = %q, 期望 %q", tt.input, result, tt.expected)
		}
	}
}

// TestEscapeLIKE_NoSQLInjection 验证无法通过 LIKE 进行 DoS
func TestEscapeLIKE_NoSQLInjection(t *testing.T) {
	// 大量 % 字符应被转义，不能触发高开销扫描
	malicious := strings.Repeat("%", 100)
	result := escapeLIKE(malicious)
	expected := strings.Repeat(`\%`, 100)
	if result != expected {
		t.Errorf("escapeLIKE 对大量 %% 未正确转义: got len=%d, want len=%d", len(result), len(expected))
	}
}

// TestIsAllowedImageType 验证图片类型 Magic Bytes 检测
func TestIsAllowedImageType(t *testing.T) {
	validTypes := []struct {
		name  string
		bytes []byte
	}{
		{"JPEG", []byte{0xFF, 0xD8, 0xFF, 0xE0, 0x00, 0x10}},
		{"PNG", []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, 0x00}},
		{"GIF87a", []byte{0x47, 0x49, 0x46, 0x38, 0x37, 0x61}},
		{"GIF89a", []byte{0x47, 0x49, 0x46, 0x38, 0x39, 0x61}},
		// WebP: RIFF....WEBP
		{"WebP", []byte{0x52, 0x49, 0x46, 0x46, 0x00, 0x00, 0x00, 0x00, 0x57, 0x45, 0x42, 0x50}},
	}
	invalidTypes := []struct {
		name  string
		bytes []byte
	}{
		{"PDF", []byte{0x25, 0x50, 0x44, 0x46}},
		{"ZIP", []byte{0x50, 0x4B, 0x03, 0x04}},
		{"EXE", []byte{0x4D, 0x5A, 0x90, 0x00}},
		{"空字节", []byte{}},
		{"随机数据", []byte{0x00, 0x01, 0x02, 0x03}},
	}

	for _, tt := range validTypes {
		if !isAllowedImageType(tt.bytes) {
			t.Errorf("isAllowedImageType %s 应返回 true", tt.name)
		}
	}
	for _, tt := range invalidTypes {
		if isAllowedImageType(tt.bytes) {
			t.Errorf("isAllowedImageType %s 应返回 false", tt.name)
		}
	}
}

// TestAmountsEqual 验证金额比较（来自 payment_handlers.go 的 amountsEqual）
// 此测试确保"10" == "10.00"，防止金额校验绕过（TASK-03 相关）
func TestAmountsEqual(t *testing.T) {
	tests := []struct {
		a, b     string
		expected bool
	}{
		{"10", "10.00", true},
		{"129", "129.0", true},
		{"10.5", "10.50", true},
		{"100", "200", false},
		{"abc", "10", false},
		{"10", "abc", false},
		{"", "", true}, // 两个都无法解析，回退到字符串比较
	}
	for _, tt := range tests {
		result := amountsEqual(tt.a, tt.b)
		if result != tt.expected {
			t.Errorf("amountsEqual(%q, %q) = %v, 期望 %v", tt.a, tt.b, result, tt.expected)
		}
	}
}
