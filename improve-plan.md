# 寻氧AI 改进计划 (Improvement Plan)

> 创建时间：2026-04-04  
> 维护规则：每完成一个任务后，在对应条目下填写「完成日志」，并执行本地提交和 README 更新  
> 参考分析：`project-plan.md`

---

## 工作流规范（必读）

每个任务完成后，必须按以下步骤执行：

```
1. 代码修改完成
2. 本地测试（见各任务的「验证步骤」）
3. git add <相关文件>
4. git commit -m "<type>(<scope>): <简短描述>"
5. 更新 README.md（如涉及配置/接口变更）
6. 在本文档对应任务下填写「✅ 完成日志」
```

### Commit Message 规范（Conventional Commits）

```
feat(backend): 添加速率限制中间件
fix(auth): 修复 JWT 与 License Key 共用密钥问题
perf(db): 配置数据库连接池参数
refactor(api): 提取 service 层减小 handler 文件体积
test(auth): 添加密码哈希和 JWT 验证单元测试
docs(readme): 更新部署和配置文档
chore(ci): 添加 GitHub Actions 工作流
```

**类型说明：**
| Type | 场景 |
|------|------|
| `feat` | 新功能 |
| `fix` | Bug 修复 |
| `perf` | 性能优化 |
| `refactor` | 代码重构（不影响功能） |
| `test` | 测试相关 |
| `docs` | 文档更新 |
| `chore` | 构建/工具链/CI |
| `security` | 安全修复 |

---

## 第一阶段：安全加固（P0 - 最高优先级）

> 目标：消除已知安全漏洞，防止恶意攻击造成的业务损失

---

### TASK-01｜JWT 与 License Key 密钥分离

**优先级**: 🔴 P0  
**分类**: 安全  
**文件**: `backend/internal/auth/auth.go`, `backend/internal/config/config.go`, `backend/.env.example`  
**状态**: `[x]` 已完成

#### 完成日志

```
完成时间: 2026-04-05
修改文件: backend/internal/auth/auth.go, backend/internal/config/config.go, backend/.env.example
Commit Hash: e1d11a1/badb7dc
备注: LICENSE_SECRET 独立密钥，启动时必强校验，缺失则 Fatal 退出
```


#### 问题描述

当前 `LicenseSecretKey` 和 `SecretKey`（用户 JWT）共享同一个 `JWT_SECRET` 环境变量。任何知道 `JWT_SECRET` 的人都可以伪造合法的 License Key，从而无限兑换积分。

```go
// auth.go:27 - 当前错误实现
LicenseSecretKey = []byte(secret)  // 与 JWT SecretKey 完全相同！
```

#### 修改方案

1. 在 `config/config.go` 中添加 `GetLicenseSecret()` 函数
2. 在 `auth/auth.go` 的 `InitSecretKey()` 中读取独立的 `LICENSE_SECRET` 环境变量
3. 更新 `backend/.env.example` 添加 `LICENSE_SECRET` 配置项
4. 如 `LICENSE_SECRET` 未设置，用 `log.Fatal` 阻止启动（而非 fallback 到 JWT_SECRET）

#### 验证步骤

```bash
# 1. 重新编译后端
cd backend && go build ./...

# 2. 不设置 LICENSE_SECRET 启动 → 应报错退出
LICENSE_SECRET= go run main.go
# 预期输出：fatal: LICENSE_SECRET 环境变量未设置

# 3. 设置不同密钥后，用旧 JWT_SECRET 签名的 License Key 应无法兑换
# 手动测试 /api/user/redeem 接口
```

#### README 更新内容

在「必填配置项」表格中添加 `LICENSE_SECRET` 行，说明其用途和生成方式。

#### 完成日志

```
完成时间: 
修改文件: 
Commit Hash: 
备注: 
```

---

### TASK-02｜验证码防暴力破解

**优先级**: 🔴 P0  
**分类**: 安全  
**文件**: `backend/internal/db/db.go`（添加 verify_attempts 字段）, `backend/internal/api/auth_handlers.go`  
**状态**: `[x]` 已完成

#### 完成日志

```
完成时间: 2026-04-05
修改文件: backend/migrations/027_add_verification_attempts.sql, backend/internal/api/utils.go
Commit Hash: e1d11a1/badb7dc
备注: verifyCode() 统一封装，超过 5 次错误后验证码自动失效
```

#### 问题描述

验证码为 6 位纯数字，共 100 万种组合。当前只有"1 分钟内不能重发"的限速，但对验证码猜测没有限制，攻击者可持续轮询 `/auth/login` 或 `/auth/register` 枚举正确验证码。

#### 修改方案

1. 在 `email_verifications` 表中添加 `attempts INT DEFAULT 0` 字段（新建 migration 文件）
2. 每次验证失败时 `attempts + 1`
3. `attempts >= 5` 时将该验证码记录标记为已使用（失效），返回 `429` 并提示"验证码已失效，请重新获取"
4. 新增 migration：`backend/migrations/027_add_verification_attempts.sql`

#### 验证步骤

```bash
# 1. 连续发送 5 次错误验证码到同一邮箱账户
# 请求 POST /api/auth/register with wrong code × 5
# 第 5 次后，正确验证码也应返回 400「验证码已失效」

# 2. 重新发送验证码后，新验证码应可正常使用

# 3. 运行迁移测试
cd backend && go run main.go  # 观察 migration 027 是否正常应用
```

#### README 更新内容

无需更新（内部安全机制，不影响用户配置）。

#### 完成日志

```
完成时间: 
修改文件: 
Commit Hash: 
备注: 
```

