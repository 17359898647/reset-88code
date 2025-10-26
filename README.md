# 88code 智能额度重置

使用 Go 语言实现的 88code 额度自动重置工具，支持智能分时重置策略。

## ✨ 功能特性

- 🕐 **智能分时策略**：根据北京时间自动选择重置策略
- 🚀 **零依赖部署**：单一二进制文件，无需额外依赖
- 🔄 **自动化执行**：GitHub Actions 定时任务自动运行
- 📊 **详细日志**：完整的请求响应和统计信息
- 🛡️ **容错机制**：单个订阅失败不影响其他订阅

## ⏰ 重置策略

### 18:30 时段（选择性重置）
- **条件**：仅重置 `resetTimes=2` 且 `当前额度 ≤ 20% 总额度` 的订阅
- **执行方式**：并发执行，速度更快
- **适用场景**：傍晚时段快速恢复低额度账户

### 23:45 时段（全量重置）
- **条件**：无条件重置所有活跃订阅
- **执行方式**：顺序执行，带请求间隔
- **适用场景**：每日完整重置，确保所有账户额度恢复

### 其他时段（手动执行）
- **条件**：无条件重置所有活跃订阅
- **执行方式**：顺序执行
- **适用场景**：手动触发时的全量重置

## 🔧 环境变量配置

### 必需环境变量

| 环境变量 | 说明 | 示例 |
|---------|------|------|
| `TOKEN` | 88code API 认证令牌 | `30846ca74de34841b6e21020ea28398b` |

### 获取 TOKEN

1. 登录 88code 账户
2. 打开浏览器开发者工具（F12）
3. 切换到 Network 标签页
4. 访问 My Subscription 页面
5. 查找请求头中的 `Authorization: Bearer <token>`
6. 复制 Bearer 后面的 token 值

## 🚀 使用方法

### GitHub Actions（推荐）

1. **Fork 本仓库**

2. **配置 Secret**
   - 进入仓库 `Settings` → `Secrets and variables` → `Actions`
   - 点击 `New repository secret`
   - 添加以下 Secret：
     - **Name**: `TOKEN`
     - **Value**: `你的 88code token`

3. **启用 Actions**
   - 确保仓库的 Actions 功能已启用
   - 程序将自动在每天 18:30 和 23:45（北京时间）运行

4. **手动触发**
   - 进入 `Actions` 页面
   - 选择 "88code 每日重置额度 (Go版本)" 工作流
   - 点击 `Run workflow` 按钮

### 本地运行

#### 前置要求
- Go 1.21 或更高版本

#### 快速开始

```bash
# 1. 克隆仓库
git clone <repository-url>
cd reset-88code

# 2. 设置环境变量并运行
TOKEN=your_token make run

# 或者创建 .env 文件
cp .env.example .env
# 编辑 .env 填写 TOKEN
export $(cat .env | xargs) && make run
```

#### 编译

```bash
# 编译 Linux 版本（用于服务器）
make build

# 手动编译其他平台
GOOS=darwin GOARCH=arm64 go build -o reset-88code main.go  # macOS
GOOS=windows GOARCH=amd64 go build -o reset-88code.exe main.go  # Windows
```

## 📂 项目结构

```
reset-88code/
├── .github/
│   └── workflows/
│       └── reset-credits-go.yml  # GitHub Actions 工作流
├── pkg/
│   └── reset/
│       ├── client.go            # HTTP 客户端
│       ├── config.go            # 配置管理
│       ├── reset.go             # 重置逻辑
│       ├── subscription.go      # 订阅查询
│       └── types.go             # 类型定义
├── bin/
│   └── main                     # 编译后的二进制（Linux）
├── main.go                      # 主程序入口
├── go.mod                       # Go 模块定义
├── Makefile                     # 构建脚本
├── .env.example                 # 环境变量示例
└── README.md                    # 项目文档
```

## 📊 执行日志示例

### 18:30 时段（选择性重置）

