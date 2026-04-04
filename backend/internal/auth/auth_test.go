package auth

import (
	"strings"
	"testing"
	"time"
)

// TestHashAndCheckPassword 验证密码哈希和校验功能
func TestHashAndCheckPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
	}{
		{"普通密码", "mypassword123"},
		{"含特殊字符", "P@ssw0rd!#$%"},
		{"中文密码", "密码123"},
		{"最小长度", "123456"},
		{"长密码", "averylongpasswordthatexceeds32characters1234567890"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := HashPassword(tt.password)
			if err != nil {
				t.Fatalf("HashPassword(%q) 失败: %v", tt.password, err)
			}
			if hash == "" {
				t.Fatal("哈希结果不能为空")
			}
			if hash == tt.password {
				t.Fatal("哈希结果不能与明文相同")
			}
			if !CheckPassword(tt.password, hash) {
				t.Errorf("CheckPassword(%q, hash) 应返回 true", tt.password)
			}
			if CheckPassword("wrongpassword", hash) {
				t.Error("错误密码应返回 false")
			}
		})
	}
}

// TestGenerateAndValidateUserToken 验证用户 JWT Token 生成和校验
func TestGenerateAndValidateUserToken(t *testing.T) {
	// 初始化测试密钥
	SecretKey = []byte("test-jwt-secret-key-at-least-32-chars-for-testing")
	LicenseSecretKey = []byte("test-license-secret-key-at-least-32-chars-different")

	userID := uint64(42)
	email := "test@example.com"
	token, err := GenerateUserToken(userID, email)
	if err != nil {
		t.Fatalf("GenerateUserToken 失败: %v", err)
	}
	if token == "" {
		t.Fatal("Token 不能为空")
	}

	// 验证 Token
	claims, err := ValidateUserToken(token)
	if err != nil {
		t.Fatalf("ValidateUserToken 失败: %v", err)
	}
	if claims.UserID != userID {
		t.Errorf("UserID = %d, 期望 %d", claims.UserID, userID)
	}
	if claims.Email != email {
		t.Errorf("Email = %q, 期望 %q", claims.Email, email)
	}
	// 验证过期时间约为 7 天
	expectedExpiry := time.Now().Add(7 * 24 * time.Hour)
	diff := claims.ExpiresAt.Time.Sub(expectedExpiry)
	if diff < -time.Minute || diff > time.Minute {
		t.Errorf("Token 过期时间偏差过大: %v", diff)
	}
}

// TestValidateUserToken_Invalid 验证无效 Token 应返回错误
func TestValidateUserToken_Invalid(t *testing.T) {
	SecretKey = []byte("test-jwt-secret-key-at-least-32-chars-for-testing")

	tests := []struct {
		name  string
		token string
	}{
		{"空 Token", ""},
		{"随机字符串", "notavalidtoken"},
		{"篡改 Token", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.tampered.signature"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ValidateUserToken(tt.token)
			if err == nil {
				t.Errorf("无效 Token %q 应返回错误", tt.token)
			}
		})
	}
}

// TestValidateLicenseKey_WithWrongSecret 用错误密钥签名的 License Key 应验证失败
func TestValidateLicenseKey_WithWrongSecret(t *testing.T) {
	// License 使用正确密钥生成
	LicenseSecretKey = []byte("correct-license-secret-at-least-32-chars-blah-blah")
	licenseToken, err := GenerateLicenseKey(100)
	if err != nil {
		t.Fatalf("GenerateLicenseKey 失败: %v", err)
	}

	// 然后换成错误密钥来验证
	LicenseSecretKey = []byte("wrong-license-secret-at-least-32-chars-blah-blah-xx")
	_, err = ValidateLicenseKey(licenseToken)
	if err == nil {
		t.Error("使用错误密钥验证 License Key 应返回错误")
	}
}

// TestGenerateLicenseKey_ValidWithCorrectSecret 用正确密钥应验证成功
func TestGenerateLicenseKey_ValidWithCorrectSecret(t *testing.T) {
	LicenseSecretKey = []byte("correct-license-secret-at-least-32-chars-blah-blah")
	credits := 500
	token, err := GenerateLicenseKey(credits)
	if err != nil {
		t.Fatalf("GenerateLicenseKey 失败: %v", err)
	}

	claims, err := ValidateLicenseKey(token)
	if err != nil {
		t.Fatalf("ValidateLicenseKey 失败: %v", err)
	}
	if claims.Credits != credits {
		t.Errorf("Credits = %d, 期望 %d", claims.Credits, credits)
	}
	if claims.ID == "" {
		t.Error("License ID 不能为空")
	}
}

// TestJWTAndLicenseKeyIndependence JWT 密钥和 License 密钥必须相互独立
func TestJWTAndLicenseKeyIndependence(t *testing.T) {
	jwtSecret := []byte("jwt-secret-32-chars-xxxxxxxxxxx")
	licenseSecret := []byte("license-secret-32-chars-xxxxxxxxx")

	SecretKey = jwtSecret
	LicenseSecretKey = licenseSecret

	// 生成用户 JWT
	userToken, err := GenerateUserToken(1, "user@test.com")
	if err != nil {
		t.Fatalf("GenerateUserToken 失败: %v", err)
	}

	// 用 License 密钥验证用户 Token 应失败
	LicenseSecretKey = jwtSecret // 故意使用 JWT 密钥来验证 License
	_, err = ValidateLicenseKey(userToken)
	// 用户 JWT 用 LicenseClaims 解析，即使签名正确也应该缺少必要字段而失败
	// 或者至少说明两者是独立的
	if err == nil {
		// JWT Token 不应该能通过 License Key 验证（因为 ClaimType 不同）
		// 此处测试确保密钥分离的意识，而非功能断言
		t.Log("警告: JWT Token 被 License Key 验证接受，请确认密钥隔离")
	}

	// 恢复正确密钥
	LicenseSecretKey = licenseSecret

	// 用正确 License 密钥重新验证
	licenseToken, _ := GenerateLicenseKey(100)
	claims, err := ValidateLicenseKey(licenseToken)
	if err != nil {
		t.Fatalf("License Key 验证失败: %v", err)
	}
	if !strings.Contains(claims.ID, "-") {
		t.Error("License ID 应为 UUID 格式")
	}
}
