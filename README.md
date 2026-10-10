# JetJob - 轻量级分布式任务调度平台

> 适合中小型团队和个人开发者的定时任务分发系统，可私有化部署。

## 功能特色

- 支持定时/周期性任务，crontab表达式
- 多节点分布式，任务高可用自动转移
- Web管理后台，任务监控&历史追踪
- 私有化部署，可按环境扩展

## 技术架构

- Go（`scheduler-go1.24.2`）：调度进程。当前可编译入口是 `cmd/scheduler`，默认 `:8090`。`cmd/jetjob` 现在不能编译：`internal/api` 有重复声明和未定义符号，`main` 还调用了不存在的 `utils.CORS` 与 `utils.JWTAuthMiddleware`。
- PHP（`backend-php7.4`，Laravel 8）：管理 API，前缀 `/api/v1`，含登录、任务、节点、角色权限。`config/jwt.php` 读取 `JWT_SECRET`。
- 前端：Vue 3.5 + Vite + Element Plus（`frontend-vue3.5.13`）。开发服务器默认 `:5173`。登录页请求 `${VITE_API_BASE}/login`，对接 Go 调度时应设为 `http://127.0.0.1:8090/api`。
- 数据库：MySQL、Redis

`docker-compose.yml` 只启动 MySQL 8 与 Redis 7，不包含 PHP、Go、Vue，因此不是整套一键部署。

产品功能说明见 `scheduler-go1.24.2/README.md`。启动命令以本文为准。

## 配置

```bash
cp .env.example .env
cp backend-php7.4/.env.example backend-php7.4/.env
cp frontend-vue3.5.13/.env.example frontend-vue3.5.13/.env.development
cp scheduler-go1.24.2/.env.example scheduler-go1.24.2/.env
cp scheduler-go1.24.2/config/config.yaml.example scheduler-go1.24.2/config/config.yaml
```

口令都是占位符，填好的 `.env` 不要提交。`cmd/scheduler` 读取 `scheduler-go1.24.2/.env` 中的 `JETJOB_MYSQL_DSN`、`JETJOB_REDIS_ADDR`、`JETJOB_TOKEN`；未设置时使用与示例相同的占位默认值。`cmd/jetjob` 读取 `config/config.yaml`，但该入口目前无法编译。

演示登录用户名 `admin`、密码 `123456`，只适合本机。PHP 复制 `.env` 后执行 `php artisan key:generate`，需要库表时再 `php artisan migrate`。

## 本地启动

依赖：Go 1.24.2、PHP 7.4+ 与 Composer、Node.js。MySQL / Redis 可本机安装，或只拉起依赖：

```bash
docker compose up -d
```

应用分开启动：

```bash
cd scheduler-go1.24.2 && go run ./cmd/scheduler
cd backend-php7.4 && composer install && php artisan serve
cd frontend-vue3.5.13 && npm install && npm run dev
```

`cmd/scheduler` 在连上 MySQL 和 Redis 之前会直接退出。Vue 登录页使用 Go 返回的 `token`。PHP 登录在 `/api/v1/login`，返回 `access_token`，当前登录页不读取该字段。

## 冒烟

不启动数据库，只跑不依赖 MySQL 的 Go 测试：

```bash
./scripts/smoke.sh
```

## License

[MIT](./LICENSE)
