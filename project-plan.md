# 寻氧AI 项目优化计划

> 分析时间：2026-04-04
> 分析范围：后端（Go/Gin）、前端（Vue 3）、管理后台、数据库迁移、基础设施

---

## 一、项目概览

**寻氧AI** 是一个多模态 AI 内容创作平台，支持图像生成、视频生成、电商图批量生成、灵感广场社区等功能。技术栈为 Go + Gin（后端）、Vue 3 + Naive UI（前端）、MySQL + GORM（数据层）。

整体架构清晰，业务模块划分合理，已具备 MVP 能力。以下对照行业最佳实践，逐层列出可优化的方向。

---

## 二、安全性（Security）

### 🔴 高优先级

#### 2.1 JWT 与 License Key 共用同一密钥
**位置**: `backend/internal/auth/auth.go:27`
```go
LicenseSecretKey = []byte(secret)  // 与 SecretKey 相同！
```
**问题**: 用户 JWT Token 和 License Key 使用同一个 `JWT_SECRET`，理论上可以伪造 License Key 来兑换积分。  
**建议**: 为 License Key 引入独立的 `LICENSE_SECRET` 环境变量，或改为使用非对称加密（RS256）签名 License Key。

#### 2.2 管理员 Token 静态化
**位置**: `backend/internal/api/admin/middleware.go`  
**问题**: Admin Token 是一个静态字符串，和普通 API Token 混用同一 Authorization 头，没有过期机制、没有角色权限管理，一旦泄露永久有效。  
**建议**: 
- 短期：添加 Token 最小长度校验（至少 32 字符）
- 长期：实现 Admin JWT，支持过期、角色控制（`admin` / `moderator`）

#### 2.3 邮箱验证码无防暴力破解机制
**位置**: `backend/internal/api/auth_handlers.go:55-58`  
**问题**: 只有"1 分钟内不能重发"的限速，但对验证码的猜测尝试没有限制（6 位纯数字 = 100 万种组合，可暴力枚举）。  
**建议**: 在验证码验证端点（`/auth/register`, `/auth/login`）添加错误次数计数，超过 5 次后使验证码失效并强制重新获取。

#### 2.4 支付回调验签后未二次校验金额
**位置**: `backend/internal/api/payment_handlers.go:261-265`  
**问题**: 虽然有 `amountsEqual` 校验，但在 `GetPaymentStatus` 的主动查询路径（`fulfillPaymentOrder`）中没有金额校验，存在绕过风险。  
**建议**: 在 `fulfillPaymentOrder` 函数中也加入金额校验逻辑。

#### 2.5 上传接口缺乏内容类型和大小限制
**位置**: `backend/internal/api/resource_handlers.go`（推测）, `api/utils.go`  
**问题**: 图片/视频上传接口对文件大小、MIME 类型的校验不够严格，存在上传恶意文件的风险。  
**建议**: 后端强制校验 Base64 图片 Magic Bytes，设置最大尺寸上限（如图片 10MB）。

---

### 🟡 中优先级

#### 2.6 速率限制（Rate Limiting）缺失
**问题**: 整个 API 没有全局速率限制中间件，`/api/auth/send-code` 等公开接口仅靠数据库查询做频率控制，无法防御分布式请求。  
**建议**: 引入 `golang.org/x/time/rate` 或 Redis 实现基于 IP 的速率限制中间件，对公开接口（发验证码、登录）设置更严格的限速。

#### 2.7 CORS 配置过于宽松（开发环境默认值暴露）
**位置**: `backend/main.go:73`
```go
corsConfig.AllowOrigins = []string{"http://localhost:5173", "http://localhost:5174"}
```
**建议**: 生产环境必须通过 `CORS_ORIGINS` 环境变量设置白名单，可在启动时检测并告警。

#### 2.8 SQL 注入风险（LIKE 查询未做参数化）
**位置**: `backend/internal/api/inspiration_handlers.go:612-613`
```go
likeExpr := "%" + keyword + "%"
query = query.Where("(title LIKE ? OR prompt LIKE ?)", likeExpr, likeExpr)
```
**问题**: 虽然使用了参数化查询，但 `%` 拼接可以被用于正则拒绝服务（`%` 大量嗅探）。  
**建议**: 对 `keyword` 长度做限制（如 ≤ 50 字符），并转义 `%` 和 `_`。