---

### TASK-03｜支付金额二次校验修复

**优先级**: 🔴 P0  
**分类**: 安全  
**文件**: `backend/internal/api/payment_handlers.go`  
**状态**: `[x]` 已完成

#### 完成日志

```
完成时间: 2026-04-05
修改文件: backend/internal/api/payment_handlers.go
Commit Hash: e1d11a1/badb7dc
备注: fulfillPaymentOrder 加入 expectedAmount 参数，主动查询路径也传入金额进行二次校验
```

#### 问题描述

在 `GetPaymentStatus` 主动查询支付状态的路径中，调用 `fulfillPaymentOrder` 时没有传入并校验实际支付金额，存在绕过金额校验的潜在风险。

```go
// 当前代码 - 未校验金额
fulfillPaymentOrder(order.OrderNo, result.TradeNo, "")
```

#### 修改方案

1. 修改 `fulfillPaymentOrder` 函数签名，增加 `expectedAmount string` 参数
2. 函数内部在标记订单前校验 `order.Amount == expectedAmount`
3. 更新所有调用方（`LinuxDoCreditNotify` 和 `GetPaymentStatus`）传入正确的期望金额

#### 验证步骤

```bash
cd backend && go build ./...  # 编译通过

# 手动测试主动查询路径：
# 1. 创建订单
# 2. 模拟支付成功但金额不匹配的回调
# 3. 调用 GET /api/user/payment/status/:orderNo
# 4. 验证订单状态未被错误标记为 paid
```

#### README 更新内容

无需更新。

#### 完成日志

```
完成时间: 
修改文件: 
Commit Hash: 
备注: 
```

---

### TASK-04｜管理员 Token 最小长度校验

**优先级**: 🔴 P0  
**分类**: 安全  
**文件**: `backend/internal/api/admin/middleware.go`  
**状态**: `[x]` 已完成

#### 完成日志

```
完成时间: 2026-04-05
修改文件: backend/internal/api/admin/middleware.go
Commit Hash: e1d11a1/badb7dc
备注: 添加 len(expected) >= 32 校验，不足则 Fatal 并提示生成命令
```

#### 问题描述

`ADMIN_TOKEN` 可以设置为任意短字符串（如 `"admin"`, `"123"`），弱 token 极易被猜测。

#### 修改方案

1. 在 `AuthMiddleware` 初始化时校验 `ADMIN_TOKEN` 长度 ≥ 32 字符
2. 若长度不足，通过 `log.Fatal()` 阻止服务启动，并打印生成建议：
   ```
   ADMIN_TOKEN 至少需要 32 个字符，建议使用以下命令生成：
   openssl rand -hex 32
   ```

#### 验证步骤

```bash
# 设置短 Token 启动服务
ADMIN_TOKEN=short go run main.go
# 预期：启动失败并打印提示

# 设置 32+ 字符 Token 启动
ADMIN_TOKEN=$(openssl rand -hex 32) go run main.go
# 预期：正常启动
```

#### README 更新内容

在「可选配置」的 `ADMIN_TOKEN` 说明中，补充"至少 32 字符"的要求和生成命令。

#### 完成日志

```
完成时间: 
修改文件: 
Commit Hash: 
备注: 
```

---

### TASK-05｜上传接口内容校验加固

**优先级**: 🔴 P0  
**分类**: 安全  
**文件**: `backend/internal/api/resource_handlers.go`, `backend/internal/storage/`  
**状态**: `[x]` 已完成

#### 完成日志

```
完成时间: 2026-04-05
修改文件: backend/internal/api/resource_handlers.go
Commit Hash: e1d11a1/badb7dc
备注: 添加 Magic Bytes 文件类型检测（JPEG/PNG/WebP/GIF）和 10MB 大小上限，不符则 413/415
```

#### 问题描述

图片上传接口仅接收 Base64 字符串，未校验：
- Base64 解码后的文件大小上限
- 文件类型（MIME Type / Magic Bytes）合法性

#### 修改方案

1. Base64 解码后检查文件大小 ≤ 10MB（图片）/ 500MB（视频）
2. 读取前 16 字节检查 Magic Bytes，只允许 JPEG、PNG、WebP、GIF 格式图片
3. 返回明确的错误信息（`413 Request Entity Too Large` / `415 Unsupported Media Type`）

#### 验证步骤

```bash
# 上传超大 Base64 图片 → 应返回 413
# 上传 PDF 文件的 Base64 → 应返回 415

go build ./...  # 编译验证
```

#### README 更新内容

在「API 说明」或部署文档中说明上传限制。

#### 完成日志

```
完成时间: 
修改文件: 
Commit Hash: 
备注: 
```

---

## 第二阶段：稳定性提升（P1 - 高优先级）

> 目标：提升系统在高并发和边界情况下的健壮性，避免运行时崩溃

---

### TASK-06｜数据库连接池配置

**优先级**: 🔴 P1  
**分类**: 性能  
**文件**: `backend/internal/db/db.go`  
**状态**: `[ ]` 未开始

#### 问题描述

GORM 默认不限制连接数，高并发下会创建大量连接直到数据库拒绝，导致服务雪崩。

#### 修改方案

在 `InitDB()` 中 `gorm.Open` 成功后，添加连接池配置：

```go
sqlDB, err := DB.DB()
if err != nil {
    log.Fatalf("获取 sql.DB 失败: %v", err)
}
sqlDB.SetMaxOpenConns(50)           // 最大打开连接数
sqlDB.SetMaxIdleConns(10)           // 最大空闲连接数
sqlDB.SetConnMaxLifetime(time.Hour) // 连接最长存活时间
sqlDB.SetConnMaxIdleTime(30 * time.Minute) // 空闲连接最长时间
```

