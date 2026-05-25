# Goravel Admin 本地运行指南

## 环境要求

- **Go**: 1.24.0+
- **Node.js**: 20.x (使用 nvm 管理)
- **Docker**: 用于 MySQL 和 Redis

## 已完成的配置

- ✅ Docker 容器已创建 (MySQL:3306, Redis:6379)
- ✅ 后端 .env 文件已配置
- ✅ 前端 .env 文件已配置

## 你需要完成的步骤

### 1. 安装 Go 1.24

在终端中运行:

```bash
# 方法1: 使用 Homebrew
brew install go@1.24

# 方法2: 手动下载
# 访问 https://go.dev/dl/ 下载 macOS arm64 版本
# 然后安装
```

安装后验证:
```bash
go version  # 应该显示 go1.24.0 或更高
```

### 2. 配置 Go 模块代理 (国内网络)

```bash
go env -w GOPROXY=https://goproxy.cn,https://mirrors.aliyun.com/goproxy,direct
```

### 3. 启动后端服务

```bash
# 进入项目目录
cd /Users/liu/code/wy-code/goravel-admin

# 生成应用 Key 和 JWT Secret
go run . artisan key:generate

# 运行数据库迁移
go run . artisan migrate

# 填充初始数据
go run . artisan db:seed

# 启动后端服务 (开发模式)
go run . --no-ansi

# 或者使用 air 热重载 (推荐)
air
```

后端服务将在 http://localhost:3000 启动

### 4. 启动前端服务

```bash
# 进入前端目录
cd /Users/liu/code/wy-code/goravel-admin/html

# 使用 nvm 切换到 Node 20
nvm use 20

# 安装依赖
npm install

# 启动开发服务器
npm run dev
```

前端服务将在 http://localhost:5173 (或类似端口) 启动

### 5. 访问后台管理系统

- **前端地址**: http://localhost:5173
- **后端API**: http://localhost:3000
- **Swagger文档**: http://localhost:3000/swagger/index.html

### 6. 登录信息

- **用户名**: admin
- **密码**: admin123

## 常见问题

### Docker 容器未启动
```bash
docker-compose -f docker-compose.dev.yml up -d
```

### 端口被占用
修改 `.env` 中的 `APP_PORT` 为其他端口

### 前端依赖安装失败
```bash
cd html
rm -rf node_modules package-lock.json
npm install
```

## 项目结构

```
goravel-admin/
├── app/                 # 后端应用代码
├── html/                # 前端 Vue 项目
├── docker-compose.dev.yml  # 开发环境 Docker 配置
├── .env                 # 后端配置文件
└── html/.env            # 前端配置文件
```