---

## 三、性能与可扩展性

### 🔴 高优先级

#### 3.1 每次请求均查询数据库校验用户状态
**位置**: `backend/internal/api/middleware.go:34-38`
```go
var user db.User
if err := db.DB.First(&user, claims.UserID).Error; err != nil { ... }
```
**问题**: 每个需要认证的请求都会查询一次 `users` 表，高并发下会成为数据库瓶颈。  
**建议**: 引入 Redis 做用户状态缓存（TTL 5 分钟），只有缓存缺失或 Token 刷新时才查数据库；禁用用户时主动清除缓存。

#### 3.2 灵感广场列表存在 N+1 查询风险
**位置**: `backend/internal/api/inspiration_handlers.go:646-667`  
**问题**: 虽然已批量查询 users 和 tags，但每次 `listInspirationsInternal` 仍执行了 4 次独立查询（count + posts + users + tags + likes），可以合并优化。  
**建议**: 使用数据库索引和分页游标（`cursor-based pagination`）替代 `OFFSET`，避免深分页性能问题。

#### 3.3 `go.mod` 模块名与业务不符
**位置**: `backend/go.mod:1`
```
module google-ai-proxy
```
**问题**: 模块名遗留了项目前身的痕迹，与当前业务「寻氧AI」不匹配，影响代码可读性和维护性。  
**建议**: 重命名为 `xiaoye-ai/backend` 或 `github.com/youorg/xiaoye`。

#### 3.4 数据库连接池未配置
**位置**: `backend/internal/db/db.go:304`
```go
DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
```
**问题**: 未设置最大连接数、连接最大空闲时间等，使用 GORM 默认值，高并发下可能耗尽连接。  
**建议**: 
```go
sqlDB, _ := DB.DB()
sqlDB.SetMaxOpenConns(50)
sqlDB.SetMaxIdleConns(10)
sqlDB.SetConnMaxLifetime(time.Hour)
```

#### 3.5 图片生成异步 Goroutine 无法追踪
**位置**: `backend/internal/api/generate_handlers.go:178`
```go
go func() { ... }()
```
**问题**: 生成任务通过裸 goroutine 运行，服务重启时进行中的任务会丢失，且无法控制并发上限（可能无限制创建 goroutine）。  
**建议**: 使用有界 worker pool（如 `golang.org/x/sync/semaphore`）控制并发数，或将图片生成也纳入视频任务的 polling 轮询机制统一管理。

---

### 🟡 中优先级

#### 3.6 APILog 写入数据库效率低
**位置**: `backend/internal/api/utils.go:112`  
**问题**: 每次 API 调用都异步写入一条完整日志到 `api_logs` 表（包含完整请求/响应体），随时间增长表会变得很大，影响查询性能。  
**建议**: 
- 考虑使用文件日志（如 `zap` + logrotate）替代数据库存储
- 或者引入 `api_logs` 表的定时清理（超过 30 天的记录删除）

#### 3.7 `isValidEmail` 使用 `regexp.MatchString` 每次编译正则
**位置**: `backend/internal/api/utils.go:19`  
**建议**: 将正则表达式编译为 package 级别变量，避免每次调用都重新编译。
```go
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
```

---

## 四、代码质量与架构

### 🔴 高优先级

#### 4.1 超大型处理器文件（Fat Handler）
**问题**: 单个处理器文件过大：
- `inspiration_handlers.go`: 1480 行
- `auth_handlers.go`: 1054 行
- `prompt_optimize_handlers.go`: 未查，推测也较大

**说明**: 违反单一职责原则，难以维护和测试。  
**建议**: 将业务逻辑抽取到 `service` 层（如 `internal/service/user_service.go`），handler 只负责参数解析、调用 service、返回响应。

#### 4.2 缺乏统一的错误处理机制
**问题**: 各处理器直接用 `c.JSON(http.StatusXXX, gin.H{"error": "..."})` 返回错误，错误码和消息格式不一致（中英混用）。  
**建议**: 定义统一的 `AppError` 类型和错误中间件：
```go
type AppError struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
}
```