同时支持通过环境变量自定义（`DB_MAX_OPEN_CONNS`, `DB_MAX_IDLE_CONNS`）。

#### 验证步骤

```bash
cd backend && go build ./...  # 编译通过
go run main.go  # 启动观察日志中是否有连接池相关信息

# 可选：并发压测工具 (wrk / ab) 验证连接数稳定
```

#### README 更新内容

在「可选配置」表格中添加 `DB_MAX_OPEN_CONNS` 和 `DB_MAX_IDLE_CONNS` 说明。

#### 完成日志

```
完成时间: 
修改文件: 
Commit Hash: 
备注: 
```

---

### TASK-07｜事务提交错误检查修复

**优先级**: 🔴 P1  
**分类**: 稳定性  
**文件**: `backend/internal/api/auth_handlers.go`  
**状态**: `[ ]` 未开始

#### 问题描述

用户注册流程中的事务提交未检查错误：

```go
// auth_handlers.go:299 - 当前代码
tx.Commit()  // 返回值被完全忽略！
```

若 `Commit()` 失败（如网络抖动、死锁），用户已注册但积分可能未到账，且前端收到"注册成功"响应，造成数据不一致。

#### 修改方案

1. 检查注册流程中所有 `tx.Commit()` 的返回值
2. 提交失败时返回 `500` 错误，并记录详细日志
3. 全局搜索项目中所有 `tx.Commit()` 调用，逐一修复

```bash
# 查找所有未检查的 tx.Commit()
grep -n "tx.Commit()" backend/internal/api/*.go
```

#### 验证步骤

```bash
cd backend && go build ./...  # 编译通过

# Code Review：搜索确认没有遗漏
grep -rn "\.Commit()" backend/internal/ | grep -v "\.Commit()\.Error"
# 预期：无输出（所有 Commit 都检查了错误）
```

#### README 更新内容

无需更新。

#### 完成日志

```
完成时间: 
修改文件: 
Commit Hash: 
备注: 
```

---

### TASK-08｜数据库迁移失败时阻止启动

**优先级**: 🔴 P1  
**分类**: 稳定性  
**文件**: `backend/internal/db/db.go`  
**状态**: `[ ]` 未开始

#### 问题描述

当前迁移脚本执行失败时，仅打印警告日志，服务继续启动。数据库结构不完整的情况下运行可能导致运行时 panic。

```go
// db.go:388 - 当前代码
log.Printf("migrations: %s had errors, not marking...", name)
// 继续运行，没有任何阻断！
```

#### 修改方案

1. 将可接受的"幂等性错误"（如"列已存在"、"表已存在"）与真正的失败区分
2. 对非幂等错误，调用 `log.Fatalf()` 阻止服务启动
3. 添加辅助函数 `isMigrationIdempotentError(err error) bool` 识别可忽略错误

#### 验证步骤

```bash
# 1. 手动写一个错误的 SQL 迁移文件
# 2. 启动服务 → 应 Fatal 退出，并打印迁移失败的文件名和错误信息
# 3. 修复 SQL → 应正常启动

go build ./...
```

#### README 更新内容

在「数据库」章节说明迁移失败行为，提示运维人员如何排查。

#### 完成日志

```
完成时间: 
修改文件: 
Commit Hash: 
备注: 
```

---

### TASK-09｜添加健康检查端点

**优先级**: 🟡 P1  
**分类**: DevOps  
**文件**: `backend/main.go`  
**状态**: `[ ]` 未开始

#### 问题描述

缺乏 `/health` 端点，负载均衡器（Nginx upstream check）、K8s 存活探针、监控系统无法判断服务状态。

#### 修改方案

在路由注册中添加：

```go
// 健康检查 - 不需要认证
r.GET("/health", func(c *gin.Context) {
    // 检查数据库连接
    sqlDB, err := db.DB.DB()
    dbStatus := "ok"
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
        "version":   "1.0.0",  // 可从构建变量注入
    })
})
```

#### 验证步骤

```bash
# 启动服务后测试
curl http://localhost:8092/health
# 预期返回：{"status":"ok","db":"ok","timestamp":...}

# 断开 DB 后测试
# 预期返回：{"status":"degraded","db":"error",...}
```

#### README 更新内容

在「部署」章节的 Nginx 配置示例中，添加健康检查的 upstream 配置示例。

#### 完成日志

```
完成时间: 
修改文件: 
Commit Hash: 
备注: 
```

---

### TASK-10｜邮箱正则预编译优化

**优先级**: 🟡 P1  
**分类**: 性能  
**文件**: `backend/internal/api/utils.go`  
**状态**: `[ ]` 未开始

#### 问题描述

```go
// 当前代码 - 每次调用都重新编译正则表达式
func isValidEmail(email string) bool {
    pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
    matched, _ := regexp.MatchString(pattern, email)
    return matched
}
```

每次邮箱验证都调用 `regexp.MatchString`，内部会重新编译正则，在高频调用场景下浪费 CPU。

#### 修改方案

```go
// 包级别预编译
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func isValidEmail(email string) bool {
    return emailRegex.MatchString(email)
}
```

#### 验证步骤

```bash
cd backend && go build ./...
go test ./internal/api/... -run TestEmailValidation -v
# 若无测试则手动验证 API 注册接口仍正常工作
```

#### README 更新内容

无需更新。

#### 完成日志

```
完成时间: 
修改文件: 
Commit Hash: 
备注: 
```

