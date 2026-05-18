# 寻氧AI 内容创作平台

[![License: AGPL-3.0](https://img.shields.io/badge/License-AGPL--3.0-blue.svg)](LICENSE)

开源多模态 AI 内容创作平台，支持 Google Gemini、火山引擎 Seedream/Seedance 等模型进行图像和视频生成。

**在线体验: [https://xiaoye.io](https://xiaoye.io)**

[中文](#功能特性) | [English](README_EN.md)

## 截图预览

| 灵感广场                                 | 生成界面                                      |
|--------------------------------------|-------------------------------------------|
| ![生成](docs/screenshots/generate.png) | ![灵感广场](docs/screenshots/inspiration.png) |

## 功能特性

- **AI 图像生成** — 多模型切换 (Gemini, Seedream)，支持参考图，最高 4K 分辨率
- **AI 视频生成** — 文生视频 & 图生视频 (Seedance, Veo 3.1)，异步任务轮询
- **电商图片批量生成** — 模板化批量生成商品图
- **提示词优化** — 基于 DeepSeek 的 AI 提示词改写
- **反向提示词** — 从图片提取生成提示词
- **灵感广场** — 社区作品展示，点赞、Remix、审核
- **积分系统** — 密钥兑换、每日签到、邀请奖励、在线充值
- **用户系统** — 邮箱注册登录、Linux.do OAuth、邮箱绑定
- **国际化** — 中文 & 英文

## 技术栈

| 层级 | 技术 |
|------|------|
| 前端 | Vue 3, Naive UI, Pinia, Vue Router, Vite |
| 后端 | Go (Gin), GORM, MySQL, JWT |
| AI 供应商 | Google Gemini, 火山引擎 (豆包) |
| 存储 | 阿里云 OSS |
| 管理后台 | Vue 3 SPA (独立构建) |

## 项目结构

```
.
├── frontend/                # Vue 3 前端应用
│   └── src/
│       ├── views/           # 页面
│       ├── components/      # 可复用组件
│       ├── composables/     # 组合式 API hooks
│       ├── stores/          # Pinia 状态管理
│       ├── locales/         # 国际化 (zh/en)
│       └── router/
├── frontend-admin/          # 管理后台审核面板
├── backend/
│   ├── cmd/                 # CLI 工具 (密钥生成等)
│   ├── internal/
│   │   ├── api/             # HTTP 处理器 & 中间件
│   │   ├── api/admin/       # 管理后台 API
│   │   ├── auth/            # JWT 认证
│   │   ├── config/          # 环境配置
│   │   ├── db/              # 数据库模型 & 初始化
│   │   ├── email/           # SMTP 邮件服务
│   │   ├── payment/         # 支付集成
│   │   ├── provider/        # AI 供应商适配器
│   │   └── storage/         # OSS 存储
│   └── migrations/          # SQL 迁移脚本
└── LICENSE
```

## 快速开始

### 环境要求

- Go 1.21+
- Node.js 18+
- MySQL 8.0+
- 阿里云 OSS 存储桶 (用于媒体存储)
- 至少一个 AI 供应商 API Key (Google Gemini 或火山引擎)

### 后端

```bash
cd backend
cp .env.example .env    # 编辑 .env 填入你的配置
go run main.go          # 启动于 :8092
```

### 前端

```bash
cd frontend
npm install
npm run dev             # 启动于 http://localhost:5173
```

### 管理后台

```bash
cd frontend-admin
npm install
npm run dev             # 启动于 http://localhost:5174
```

## 配置说明

所有配置通过环境变量管理。复制 `backend/.env.example` 为 `backend/.env` 并填入对应值。

### 必填项

| 变量 | 说明 |
|------|------|
| `DB_USER` | MySQL 用户名 |
| `DB_PASSWORD` | MySQL 密码 |
| `DB_NAME` | MySQL 数据库名 |
| `JWT_SECRET` | JWT 签名密钥 (使用强随机字符串) |
| `GOOGLE_API_KEY` | Google Gemini API Key ([获取](https://aistudio.google.com/app/apikey)) |

### 可选项

| 变量 | 说明 | 默认值 |
|------|------|--------|
| `DB_HOST` | MySQL 主机 | `localhost` |
| `DB_PORT` | MySQL 端口 | `3306` |
| `PORT` | 后端服务端口 | `8092` |
| `HTTP_PROXY` | 外部 API 调用代理 | - |
| `ARK_API_KEY` | 火山引擎 API Key (Seedream/Seedance) | - |
| `DEEPSEEK_API_KEY` | DeepSeek API Key (提示词优化) | - |
| `DEEPSEEK_BASE_URL` | DeepSeek API 基础地址 | `https://api.deepseek.com` |
| `DEEPSEEK_MODEL` | DeepSeek 模型名 | `deepseek-chat` |
| `PROMPT_OPTIMIZE_CREDITS` | 提示词优化消耗积分 | `1` |
| `OSS_ENDPOINT` | 阿里云 OSS Endpoint | - |
| `OSS_ACCESS_KEY_ID` | OSS Access Key ID | - |
| `OSS_ACCESS_KEY_SECRET` | OSS Access Key Secret | - |
| `OSS_BUCKET_NAME` | OSS 存储桶名 | - |
| `OSS_REGION` | OSS 地域 | - |
| `OSS_PUBLIC_DOMAIN` | OSS 自定义域名 | - |
| `SMTP_HOST` | SMTP 服务器地址 | - |
| `SMTP_PORT` | SMTP 端口 | - |
| `SMTP_USER` | SMTP 用户名 | - |
| `SMTP_PASSWORD` | SMTP 密码 / 授权码 | - |
| `SMTP_FROM_EMAIL` | 发件人邮箱 | - |
| `SMTP_FROM_NAME` | 发件人名称 | - |
| `ADMIN_TOKEN` | 管理后台认证 Token | - |
| `INSPIRATION_AUTO_APPROVE` | 灵感帖子自动审核通过 | `false` |
| `CORS_ORIGINS` | CORS 允许来源 (逗号分隔) | `http://localhost:5173,http://localhost:5174` |
| `LINUXDO_CLIENT_ID` | Linux.do OAuth Client ID | - |
| `LINUXDO_CLIENT_SECRET` | Linux.do OAuth Client Secret | - |
| `OAUTH_REDIRECT_URL` | OAuth 回调地址 | - |
| `LINUXDO_CREDIT_PID` | Linux.do Credit 商户 ID | - |
| `LINUXDO_CREDIT_KEY` | Linux.do Credit 商户密钥 | - |
| `LINUXDO_CREDIT_NOTIFY_URL` | 支付回调地址 | - |
| `LINUXDO_CREDIT_RETURN_URL` | 支付完成跳转地址 | - |

## 数据库

后端使用 GORM 自动迁移，首次启动时自动创建表结构。`backend/migrations/` 目录下的 SQL 脚本供参考和手动变更使用。

## 部署

### 生产构建

```bash
# 前端
cd frontend && npm run build        # 产物: frontend/dist/

# 管理后台
cd frontend-admin && npm run build  # 产物: frontend-admin/dist/

# 后端
cd backend && go build -o server main.go
```

### Nginx 配置示例

```nginx
server {
    listen 443 ssl;
    server_name your-domain.com;

    # 前端
    location / {
        root /path/to/frontend/dist;
        try_files $uri $uri/ /index.html;
    }

    # 后端 API
    location /api/ {
        proxy_pass http://127.0.0.1:8092;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

生产环境在 `.env` 中设置 `CORS_ORIGINS=https://your-domain.com`。

## 参与贡献

1. Fork 本仓库
2. 创建功能分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add amazing feature'`)
4. 推送分支 (`git push origin feature/amazing-feature`)
5. 发起 Pull Request

## 代码架构深度解析

> 本节从代码层面深入剖析项目的设计决策与实现细节，适合二次开发者或希望深入理解系统的读者。

### 系统整体架构

```
┌─────────────────────────────────────────────────────────────────┐
│                         用户端 / 管理端                           │
│   frontend/ (Vue3+NaiveUI)    frontend-admin/ (Vue3+AntDesign)  │
└────────────────────┬────────────────────────────────────────────┘
                     │ HTTP/REST (CORS)
┌────────────────────▼────────────────────────────────────────────┐
│                      backend/  Go + Gin                          │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌────────────────┐  │
│  │   API    │  │  Auth    │  │ Payment  │  │ Background     │  │
│  │ Handlers │  │  JWT     │  │ Linux.do │  │ Workers (poller│  │
│  └────┬─────┘  └──────────┘  └──────────┘  │ cleanup tasks) │  │
│       │                                     └────────────────┘  │
│  ┌────▼──────────────────────────────────────────────────────┐  │
│  │              Provider 抽象层 (interface.go)                │  │
│  │  Gemini(图) │ Volcengine(图) │ GoogleVideo │ VolcVideo     │  │
│  └─────────────────────────────────────────────────────────┘  │
│  ┌──────────────────────────┐  ┌─────────────────────────────┐  │
│  │  db/ (GORM + MySQL)      │  │  storage/ (阿里云 OSS)       │  │
│  └──────────────────────────┘  └─────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────┘
```

### 路由层设计

`main.go` 将所有路由清晰分为四个权限层级：

| 层级 | 路径前缀 | 鉴权方式 | 限速策略 |
|------|---------|---------|---------|
| 公开接口 | `/api/pricing`, `/api/models`, `/api/inspirations` | 无 | 全局 100 req/min/IP |
| Auth 接口 | `/api/auth/*` | 无 | 严格 10 req/min/IP（防暴力破解） |
| 用户接口 | `/api/user/*`, `/api/generate`, `/api/generations` | JWT Bearer Token | 全局限速 |
| 管理接口 | `/api/admin/*` | `X-Admin-Token` 请求头 | 全局限速 |

**后台 Worker（启动时注册）：**
- `StartVideoTaskPoller()` — 轮询视频任务状态，驱动 queued → success/failed 状态机
- `StartVerificationCleanup()` — 清理过期邮箱验证码
- `StartGenerationCleanup()` — 清理长时间 pending 的生成记录
- `StartAPILogCleanup()` — 定期删除 30 天前的 API 日志

### 统一生成接口 `POST /api/generate`

整个创作能力汇聚于单一接口，通过 `type` 字段路由：

```
UnifiedGenerate (generate_handlers.go)
  ├── type="image"      → handleUnifiedImageGenerate
  │     └── 后台 goroutine（有并发槽位控制）
  ├── type="video"      → handleUnifiedVideoGenerate
  │     └── Provider.CreateVideoTask() → 任务 ID → Poller 轮询
  └── type="ecommerce"  → handleUnifiedEcommerceGenerate
        └── 后台 goroutine（支持 5~15 张批量输出）
```

**图片/电商钻石扣费流程（防超扣设计）：**

```
1. 检查余额（读取用户 credits）
2. 乐观锁扣费：UPDATE users SET credits=credits-N WHERE id=? AND credits>=N
3. 写积分流水账本（credit_transactions）
4. 创建 Generation 记录（status=generating）
5. 申请并发槽位（generate_pool.go，防 goroutine 无限累积）
6. 启动 goroutine 执行生成
   └── 任意步骤失败 → refundCredits() 退款 + 更新记录为 failed
```

### AI Provider 抽象层

`internal/provider/` 定义了两套接口，以适配图像和视频的不同调用模式：

```go
// 图像生成（同步调用）
ImageGenerator       → GenerateImage(prompt, opts) (*ImageResult, error)
MultiImageGenerator  → GenerateMultiImage(prompt, inputs, count, opts) → 电商多图

// 视频生成（异步任务）
VideoProvider        → CreateVideoTask(req) → TaskID
                    → QueryVideoTask(taskID) → 状态 + 结果 URL
                    → CalculateCredits(resolution, duration, audio) → 积分数
```

**当前接入的模型：**

| 模型 ID | 文件 | 用途 |
|---------|------|------|
| `gemini-3-pro-image-preview` | `gemini.go` | 图像生成（默认首选） |
| `gemini-3.1-flash-image-preview` | `gemini.go` | 图像生成（Flash 版本） |
| `doubao-seedream-4-5` | `volcengine.go` | 图像生成 + 电商多图 |
| `doubao-seedance-1-5-*` | `volcengine_video.go` | 视频生成 |
| Veo 3.1（Google Video） | `google_video.go` | 视频生成 |

Provider 使用注册表模式（`init()` 自动注册），`GetDefault()` 按优先级列表选择首个可用模型。

### 数据模型速查

```
User                  用户信息 + 钻石余额 + 邀请码 + 签到连续天数
CreditTransaction     积分流水账本（type: generate-cost/refund/register-gift/invite-reward/checkin/license-redeem/payment）
Generation            统一生成历史（type: image/video，status: generating/queued/success/failed）
InspirationPost       灵感广场 UGC 帖子（含审核状态 review_status: pending/approved/rejected）
InspirationLike       用户点赞关系表
InspirationTag        标签字典（带 slug + 使用计数）
InspirationPostTag    帖子↔标签多对多关系
InspirationReviewLog  审核操作日志（记录每次状态流转）
PaymentOrder          支付订单（provider: linuxdo）
License               兑换码（status: active/redeemed/disabled）
UserNotification      站内通知（bizKey 做幂等）
SystemSetting         系统动态 KV 配置（管理后台可实时修改）
EmailVerification     邮箱验证码（含尝试次数，防暴力枚举）
APILog                API 调用日志（含请求/响应体，定期清理）
```

### 数据库迁移系统

项目自研了一个轻量级迁移引擎（`db/db.go:runMigrations()`），无需外部依赖：

- 使用 `schema_migrations` 表追踪已执行版本
- 按文件名字典序执行 `.sql` 文件，每次启动只跑新增文件
- 幂等性错误（`duplicate column name`、`already exists` 等）静默跳过，兼容 MySQL < 8.0.4
- 真正失败时调用 `log.Fatalf()` 阻止服务启动，确保数据库结构一致

> 当前已有 **28 个迁移版本**，完整记录了从 License 系统 → 用户体系 → 生成统一模型 → 灵感广场 → 支付系统的演化路径。

### 前端架构（用户端）

**状态管理（Pinia Stores）：**

| Store | 职责 |
|-------|------|
| `userStore` | 用户信息、JWT Token、登录状态（localStorage 持久化） |
| `themeStore` | 明/暗/跟随系统主题（3 态切换，持久化） |
| `localeStore` | 中英文语言切换（持久化） |

**核心组件职责：**

| 组件 | 大小 | 职责 |
|------|------|------|
| `ComposerBar.vue` | ~49KB | 创作台核心交互：图/视/电商三模式参数配置 + 图片上传 |
| `AppSidebar.vue` | ~30KB | 左侧导航栏（含用户信息、积分显示、历史记录入口） |
| `AuthModal.vue` | ~36KB | 注册/登录/重置密码弹窗（邮箱验证码 + OAuth 双路径） |
| `ShareGenerationDialog.vue` | ~33KB | 作品发布到灵感广场对话框（含标签、描述配置） |
| `PricingModal.vue` | ~22KB | 定价/充值弹窗 |

**落地页（`Landing.vue`）设计要点：**
- Hero 区使用视频轮播（3 个视频，3s 切换，淡入淡出）
- IntersectionObserver 驱动 `.animate-on-scroll` 进入视口时触发滑入动画
- 支持独立的语言切换器（滑块式动画）和主题切换按钮

### 管理后台架构

独立的轻量 SPA（无 Pinia，无复杂路由守卫），使用 Ant Design Vue。

```
AdminLayout.vue          → 侧边导航框架
├── InspirationReview    → 灵感帖子审核（通过/拒绝，支持预览）
├── UserList             → 用户管理（搜索、钻石调整、封禁/解冻）
├── GenerationList       → 生成记录查看
└── Settings             → SystemSetting KV 动态配置
```

鉴权通过 `useAdmin` composable 统一封装，每个请求自动携带 `X-Admin-Token`。

### 积分（钻石）体系

| 事件 | 变动 |
|------|------|
| 新用户注册 | +10 |
| 每日签到 | +1（连续签到 streak 递增） |
| 邀请新用户注册 | +10（上限 500） |
| 兑换 License 码 | +码内余额 |
| 线上支付 | 按套餐方案 |
| 图像生成 1K | -1 |
| 图像生成 2K | -2 |
| 图像生成 4K | -4 |
| 视频生成（按分辨率×时长） | -N（Provider 计算） |
| 生成失败 | 自动退款 |

所有变动均写入 `credit_transactions` 流水账本，可通过 `/api/user/credits/transactions` 查询。

---

## 未来功能

- **画布支持** — 可视化画布编辑器，支持图层操作、局部重绘与自由拼接

## 致谢

感谢 [Linux.do](https://linux.do) 社区的交流与支持，项目的许多想法和改进都来自社区成员的反馈。

## 开源协议

本项目基于 [AGPL-3.0](LICENSE) 协议开源。如果你将修改后的版本部署为网络服务，必须向该服务的用户公开源代码。