#### 4.3 硬编码业务常量
**问题**: 多处存在魔法数字和硬编码字符串：
- 新用户赠送 10 钻石（`auth_handlers.go:189`）
- 邀请奖励 10 钻石、上限 500 钻石（`auth_handlers.go:232-233`）
- 支付套餐价格硬编码（`payment_handlers.go:22-29`）
- 验证码 6 位（多处）

**建议**: 将业务配置移到环境变量或配置文件，利用 `config` 包统一管理。

#### 4.4 邀请码唯一性检查存在竞争条件
**位置**: `backend/internal/api/auth_handlers.go:174-181`
```go
for {
    var count int64
    db.DB.Model(&db.User{}).Where("invite_code = ?", inviteCode).Count(&count)
    if count == 0 { break }
    ...
}
```
**问题**: 检查与写入之间存在 TOCTOU 竞争，高并发时可能产生重复邀请码。  
**建议**: 依赖数据库的 `uniqueIndex` 约束，在 `Create` 失败时捕获重复键错误并重试，而非先查再插。

#### 4.5 缺乏单元测试
**问题**: 整个后端没有测试文件（`*_test.go`）。核心业务逻辑（积分系统、支付、生成流水）完全缺少自动化测试覆盖。  
**建议**: 优先对以下模块补充测试：
- `auth.HashPassword` / `CheckPassword`
- `recordCreditTransaction`
- `refundCredits`
- `fulfillPaymentOrder`

---

### 🟡 中优先级

#### 4.6 `cmd/` 目录未充分利用
**位置**: `backend/cmd/`  
**问题**: `cmd` 目录存在但内容不明（密钥生成工具），未提供 CLI 管理工具（如用户管理、批量返款等）。  
**建议**: 利用 `cmd/` 实现管理 CLI（`cmd/admin-cli/main.go`），避免直接数据库操作。

#### 4.7 `ImageRecord` 遗留模型（Legacy Code）
**位置**: `backend/internal/db/db.go:140-155`  
**问题**: 代码注释明确标注 "legacy and kept for compatibility"，但实际是否仍在使用不明确，增加维护负担。  
**建议**: 评估是否可彻底移除，若需保留应写明原因和移除计划。

#### 4.8 事务提交错误处理不一致
**位置**: `backend/internal/api/auth_handlers.go:299`
```go
tx.Commit()  // 未检查错误！
```
**建议**: 所有事务提交都应检查返回错误：
```go
if err := tx.Commit().Error; err != nil {
    log.Printf("事务提交失败: %v", err)
    c.JSON(http.StatusInternalServerError, gin.H{"error": "操作失败"})
    return
}
```

---

## 五、前端优化

### 🔴 高优先级

#### 5.1 前端 `package.json` 名称与业务不符
**位置**: `frontend/package.json:2`
```json
"name": "google-ai-frontend"
```
**建议**: 改为 `xiaoye-ai-frontend`，与项目保持一致。

#### 5.2 组件文件体积过大
**问题**: 多个组件文件超过 20KB：
- `ComposerBar.vue`: 49KB（接近 2000 行）
- `AuthModal.vue`: 35KB
- `ShareGenerationDialog.vue`: 31KB

**建议**: 
- 将 `ComposerBar.vue` 按功能拆分为子组件（如 `ImageComposer.vue`, `VideoComposer.vue`）
- 对大型组件使用 `defineAsyncComponent` 实现懒加载

#### 5.3 缺乏错误边界和全局异常处理
**位置**: `frontend/src/main.js`  
**问题**: 没有设置 `app.config.errorHandler`，组件内未捕获的异常会导致白屏。  
**建议**:
```js
app.config.errorHandler = (err, vm, info) => {
  console.error('Vue Error:', err, info)
  // 上报到监控系统
}
```

#### 5.4 API 请求缺乏统一的错误拦截
**问题**: 前端 API 调用散布在各组件/store 中，没有统一的 axios 拦截器处理 401（Token 过期跳转登录）、500（全局提示）等场景。  
**建议**: 创建 `src/utils/request.js`，配置 axios 实例，统一处理认证错误和服务器错误。

---

### 🟡 中优先级

