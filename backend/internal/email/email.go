package email

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	smtpHost     string
	smtpPort     string
	smtpUser     string
	smtpPassword string
	fromEmail    string
	fromName     string

	cfWorkerBase    string
	cfAdminPassword string
	cfDefaultDomain string
	cfJwtSecret     string
)

// InitEmail 初始化邮件配置
func InitEmail() {
	smtpHost = os.Getenv("SMTP_HOST")
	if smtpHost == "" {
		smtpHost = "smtp.qq.com"
	}
	smtpPort = os.Getenv("SMTP_PORT")
	if smtpPort == "" {
		smtpPort = "587"
	}
	smtpUser = os.Getenv("SMTP_USER")
	smtpPassword = os.Getenv("SMTP_PASSWORD")
	fromEmail = os.Getenv("SMTP_FROM_EMAIL")
	if fromEmail == "" {
		fromEmail = smtpUser
	}
	fromName = os.Getenv("SMTP_FROM_NAME")
	if fromName == "" {
		fromName = "寻氧AI"
	}

	cfWorkerBase = os.Getenv("CF_WORKER_BASE")
	cfAdminPassword = os.Getenv("CF_ADMIN_PASSWORD")
	cfDefaultDomain = os.Getenv("CF_DEFAULT_DOMAIN")
	cfJwtSecret = os.Getenv("CF_JWT_SECRET")

	if cfWorkerBase != "" {
		log.Printf("邮件服务已配置: 采用 CF Worker 代理模式 (%s)", cfWorkerBase)
	} else if smtpUser != "" && smtpPassword != "" {
		log.Printf("邮件服务已配置: %s:%s", smtpHost, smtpPort)
	} else {
		log.Printf("警告: 邮件服务未配置 (CF_WORKER / SMTP 未设置)")
	}
}

// IsConfigured 检查邮件服务是否配置
func IsConfigured() bool {
	return cfWorkerBase != "" || (smtpUser != "" && smtpPassword != "")
}

// SendVerificationCode 发送验证码邮件
func SendVerificationCode(to, code, purpose string) error {
	if !IsConfigured() {
		log.Printf("邮件服务未配置，验证码 [%s] 发送到 [%s] (类型: %s)", code, to, purpose)
		return nil // 在开发环境下不报错
	}

	var subject, body string

	switch purpose {
	case "register":
		subject = "寻氧AI - 注册验证码"
		body = fmt.Sprintf(`
			<html>
			<body style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto;">
				<div style="background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%); padding: 30px; text-align: center;">
					<h1 style="color: white; margin: 0;">🎨 寻氧AI</h1>
				</div>
				<div style="padding: 30px; background: #f9f9f9;">
					<h2 style="color: #333;">欢迎注册 寻氧AI</h2>
					<p style="color: #666;">您的验证码是：</p>
					<div style="background: white; padding: 20px; text-align: center; border-radius: 8px; margin: 20px 0;">
						<span style="font-size: 32px; font-weight: bold; color: #667eea; letter-spacing: 8px;">%s</span>
					</div>
					<p style="color: #999; font-size: 14px;">验证码有效期为10分钟，请尽快使用。</p>
					<p style="color: #999; font-size: 14px;">如果这不是您的操作，请忽略此邮件。</p>
				</div>
				<div style="background: #333; padding: 20px; text-align: center;">
					<p style="color: #999; margin: 0; font-size: 12px;">© %d 寻氧AI. All rights reserved.</p>
				</div>
			</body>
			</html>
		`, code, time.Now().Year())
	case "login":
		subject = "寻氧AI - 登录验证码"
		body = fmt.Sprintf(`
			<html>
			<body style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto;">
				<div style="background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%); padding: 30px; text-align: center;">
					<h1 style="color: white; margin: 0;">🎨 寻氧AI</h1>
				</div>
				<div style="padding: 30px; background: #f9f9f9;">
					<h2 style="color: #333;">登录验证</h2>
					<p style="color: #666;">您的登录验证码是：</p>
					<div style="background: white; padding: 20px; text-align: center; border-radius: 8px; margin: 20px 0;">
						<span style="font-size: 32px; font-weight: bold; color: #667eea; letter-spacing: 8px;">%s</span>
					</div>
					<p style="color: #999; font-size: 14px;">验证码有效期为10分钟，请尽快使用。</p>
					<p style="color: #999; font-size: 14px;">如果这不是您的操作，请忽略此邮件。</p>
				</div>
				<div style="background: #333; padding: 20px; text-align: center;">
					<p style="color: #999; margin: 0; font-size: 12px;">© %d 寻氧AI. All rights reserved.</p>
				</div>
			</body>
			</html>
		`, code, time.Now().Year())
	case "reset":
		subject = "寻氧AI - 重置密码验证码"
		body = fmt.Sprintf(`
			<html>
			<body style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto;">
				<div style="background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%); padding: 30px; text-align: center;">
					<h1 style="color: white; margin: 0;">🎨 寻氧AI</h1>
				</div>
				<div style="padding: 30px; background: #f9f9f9;">
					<h2 style="color: #333;">重置密码</h2>
					<p style="color: #666;">您正在重置密码，验证码是：</p>
					<div style="background: white; padding: 20px; text-align: center; border-radius: 8px; margin: 20px 0;">
						<span style="font-size: 32px; font-weight: bold; color: #667eea; letter-spacing: 8px;">%s</span>
					</div>
					<p style="color: #999; font-size: 14px;">验证码有效期为10分钟，请尽快使用。</p>
					<p style="color: #999; font-size: 14px;">如果这不是您的操作，请立即检查账号安全。</p>
				</div>
				<div style="background: #333; padding: 20px; text-align: center;">
					<p style="color: #999; margin: 0; font-size: 12px;">© %d 寻氧AI. All rights reserved.</p>
				</div>
			</body>
			</html>
		`, code, time.Now().Year())
	default:
		subject = "寻氧AI - 验证码"
		body = fmt.Sprintf(`您的验证码是: %s，有效期10分钟。`, code)
	}

	return sendHTML(to, subject, body)
}

