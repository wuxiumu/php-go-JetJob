# **JetJob / php-go-scheduler**



> 轻量级分布式任务调度平台 | Go + MySQL + Redis + Vue3

> 支持多类型任务（Shell/HTTP/文件）、定时/依赖/手动触发、节点心跳、WebSocket 日志推送



------





## **✨ 项目简介**





JetJob 是一套易于二次开发、面向工程实践的分布式任务调度平台。支持多种任务类型、节点弹性扩展，任务定时/依赖/并发控制，前后端完全分离，适合企业级调度、自动化运维、开发测试等场景。



------





## **🏗️ 技术栈**





- **后端：**



- Go 1.18+ / 1.20+
- Gin (Web 框架)
- GORM (MySQL ORM)
- go-redis (Redis 客户端)
- gorilla/websocket (日志推送)



- **数据库：** MySQL 5.7/8.0+

- **缓存：** Redis 5.0+

- **前端：** Vue3 + Vite + Element Plus

- **其他：** 支持 Docker 部署，支持 JWT/Token 认证





------





## **🧩 核心功能**





- **任务调度：**



- Shell、HTTP、文件型任务
- 支持 Cron 定时、依赖、手动三种触发模式
- 任务重试、并发数限制
- 任务参数自定义、任务输出日志保存



- **节点管理：**



- 节点注册/心跳、离线自动检测
- 节点状态、负载信息可视化



- **任务管理：**



- 任务新建/编辑/删除，任务历史与状态
- 支持批量操作



- **Worker/Agent：**



- 支持任务拉取、执行和状态上报
- 支持横向扩容



- **权限认证：**



- 支持 Token 鉴权
- 管理员、普通用户分级管理（可选 RBAC 扩展）



- **日志中心：**



- 任务日志实时 WebSocket 推送
- 操作日志和调度审计



- **前后端分离：**



- 环境变量灵活切换接口地址
- 本地与线上快速切换







------





## **📦 目录结构**



```
php-go-scheduler/
├── cmd/jetjob/                # 启动入口 main.go
├── config/                    # 配置文件
├── internal/
│   ├── api/                   # 路由和 handler
│   ├── app/                   # 业务逻辑/服务层
│   ├── model/                 # 数据结构
│   ├── storage/               # DB/Redis 连接
│   └── utils/                 # 配置、token校验等工具
├── frontend/                  # Vue3 前端项目
├── go.mod / go.sum
├── README.md
└── ...
```



------





## **🚀 快速启动**







### **1. 克隆项目 & 初始化依赖**



```
git clone https://github.com/yourname/php-go-scheduler.git
cd php-go-scheduler
go mod tidy
```



### **2. 配置数据库和 Redis**





- 复制一份默认配置：



```
cp config/config.yaml.example config/config.yaml
```



- 修改 config.yaml，填写 MySQL 和 Redis 连接信息，设置初始 token。







### **3. 启动后端服务**



```
go run ./cmd/jetjob/main.go
```



- 默认监听 :8090，日志输出到控制台。







### **4. 启动前端管理台**



```
cd frontend
npm install
npm run dev
```



- 默认端口 5173，配置接口地址见 .env.development。





------





## **🔑 API & 权限**







### **Token 鉴权**





- 登陆后端管理台，需输入初始 Token
- 所有 API 自动携带 Authorization: Bearer {token}







### **典型接口文档**



| **接口**             | **方法** | **说明**        |
| -------------------- | -------- |---------------|
| /api/register           | POST     | 用户注册          |
| /api/login           | POST     | 用户登录，换取 token |
| /api/tasks           | GET      | 任务列表          |
| /api/tasks           | POST     | 创建任务          |
| /api/tasks/:id       | GET      | 任务详情          |
| /api/tasks/:id       | PUT      | 编辑任务          |
| /api/tasks/:id       | DELETE   | 删除任务          |
| /api/nodes/register  | POST     | 节点注册          |
| /api/nodes/heartbeat | POST     | 节点心跳          |
| /ws/logs             | GET      | WebSocket 日志  |
| …                    |          | 更多见代码注释       |



- API 支持前端环境变量配置，便于本地和生产环境切换。





------





## **📝 任务类型说明**





- **Shell 任务：**

  在 Worker 机器上执行指定 shell 命令或脚本，可带参数和输出采集。

- **HTTP 任务：**

  以 HTTP 客户端方式请求指定接口，支持自定义 header/body、支持 webhook 场景。

- **文件任务：**

  支持上传、下载、脚本/二进制分发等。





**每种任务都支持参数配置、输出日志存储、失败重试、任务依赖链路。**



------





## **🖥️ 节点与 Worker 说明**





- Worker 节点通过 API 注册，定时心跳上报。
- 节点失联自动标记为 offline。
- 支持多节点横向扩容、主备热切换。





------





## **📊 日志推送与可观测性**





- 所有任务执行日志支持 WebSocket 实时推送，管理台一键查看。
- 支持历史日志检索、关键错误告警（可扩展邮件/短信/钉钉）。





------





## **⚡ 前端一键对接（环境变量）**





- 前端 frontend/.env.development



```
VITE_API_BASE=http://localhost:8090
```



-
- axios 自动带 token，全局拦截 401/错误弹窗
- 切换生产只需 .env.production 配置接口地址





------





## **🛠️ 开发/运维建议**





- **开发环境推荐 Docker 一键启动 MySQL/Redis，避免本地依赖冲突**
- **生产部署建议 Supervisor 或 systemd 保活，日志写入独立目录**
- **集成 CI：每次提交自动 go test ./… 保证主干代码稳定**





------





## **📈 扩展与自定义**





- 可对接企业现有权限系统、CMDB、告警平台
- 支持二次开发新任务类型（如 Python、Go 脚本、K8s CronJob 等）
- 支持更多通知方式（钉钉/企业微信/邮件）





------





## **🧩 常见问题 & FAQ**





- **Q: 前端接口 401 无法登陆？**

  A: 请确认 config.yaml 中 token 设置正确，前端本地输入一致。

- **Q: 任务调度无响应？**

  A: 请检查 Worker 节点是否在线，任务是否为激活状态。

- **Q: 日志推送收不到？**

  A: 请确认浏览器 ws 连接 ws://localhost:8090/ws/logs 已建立，后端有日志输出。

- **Q: 跨域请求失败？**

  A: 后端 Gin CORS 中间件是否已配置（见 main.go 示例）。





------





## **📸 示例截图与架构图**





> 建议添加平台首页、任务管理、节点监控等前端截图，以及 draw.io 架构图。



------





## **🏆 Star/贡献**





欢迎 PR、issue、star！

如有定制化需求、企业集成咨询，请联系项目作者。



------





## **License**





[MIT](./LICENSE)



------



> 本项目长期维护，适用于团队二次开发、企业生产环境试用。



------



如需**补充“接口字段说明表格”、更详细的流程图、部署脚本、FAQ、演示视频文案**等，随时告知！