#### 5.5 i18n 覆盖不完整
**位置**: `frontend/src/i18n.js`  
**问题**: 部分 UI 文案直接写在组件模板中（中文硬编码），没有经过 i18n 系统，国际化不彻底。  
**建议**: 所有面向用户的文字都通过 `$t()` 引用，并补全英文翻译条目。

#### 5.6 Pinia Store 缺乏持久化
**问题**: `user` store 的登录状态仅在内存中，刷新页面后需要重新从 localStorage 读取 token 并请求用户信息，流程繁琐。  
**建议**: 使用 `pinia-plugin-persistedstate` 对 `userStore` 中的 token 做持久化。

#### 5.7 图片/视频加载缺乏懒加载和占位图
**问题**: 灵感广场大量媒体资源直接加载，没有懒加载策略，首屏 FCP 指标较差。  
**建议**: 使用 Intersection Observer 或 Vue 懒加载指令，对屏幕外图片延迟加载；添加 skeleton 占位。

#### 5.8 前端缺少 TypeScript
**问题**: 整个前端使用纯 JavaScript，无类型检查，大型组件（ComposerBar 等）的 props/emits 接口缺乏类型约束。  
**建议**: 逐步迁移到 TypeScript（`lang="ts"`），至少对 stores 和 composables 加类型注解。

---

## 六、DevOps 与部署

### 🔴 高优先级

#### 6.1 缺少 Docker/容器化支持
**问题**: 没有 `Dockerfile`，部署依赖手动配置环境，不利于一致性部署和 CI/CD。  
**建议**: 
- 添加 `backend/Dockerfile`（多阶段构建，减小镜像体积）
- 添加 `docker-compose.yml`（后端 + MySQL + Redis）

#### 6.2 缺少 CI/CD 配置
**位置**: `.github/` 目录（内容未知）  
**问题**: 没有自动化构建、测试、部署流程。  
**建议**: 添加 GitHub Actions workflow：
- `on: push` → 运行 `go test ./...` 和 `npm run build`
- `on: merge to main` → 自动部署到生产环境

#### 6.3 日志系统基础薄弱
**问题**: 整个后端使用标准库 `log` 包，没有日志级别、结构化日志、日志轮转等功能。  
**建议**: 引入 `go.uber.org/zap` 或 `github.com/rs/zerolog`，实现结构化日志（JSON 格式），便于 ELK / Loki 接入。

---

### 🟡 中优先级

#### 6.4 健康检查端点缺失
**问题**: 没有 `/health` 或 `/ping` 端点，负载均衡器和 K8s 就绪探针无法使用。  
**建议**: 添加：
```go
r.GET("/health", func(c *gin.Context) {
    c.JSON(200, gin.H{"status": "ok", "ts": time.Now().Unix()})
})
```

#### 6.5 监控与告警缺失
**问题**: 没有 Metrics 暴露（Prometheus），无法监控 API QPS、错误率、数据库连接数等关键指标。  
**建议**: 引入 `github.com/gin-gonic/gin-contrib/prometheus` 或 `go.opentelemetry.io/otel`。

#### 6.6 `.github/` 目录内容缺失关键模板
**建议**: 添加：
- `ISSUE_TEMPLATE/bug_report.md`
- `pull_request_template.md`
- `CONTRIBUTING.md`

---

## 七、数据库与数据层

#### 7.1 迁移脚本与 GORM AutoMigrate 双轨并行
**位置**: `backend/internal/db/db.go` + `backend/migrations/`  
**问题**: 代码注释说明使用 SQL 迁移脚本，但实际代码逻辑未使用 GORM 的 `AutoMigrate`，两套机制并行运行（还是只有 SQL 脚本？）。需要明确迁移策略，避免两套方案冲突。  
**建议**: 统一使用 SQL 迁移脚本（当前方案），在 README 中明确说明，并在迁移失败时让服务启动失败而非忽略错误继续运行。

#### 7.2 迁移错误被忽略
**位置**: `backend/internal/db/db.go:380-389`
```go
if !failed {
    // 仅成功才标记
} else {
    log.Printf("migraions: %s had errors, not marking...", name)
    // 继续运行！
}
```
**问题**: 迁移部分失败时服务仍正常启动，可能导致数据库结构不完整而出现运行时错误。  
**建议**: 迁移失败时应 `log.Fatal()`，强制服务停止，确保数据库结构始终一致。