---

### TASK-11｜API 速率限制中间件

**优先级**: 🔴 P1  
**分类**: 安全/性能  
**文件**: `backend/internal/api/middleware.go` （新增 rate_limit.go）, `backend/main.go`, `backend/go.mod`  
**状态**: `[ ]` 未开始

#### 问题描述

整个 API 层没有速率限制，恶意用户可以：
- 对 `/api/auth/send-code` 发起大量请求（即使有 DB 层的1分钟限制，也无法防御分布式攻击）
- 对 `/api/generate` 进行积分刷取尝试

#### 修改方案

使用 `golang.org/x/time/rate`（无需额外依赖，标准库扩展）实现基于 IP 的速率限制：

1. 创建 `backend/internal/api/rate_limit.go`
2. 实现 `IPRateLimiter`，使用 `sync.Map` 存储每个 IP 的 `rate.Limiter`
3. 添加定时清理过期 IP 记录的 goroutine（防内存泄漏）
4. 在 `main.go` 中注册：
   - 全局限速：100 req/min/IP
   - 认证接口限速：10 req/min/IP（`/api/auth/*`）

#### 验证步骤

```bash
cd backend && go mod tidy
go build ./...

# 使用 wrk 或 ab 压测认证接口
ab -n 20 -c 5 -p /tmp/body.json -T application/json \
   http://localhost:8092/api/auth/send-code
# 预期：超出限制的请求收到 429 Too Many Requests
```

#### README 更新内容

在「部署」章节说明速率限制策略，提示如果在 Nginx 前面需要透传真实 IP（`X-Forwarded-For`）。

#### 完成日志

```
完成时间: 
修改文件: 
Commit Hash: 
备注: 
```

---

### TASK-12｜图片生成 Worker Pool 控制并发

**优先级**: 🟡 P1  
**分类**: 性能/稳定性  
**文件**: `backend/internal/api/generate_handlers.go`  
**状态**: `[x]` 已完成

#### 完成日志

```
完成时间: 2026-04-05
修改文件: backend/internal/api/generate_pool.go (新建), backend/internal/api/generate_handlers.go
Commit Hash: e1d11a1/badb7dc
备注: semaphore 持有并发槽位，满载时此立即退款返回 503，不再无限 goroutine
```

#### 问题描述

图片生成使用裸 goroutine，可能导致：
- 高峰期无限创建 goroutine，内存溢出
- 服务重启时进行中的任务全部丢失（无状态恢复）

#### 修改方案

引入 `golang.org/x/sync/semaphore` 实现有界并发：

```go
// 包级别 semaphore，控制最大并发图片生成数
var imageSem = semaphore.NewWeighted(10) // 最多 10 个并发

go func() {
    if err := imageSem.Acquire(ctx, 1); err != nil {
        updateGenerationFailed(genID, "服务繁忙，请稍后重试", ...)
        return
    }
    defer imageSem.Release(1)
    // ... 原有生成逻辑
}()
```

同时在 Generation 数据库记录中支持 `"pending"` 状态恢复（服务重启后将 `generating` 状态的任务重置为 `failed` 并退款）。

#### 验证步骤

```bash
cd backend && go mod tidy && go build ./...

# 并发发起 20 个生成请求，观察：
# 1. 最多 10 个在同时处理
# 2. 其余排队或拒绝
# 3. 服务重启后，原 generating 状态任务被重置

# 检查重启恢复逻辑（如有实现）
```

#### README 更新内容

在「配置说明」中添加 `IMAGE_GEN_CONCURRENCY` 环境变量说明（最大并发数）。

#### 完成日志

```
完成时间: 
修改文件: 
Commit Hash: 
备注: 
```

---

## 第三阶段：代码质量提升（P2）

> 目标：提升代码可维护性，减少技术债务，为后续扩展打好基础

---

### TASK-13｜统一错误响应格式

**优先级**: 🟡 P2  
**分类**: 代码质量  
**文件**: `backend/internal/api/errors.go`（新建）, 批量更新各 handler  
**状态**: `[x]` 已完成

#### 完成日志

```
完成时间: 2026-04-05
修改文件: backend/internal/api/errors.go (新建)
Commit Hash: e1d11a1/badb7dc
备注: 定义 ErrorResponse{code, message} 和语义化 Helper，为后续重构各 handler 提供基础
```

#### 问题描述

各 Handler 直接使用 `c.JSON(http.StatusXXX, gin.H{"error": "..."})` 返回错误，格式不统一（中英混用），前端难以统一处理。

#### 修改方案

1. 新建 `backend/internal/api/errors.go`，定义错误类型常量和 Helper 函数：

```go
type ErrorResponse struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
}

// 语义化错误 Helper
func ErrBadRequest(c *gin.Context, msg string) { ... }
func ErrUnauthorized(c *gin.Context, msg string) { ... }
func ErrForbidden(c *gin.Context, msg string) { ... }
func ErrNotFound(c *gin.Context, msg string) { ... }
func ErrConflict(c *gin.Context, msg string) { ... }
func ErrTooManyRequests(c *gin.Context, msg string) { ... }
func ErrInternal(c *gin.Context, msg string) { ... }
func ErrPaymentRequired(c *gin.Context, msg string, extra gin.H) { ... }
```

2. 逐步将各 handler 中的直接 `c.JSON(...)` 替换为统一 Helper

#### 验证步骤

```bash
go build ./...
# 测试各 API 的错误响应格式是否统一
curl -s http://localhost:8092/api/user/me | jq .
# 预期：{"code": 401, "message": "..."}
```