```
开始重置额度...
Token: 30846ca7****
发送请求到: https://www.88code.org/admin-api/cc-admin/system/subscription/my
响应状态码: 200
✅ 获取订阅信息成功，共 8 个订阅
活跃订阅数: 4 / 8
  [活跃] ID: 31645, 重置次数: 2, 当前额度: 13.32, 总额度: 100
  [活跃] ID: 31429, 重置次数: 2, 当前额度: 100.00, 总额度: 100
  [活跃] ID: 27239, 重置次数: 2, 当前额度: 20.00, 总额度: 20
  [活跃] ID: 178, 重置次数: 2, 当前额度: 18.85, 总额度: 100
当前北京时间: 2025-10-27 18:30:15
🕐 当前时间段: 18:00-18:59，执行 18:30 重置策略
执行 18:30 重置策略，共 4 个订阅
分类完成 - 重置次数 0: 0 个, 1: 0 个, 2: 4 个
  订阅 ID 31645 需要重置 (当前额度: 13.32, 总额度: 100, 使用率: 13.3%)
  订阅 ID 178 需要重置 (当前额度: 18.85, 总额度: 100, 使用率: 18.9%)
开始并发重置 2 个订阅
正在重置订阅 ID 31645...
正在重置订阅 ID 178...
响应状态码: 200
响应内容: {"code":0,"msg":"操作成功","ok":true}
✅ 订阅 ID 31645 重置成功
响应状态码: 200
响应内容: {"code":0,"msg":"操作成功","ok":true}
✅ 订阅 ID 178 重置成功

📊 重置完成统计：
   总订阅数: 4
   成功数量: 2
   失败数量: 0
🎉 所有订阅额度重置成功！
```

### 23:45 时段（全量重置）

```
当前北京时间: 2025-10-27 23:45:30
🕐 当前时间段: 23:00-00:59，执行 23:45 重置策略
执行 23:45 重置策略，共 4 个订阅
无条件重置所有活跃订阅
正在重置第 1/4 个订阅 (ID: 31645, 重置次数: 2, 当前额度: 66.61, 总额度: 100)...
响应状态码: 200
响应内容: {"code":0,"msg":"操作成功","ok":true}
✅ 订阅 ID 31645 重置成功
...

📊 重置完成统计：
   总订阅数: 4
   成功数量: 4
   失败数量: 0
🎉 所有订阅额度重置成功！
```

## 🔒 安全说明

- TOKEN 是敏感信息，请勿泄露或提交到代码仓库
- 使用 GitHub Secrets 安全存储 TOKEN
- `.gitignore` 已配置忽略 `.env` 文件
- 日志中 TOKEN 会被部分掩码显示（仅显示前8位）

## 🛠️ 技术栈

- **语言**: Go 1.21+
- **标准库**: net/http, encoding/json, sync
- **CI/CD**: GitHub Actions
- **时区**: Asia/Shanghai (北京时间)

## 📝 开发

### 本地测试

```bash
# 运行测试
TOKEN=your_token make run

# 编译本地版本
go build -o reset-88code main.go

# 运行编译后的程序
TOKEN=your_token ./reset-88code
```

### 代码结构

- `pkg/reset/client.go` - HTTP 客户端封装，自动注入 Bearer Token
- `pkg/reset/config.go` - 环境变量加载和验证
- `pkg/reset/subscription.go` - 订阅信息获取和过滤
- `pkg/reset/reset.go` - 重置策略实现
- `pkg/reset/types.go` - 数据结构定义

## ⚠️ 注意事项

- 确保 TOKEN 有效，过期的 TOKEN 会导致 401 错误
- GitHub Actions 执行时间可能有 1-5 分钟延迟
- 并发重置仅在 18:30 时段使用，避免请求过于频繁
- 手动执行会触发全量重置，请谨慎操作

## 📄 License

MIT

---

🤖 Generated with [Claude Code](https://claude.com/claude-code)