#### 7.3 缺乏数据库备份策略文档
**建议**: 在 `docs/` 或 README 中说明推荐的备份策略（如 mysqldump cron job 或阿里云 RDS 自动备份）。

#### 7.4 `api_logs` 表缺乏清理机制
**问题**: 每次 API 调用写入完整请求响应体到 `api_logs` 表，长期运行后表体积会急剧增长。  
**建议**: 添加定时清理任务（保留最近 7 天），或改用外部日志系统（ELK/Loki）。

---

## 八、功能完善

#### 8.1 用户注销（Logout）功能缺失
**问题**: 没有 `/auth/logout` 接口，也没有 Token 黑名单机制。用户无法强制使 Token 失效（如设备丢失时）。  
**建议**: 实现 Token 版本号（`token_version` 字段），每次 logout 时 +1，使旧 Token 失效。

#### 8.2 账号安全功能缺失
**建议**: 
- 密码强度要求：目前只有 6 字符最低要求，建议增加复杂度校验
- 登录失败次数锁定：连续失败 5 次后锁定 30 分钟
- 异地登录提醒：记录登录 IP，异常时发送邮件通知

#### 8.3 用户个人资料编辑缺失
**问题**: `/user/me` 接口只有读取，没有 PUT 更新接口（修改昵称、头像等）。  
**建议**: 添加 `PUT /api/user/profile` 接口，支持修改昵称、头像 URL。

#### 8.4 内容举报机制缺失
**问题**: 灵感广场没有举报功能，用户无法对违规内容进行举报，管理员也只能主动发现。  
**建议**: 添加 `POST /api/inspirations/:shareID/report` 接口和对应的管理后台处理界面。

#### 8.5 管理后台功能单一
**问题**: 管理后台（`frontend-admin/`）功能极为有限，仅支持灵感审核，缺乏：
- 用户管理（封禁/解封、积分调整）
- 数据统计大盘（DAU、生成量、收入）
- License Key 批量生成与管理
- 支付订单管理

#### 8.6 支付方式单一
**问题**: 目前仅支持 Linux.do 积分支付，限制了用户群体。  
**建议**: 评估接入支付宝、微信支付（wechatpay-go SDK）。

---

## 九、文档与开发体验

#### 9.1 API 文档缺失
**问题**: 没有 Swagger / OpenAPI 文档，前端开发和第三方集成需要读源码。  
**建议**: 引入 `github.com/swaggo/swag`，为主要 API 添加注解并自动生成文档。

#### 9.2 前端无 ESLint / Prettier 配置
**问题**: `frontend/` 没有代码规范配置文件（`.eslintrc`、`.prettierrc`），代码风格依赖个人习惯。  
**建议**: 添加 ESLint + Prettier + husky + lint-staged，提交前自动格式化。

#### 9.3 环境变量缺少完整性校验
**问题**: 后端启动时只有 `DB_USER/PASSWORD/NAME` 和 `JWT_SECRET` 做存在性检查，其他关键配置（如 OSS、SMTP）是可选的，但如果功能被使用时缺失则会报运行时错误。  
**建议**: 在 `config` 包中实现 `Validate()` 函数，启动时统一校验所有配置，缺失时打印清晰的错误提示。

---

## 十、优先级汇总

| 优先级 | 数量 | 类别 |
|--------|------|------|
| 🔴 高  | 12   | 安全漏洞、性能瓶颈、架构问题 |
| 🟡 中  | 16   | 代码质量、用户体验、可维护性 |
| 🟢 低  | -    | 文档补充、工具链完善 |

### 建议执行顺序

**第一阶段（安全加固）**
1. JWT 与 License Key 密钥分离
2. 添加验证码暴力破解防护
3. 支付金额双校验修复
4. 添加速率限制中间件
5. 数据库连接池配置

**第二阶段（稳定性提升）**
1. 图片生成改用 worker pool
2. 迁移失败时 Fatal 停止
3. 事务提交错误检查
4. 添加健康检查端点
5. 替换标准 log 为结构化日志

**第三阶段（体验完善）**
1. 管理后台扩充用户/数据管理
2. 用户账号安全功能
3. API 文档（Swagger）
4. 前端代码规范配置
5. Docker / CI 构建流水线