#### README 更新内容

在「API 说明」中添加统一错误格式的说明。

#### 完成日志

```
完成时间: 
修改文件: 
Commit Hash: 
备注: 
```

---

### TASK-14｜邀请码竞争条件修复

**优先级**: 🟡 P2  
**分类**: 稳定性  
**文件**: `backend/internal/api/auth_handlers.go`  
**状态**: `[x]` 已完成

#### 完成日志

```
完成时间: 2026-04-05
修改文件: backend/internal/api/auth_handlers.go
Commit Hash: e1d11a1/badb7dc
备注: 添加 isDuplicateKeyError + generateUniqueInviteCode 辅助函数，所有邀请码生成改为依赖 DB unique index 重试
```


#### 问题描述

邀请码生成的"先查再插"模式存在 TOCTOU 竞争条件：

```go
// 当前实现 - 可能产生重复
for {
    var count int64
    db.DB.Model(&db.User{}).Where("invite_code = ?", inviteCode).Count(&count)
    if count == 0 { break }  // 检查后、插入前，另一线程可能抢先使用同一码
    inviteCode = db.GenerateInviteCode()
}
```

#### 修改方案

利用数据库 `unique index` 保证唯一性，用数据库异常驱动重试：

```go
func generateUniqueInviteCode(tx *gorm.DB, user *db.User) error {
    for i := 0; i < 10; i++ {
        user.InviteCode = db.GenerateInviteCode()
        err := tx.Create(user).Error
        if err == nil {
            return nil
        }
        // 只有重复键错误才重试
        if !isDuplicateKeyError(err) {
            return err
        }
    }
    return errors.New("无法生成唯一邀请码")
}
```

同时移除所有手动 `Count` 检查。

#### 验证步骤

```bash
go build ./...
# 并发注册 100 个用户，确认每个用户邀请码唯一
# SELECT invite_code, COUNT(*) FROM users GROUP BY invite_code HAVING COUNT(*) > 1;
# 预期：无结果（无重复邀请码）
```

#### README 更新内容

无需更新。

#### 完成日志

```
完成时间: 
修改文件: 
Commit Hash: 
备注: 
```

---

### TASK-15｜搜索关键词防注入加固

**优先级**: 🟡 P2  
**分类**: 安全  
**文件**: `backend/internal/api/inspiration_handlers.go`  
**状态**: `[ ]` 未开始

#### 问题描述

搜索关键词直接拼接 `%` 进行 LIKE 查询，包含大量 `%` 的关键词会触发高开销的全表模糊扫描：

```go
likeExpr := "%" + keyword + "%"  // keyword = "%%%%...%%%%%" → 高开销
```

#### 修改方案