// sendHTML 发送HTML格式邮件
func sendHTML(to, subject, body string) error {
	if cfWorkerBase != "" {
		return sendHTMLViaCFWorker(to, subject, body)
	}

	auth := smtp.PlainAuth("", smtpUser, smtpPassword, smtpHost)

	headers := make(map[string]string)
	headers["From"] = fmt.Sprintf("%s <%s>", fromName, fromEmail)
	headers["To"] = to
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=UTF-8"

	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	addr := fmt.Sprintf("%s:%s", smtpHost, smtpPort)
	
	var err error
	if smtpPort == "465" {
		// 针对 465 隐式 TLS 的处理方式，原生 smtp.SendMail 仅支持 587 STARTTLS
		tlsconfig := &tls.Config{
			InsecureSkipVerify: true,
			ServerName:         smtpHost,
		}
		conn, errConn := tls.Dial("tcp", addr, tlsconfig)
		if errConn != nil {
			err = errConn
		} else {
			defer conn.Close()
			client, errClient := smtp.NewClient(conn, smtpHost)
			if errClient != nil {
				err = errClient
			} else {
				defer client.Close()
				if errAuth := client.Auth(auth); errAuth != nil {
					err = errAuth
				} else if errMail := client.Mail(fromEmail); errMail != nil {
					err = errMail
				} else if errRcpt := client.Rcpt(to); errRcpt != nil {
					err = errRcpt
				} else {
					w, errData := client.Data()
					if errData != nil {
						err = errData
					} else {
						_, errWrite := w.Write([]byte(message))
						errClose := w.Close()
						client.Quit()
						if errWrite != nil {
							err = errWrite
						} else {
							err = errClose
						}
					}
				}
			}
		}
	} else {
		// 标准端口或支持 STARTTLS 的 587 端口走原生库
		err = smtp.SendMail(addr, auth, fromEmail, []string{to}, []byte(message))
	}

	if err != nil {
		// QQ邮箱的SMTP服务器有时会在邮件发送成功后返回不完整的响应
		// 导致 "short response" 错误，但邮件实际上已经发送成功
		// 参考: https://github.com/golang/go/issues/24845
		errStr := err.Error()
		if strings.HasPrefix(errStr, "short response") {
			log.Printf("邮件已发送成功 (忽略 short response 错误): %s <- %s", to, subject)
			return nil
		}
		log.Printf("发送邮件失败 [%s]: %v", to, err)
		return fmt.Errorf("发送邮件失败: %w", err)
	}

	log.Printf("邮件发送成功: %s <- %s", to, subject)
	return nil
}

// generateCFMailToken 模拟 CF Worker 的 Jwt().sign({ address }) 逻辑
func generateCFMailToken(address, secret string) string {
	claims := jwt.MapClaims{
		"address": address,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(secret))
	return tokenString
}

// sendHTMLViaCFWorker 通过 Cloudflare Worker 大部分标准的 API 格式发送
func sendHTMLViaCFWorker(to, subject, body string) error {
	endpoint := strings.TrimRight(cfWorkerBase, "/") + "/api/send_mail"
	
	senderEmail := fromEmail
	// 对于 CF Worker，必须使用在其域名下经过校验的发件地址
	if !strings.HasSuffix(senderEmail, "@"+cfDefaultDomain) {
		senderEmail = "no-reply@" + cfDefaultDomain
	}

	// 适配最常见开源 CF 邮件代理 (如 dreamhunter2333/cloudflare_temp_email) 的精确 JSON 格式和认证
	payload := map[string]interface{}{
		"to_mail":   to,
		"to_name":   to,
		"from_name": fromName,
		// 发信地址会在 worker 内部结合 address 参数或者管理员覆写被处理
		"subject":   subject,
		"content":   body,
		"is_html":   true,
	}
	
	jsonData, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("构建 CF worker 请求失败: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	
	if cfJwtSecret != "" {
		tokenString := generateCFMailToken(senderEmail, cfJwtSecret)
		req.Header.Set("Authorization", "Bearer "+tokenString)
	} else {
		// 回退到 admin header 如果没有提供 jwt secret
		req.Header.Set("x-admin-password", cfAdminPassword)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("CF Worker 邮件网络请求失败 [%s]: %v", to, err)
		return fmt.Errorf("网络请求失败: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		log.Printf("CF Worker 返回错误 [%d]: %s", resp.StatusCode, string(respBody))
		return fmt.Errorf("CF Worker 发送失败 (状态码 %d)", resp.StatusCode)
	}

	log.Printf("CF Worker 邮件发送成功: %s <- %s", to, subject)
	return nil
}