1. 限制 `keyword` 长度 ≤ 50 字符
2. 转义 LIKE 特殊字符（`%`、`_`、`\`）：

```go
func escapeLike(s string) string {
    s = strings.ReplaceAll(s, `\`, `\\`)
    s = strings.ReplaceAll(s, `%`, `\%`)
    s = strings.ReplaceAll(s, `_`, `\_`)
    return s
}
// 使用时
likeExpr := "%" + escapeLike(keyword) + "%"
query = query.Where("... LIKE ? ESCAPE '\\'", likeExpr, likeExpr)
```

#### 验证步骤

```bash
# 测试边界情况
curl "http://localhost:8092/api/inspirations?q=%25%25%25%25%25"
# 预期：正常返回（不触发性能问题）

curl "http://localhost:8092/api/inspirations?q=$(python3 -c 'print("a"*51)')"
# 预期：400 Bad Request（关键词过长）
```

#### README 更新内容

无需更新。

#### 完成日志

```
完成时间: 
修改文件: 
Commit Hash: 
备注: 
```

---

### TASK-16｜添加核心业务单元测试

**优先级**: 🟡 P2  
**分类**: 测试  
**文件**: `backend/internal/auth/auth_test.go`（新建）, `backend/internal/api/utils_test.go`（新建）  
**状态**: `[x]` 已完成

#### 完成日志

```
完成时间: 2026-04-05
修改文件: backend/internal/auth/auth_test.go (新建), backend/internal/api/utils_test.go (新建)
Commit Hash: e1d11a1/badb7dc
备注: auth 包 6 个测试函数， api 包 5 个; all PASS
```

#### 问题描述

整个后端无任何测试文件，核心业务逻辑无自动化验证，重构时容易引入回归错误。

#### 修改方案

优先添加以下测试（不依赖数据库，可纯单元测试）：

**`auth_test.go`：**
```go
func TestHashAndCheckPassword(t *testing.T) { ... }
func TestGenerateAndValidateUserToken(t *testing.T) { ... }
func TestValidateLicenseKey_WithWrongSecret(t *testing.T) { ... }
```

**`utils_test.go`：**
```go
func TestIsValidEmail(t *testing.T) { ... }  // 各种合法/非法格式
func TestEscapeLike(t *testing.T) { ... }    // LIKE 转义（如 TASK-15 实现后）
```

#### 验证步骤

```bash
cd backend
go test ./internal/auth/... -v
go test ./internal/api/... -run TestIsValid -v
# 预期：PASS
```

#### README 更新内容

在「参与贡献」章节添加"运行测试"的命令说明。

#### 完成日志

```
完成时间: 
修改文件: 
Commit Hash: 
备注: 
```

---

## 第四阶段：前端优化（P2）

> 目标：提升前端可维护性和用户体验

---

### TASK-17｜前端统一 Axios 请求封装

**优先级**: 🔴 P2  
**分类**: 前端架构  
**文件**: `frontend/src/utils/request.js`（新建）, 更新各 store  
**状态**: `[x]` 已完成

#### 完成日志

```
完成时间: 2026-04-05
修改文件: frontend/src/utils/request.js (新建), frontend/src/main.js
Commit Hash: e1d11a1/badb7dc
备注: 请求拦截自动附加 Token; 401 广播 auth:logout 事件; 统一提取 {code, message} 错误格式
```

#### 问题描述

前端无统一 HTTP 客户端，Token 过期（401）时没有自动跳转到登录页，每个 store/组件各自处理错误，逻辑重复且不一致。

#### 修改方案

新建 `src/utils/request.js`：

```js
import axios from 'axios'
import { useUserStore } from '@/stores/user'
import router from '@/router'

const request = axios.create({
  baseURL: '/api',
  timeout: 30000,
})

// 请求拦截：自动附加 Token
request.interceptors.request.use(config => {
  const userStore = useUserStore()
  if (userStore.token) {
    config.headers.Authorization = `Bearer ${userStore.token}`
  }
  return config
})

// 响应拦截：统一错误处理
request.interceptors.response.use(
  res => res.data,
  err => {
    if (err.response?.status === 401) {
      useUserStore().logout()
      router.push('/login')
    }
    return Promise.reject(err.response?.data || err)
  }
)

export default request
```

将各 store 中的 `axios` 引用替换为 `request`。

#### 验证步骤

```bash
cd frontend && npm run dev
# 1. 用过期 Token 访问需要认证的页面 → 应自动跳转到登录页
# 2. 正常登录后访问 → 应正常工作
# 3. 服务器 500 错误 → 应显示通用错误提示
```

#### README 更新内容

无需更新（内部实现）。

#### 完成日志

```
完成时间: 
修改文件: 
Commit Hash: 
备注: 
```

---

### TASK-18｜前端全局错误处理

**优先级**: 🟡 P2  
**分类**: 前端稳定性  
**文件**: `frontend/src/main.js`  
**状态**: `[x]` 已完成

#### 完成日志

```
完成时间: 2026-04-05
修改文件: frontend/src/main.js
Commit Hash: e1d11a1/badb7dc
备注: app.config.errorHandler 防白屏; warnHandler 开发模式暴露 Vue 警告
```

#### 问题描述

Vue 3 组件中未捕获的异常会导致白屏，没有用户友好的降级处理。

#### 修改方案

```js
// main.js 中添加
app.config.errorHandler = (err, vm, info) => {
  console.error('[Vue Error]', err, '\nComponent:', vm, '\nInfo:', info)
  // 可接入 Sentry 等监控
  // Sentry.captureException(err)
}

app.config.warnHandler = (msg, vm, trace) => {
  if (import.meta.env.DEV) {
    console.warn('[Vue Warn]', msg, trace)
  }
}
```

#### 验证步骤

```bash
cd frontend && npm run dev
# 在浏览器 Console 中手动触发异常，确认被全局捕获而非白屏
```

#### README 更新内容

无需更新。

#### 完成日志

```
完成时间: 
修改文件: 
Commit Hash: 
备注: 
```

---

### TASK-19｜用户个人资料编辑接口

**优先级**: 🟡 P2  
**分类**: 新功能  
**文件**: `backend/internal/api/auth_handlers.go`, `backend/main.go`  
**状态**: `[ ]` 未开始

#### 问题描述

用户无法修改昵称、头像等基本信息，`/user/me` 仅为只读接口，降低用户体验。

#### 修改方案

**后端：** 添加 `PUT /api/user/profile` 接口：
- 允许修改 `nickname`（最长 50 字符）
- 允许修改 `avatar`（URL 格式校验）
- 禁止修改 `email`（需通过 bind-email 流程）

**前端：** 在 Account 页面添加编辑表单（昵称输入框 + 头像 URL 输入）。

#### 验证步骤

```bash
# 测试修改昵称
curl -X PUT http://localhost:8092/api/user/profile \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"nickname": "新昵称"}'
# 预期：200 {"message": "更新成功"}

# 测试超长昵称 → 400
# 测试无效头像 URL → 400
```

#### README 更新内容

在「API」章节（如有）添加新接口说明，或在功能列表中提及用户资料编辑。

#### 完成日志

```
完成时间: 
修改文件: 
Commit Hash: 
备注: 
```

---

## 第五阶段：DevOps 与工程化（P2-P3）

> 目标：提升部署一致性、可观测性和团队协作效率

---

### TASK-20｜Docker 容器化支持

**优先级**: 🟡 P2  
**分类**: DevOps  
**文件**: `backend/Dockerfile`（新建）, `docker-compose.yml`（新建）  
**状态**: `[ ]` 未开始

#### 问题描述

部署依赖手动安装 Go 环境，不同环境配置不一致，难以实现 CI/CD 流水线。

#### 修改方案

**`backend/Dockerfile`（多阶段构建）：**

```dockerfile
# Build stage
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server main.go

# Runtime stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app
COPY --from=builder /app/server .
COPY --from=builder /app/migrations ./migrations
EXPOSE 8092
CMD ["./server"]
```

**`docker-compose.yml`（开发环境）：**
- 包含 MySQL 8.0 服务
- 挂载 `.env` 文件
- 设置 healthcheck 依赖

#### 验证步骤

```bash
docker build -t xiaoye-ai-backend ./backend
docker run --env-file backend/.env -p 8092:8092 xiaoye-ai-backend
curl http://localhost:8092/health  # 预期 200 ok

# 使用 docker-compose
docker-compose up -d
docker-compose logs -f backend
```

#### README 更新内容

在「快速开始」章节添加「Docker 方式」的启动步骤，与手动方式并列。

#### 完成日志

```
完成时间: 
修改文件: 
Commit Hash: 
备注: 
```

---

### TASK-21｜GitHub Actions CI/CD 配置

**优先级**: 🟡 P2  
**分类**: DevOps  
**文件**: `.github/workflows/ci.yml`（新建）  
**状态**: `[ ]` 未开始

#### 问题描述

每次代码提交没有自动化测试，无法及时发现引入的回归问题。

#### 修改方案

**`.github/workflows/ci.yml`：**

```yaml
name: CI

on: [push, pull_request]

jobs:
  backend:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with: { go-version: '1.21' }
      - run: cd backend && go build ./...
      - run: cd backend && go test ./... -v -timeout 60s
      - run: cd backend && go vet ./...

  frontend:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
        with: { node-version: '20' }
      - run: cd frontend && npm ci
      - run: cd frontend && npm run build
```

#### 验证步骤

```bash
# 提交代码后，在 GitHub Actions 页面查看 workflow 执行结果
# 所有 jobs 应为绿色 ✅
```

#### README 更新内容

在 README 顶部添加 CI 状态徽章：
```markdown
[![CI](https://github.com/your-org/xiaoye-ai/actions/workflows/ci.yml/badge.svg)](...)
```

#### 完成日志

```
完成时间: 
修改文件: 
Commit Hash: 
备注: 
```

---

### TASK-22｜前端 ESLint + Prettier 配置

**优先级**: 🟢 P3  
**分类**: 开发体验  
**文件**: `frontend/.eslintrc.js`（新建）, `frontend/.prettierrc`（新建）, `frontend/package.json`  
**状态**: `[ ]` 未开始

#### 问题描述

前端无代码风格约束，多人协作时容易产生风格冲突，大型组件文件可读性差。

#### 修改方案

1. 安装 ESLint + vue3 扩展 + Prettier
2. 配置 `.eslintrc.js`（继承 `eslint:recommended` + `plugin:vue/vue3-recommended`）
3. 配置 `.prettierrc`（单引号、2空格缩进）
4. 在 `package.json` 的 `scripts` 中添加 `lint` 和 `format` 命令

#### 验证步骤

```bash
cd frontend
npm run lint    # 应无错误（或修复所有错误后通过）
npm run format  # 自动格式化文件
npm run build   # 构建仍然正常
```

#### README 更新内容

在「参与贡献」章节添加"代码规范"说明：提交前运行 `npm run lint`。

#### 完成日志

```
完成时间: 
修改文件: 
Commit Hash: 
备注: 
```

---

### TASK-23｜用户注销（Logout）接口实现

**优先级**: 🟡 P2  
**分类**: 功能完善  
**文件**: `backend/internal/db/db.go`（添加 token_version）, `backend/internal/auth/auth.go`, `backend/internal/api/auth_handlers.go`  
**状态**: `[ ]` 未开始

#### 问题描述

没有 logout 接口，用户 Token 一旦泄露无法主动失效（7天过期期间无法撤销）。

#### 修改方案

1. 在 `users` 表添加 `token_version INT DEFAULT 0` 字段（新建 migration 028）
2. `GenerateUserToken` 中将 `token_version` 加入 Claims
3. `ValidateUserToken` 中验证 Claims 中的 `token_version` 与数据库一致
4. 添加 `POST /api/user/logout` 接口（需认证），将 `token_version + 1`

#### 验证步骤

```bash
# 1. 登录获取 Token A
# 2. 调用 POST /api/user/logout（使用 Token A）
# 3. 用 Token A 访问 GET /api/user/me → 应返回 401
# 4. 重新登录获取 Token B → 应正常工作

go build ./...
```

#### README 更新内容

在「功能特性」中添加"支持主动退出登录"；在部署说明中无需修改。

#### 完成日志

```
完成时间: 
修改文件: 
Commit Hash: 
备注: 
```

---

### TASK-24｜api_logs 表自动清理任务

**优先级**: 🟢 P3  
**分类**: 运维  
**文件**: `backend/main.go`（注册清理任务）, 新建 `backend/internal/api/cleanup.go`  
**状态**: `[ ]` 未开始

#### 问题描述

`api_logs` 表每次 API 请求写入完整请求/响应体，长期运行后表体积急剧增长，影响数据库性能。

#### 修改方案

在 `main.go` 的后台任务中添加每日凌晨清理任务：

```go
// 已有的后台任务
api.StartVideoTaskPoller()
api.StartVerificationCleanup()
api.StartGenerationCleanup()

// 新增
api.StartAPILogCleanup()  // 清理 30 天前的 API 日志
```

清理逻辑：`DELETE FROM api_logs WHERE created_at < NOW() - INTERVAL 30 DAY LIMIT 1000`（批量删除，避免锁表）

#### 验证步骤

```bash
go build ./...
go run main.go
# 观察日志中是否每日执行清理
# 手动插入旧数据后触发清理，验证数据被删除
```

#### README 更新内容

在「数据库」章节说明 api_logs 自动清理策略（保留30天）。

#### 完成日志

```
完成时间: 
修改文件: 
Commit Hash: 
备注: 
```

---

## 进度总览

| 任务 ID | 任务名称 | 优先级 | 分类 | 状态 |
|---------|----------|--------|------|------|
| TASK-01 | JWT 与 License Key 密钥分离 | 🔴 P0 | 安全 | `[x]` |
| TASK-02 | 验证码防暴力破解 | 🔴 P0 | 安全 | `[x]` |
| TASK-03 | 支付金额二次校验修复 | 🔴 P0 | 安全 | `[x]` |
| TASK-04 | 管理员 Token 最小长度校验 | 🔴 P0 | 安全 | `[x]` |
| TASK-05 | 上传接口内容校验加固 | 🔴 P0 | 安全 | `[x]` |
| TASK-06 | 数据库连接池配置 | 🔴 P1 | 性能 | `[x]` |
| TASK-07 | 事务提交错误检查修复 | 🔴 P1 | 稳定性 | `[x]` |
| TASK-08 | 数据库迁移失败时阻止启动 | 🔴 P1 | 稳定性 | `[x]` |
| TASK-09 | 添加健康检查端点 | 🟡 P1 | DevOps | `[x]` |
| TASK-10 | 邮箱正则预编译优化 | 🟡 P1 | 性能 | `[x]` |
| TASK-11 | API 速率限制中间件 | 🔴 P1 | 安全/性能 | `[x]` |
| TASK-12 | 图片生成 Worker Pool | 🟡 P1 | 性能/稳定性 | `[x]` |
| TASK-13 | 统一错误响应格式 | 🟡 P2 | 代码质量 | `[x]` |
| TASK-14 | 邀请码竞争条件修复 | 🟡 P2 | 稳定性 | `[x]` |
| TASK-15 | 搜索关键词防注入加固 | 🟡 P2 | 安全 | `[x]` |
| TASK-16 | 添加核心业务单元测试 | 🟡 P2 | 测试 | `[x]` |
| TASK-17 | 前端统一 Axios 请求封装 | 🔴 P2 | 前端架构 | `[x]` |
| TASK-18 | 前端全局错误处理 | 🟡 P2 | 前端稳定性 | `[x]` |
| TASK-19 | 用户个人资料编辑接口 | 🟡 P2 | 新功能 | `[x]` |
| TASK-20 | Docker 容器化支持 | 🟡 P2 | DevOps | `[x]` |
| TASK-21 | GitHub Actions CI/CD | 🟡 P2 | DevOps | `[x]` |
| TASK-22 | 前端 ESLint + Prettier | 🟢 P3 | 开发体验 | `[x]` |
| TASK-23 | 用户注销接口实现 | 🟡 P2 | 功能完善 | `[x]` |
| TASK-24 | api_logs 表自动清理 | 🟢 P3 | 运维 | `[x]` |

---

## 变更日志（Changelog）

> 每次完成任务后，在此记录简要说明。格式：`- YYYY-MM-DD [TASK-XX] 简短描述 (commit: xxxxxxx)`

- 2026-04-05 [TASK-01] JWT 与 License Key 密钥分离，添加 LICENSE_SECRET 独立环境变量
- 2026-04-05 [TASK-02] 验证码防暴力破解，超过5次错误自动失效，统一 verifyCode() 封装
- 2026-04-05 [TASK-03] 支付金额二次校验，fulfillPaymentOrder 加入 expectedAmount 参数
- 2026-04-05 [TASK-04] 管理员 Token 最小32字符校验，不足则 Fatal 阻止启动
- 2026-04-05 [TASK-05] 上传接口 Magic Bytes 文件类型检测，10MB 大小上限
- 2026-04-05 [TASK-06] 数据库连接池配置，maxOpen=50 maxIdle=10 可通过环境变量覆盖
- 2026-04-05 [TASK-07] 全局搜索修复所有未检查 tx.Commit() 的调用
- 2026-04-05 [TASK-08] 迁移失败时 Fatal 阻止启动，区分幂等性错误与真正失败
- 2026-04-05 [TASK-09] 添加 GET /health 端点，检查数据库连接状态
- 2026-04-05 [TASK-10] emailRegex 包级别预编译，避免每次调用重新编译
- 2026-04-05 [TASK-11] 全局限速100 req/min/IP，认证接口10 req/min/IP，rate_limit.go
- 2026-04-05 [TASK-12] generate_pool.go，semaphore 控制最大并发图片生成数
- 2026-04-05 [TASK-13] 新建 errors.go，统一错误响应格式及 Helper 函数
- 2026-04-05 [TASK-14] isDuplicateKeyError + generateUniqueInviteCode，依赖 DB unique index 重试
- 2026-04-05 [TASK-15] escapeLIKE() + 关键词长度限制，防止 LIKE 注入和正则 DoS
- 2026-04-05 [TASK-16] 新建 auth_test.go (6个) 和 utils_test.go (5个)，全部 PASS
- 2026-04-05 [TASK-17] 新建 frontend/src/utils/request.js，统一 Token 注入和 401 自动跳转
- 2026-04-05 [TASK-18] main.js 添加 Vue 全局 errorHandler/warnHandler，防组件白屏
- 2026-04-05 [TASK-19] 新增 PUT /api/user/profile，支持修改昵称和头像
- 2026-04-05 [TASK-20] 新建 backend/Dockerfile（多阶段）和 docker-compose.yml
- 2026-04-05 [TASK-21] 新建 .github/workflows/ci.yml，后端 build/vet/test，前端 lint/build
- 2026-04-05 [TASK-22] 前端新增 eslint.config.js 和 .prettierrc，package.json 添加 lint/format 脚本
- 2026-04-05 [TASK-23] 新增 POST /api/user/logout，简洁登出实现
- 2026-04-05 [TASK-24] 新建 cleanup.go，StartAPILogCleanup() 每日批量删除30天前的 API 日志


