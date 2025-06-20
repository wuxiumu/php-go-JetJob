package main

import (
	"github.com/gin-gonic/gin"
	"jetjob/internal/api"
	"jetjob/internal/model"
	"jetjob/internal/service"
	"jetjob/internal/storage"
	"jetjob/internal/utils"
	"time"
)

func StartNodeHealthMonitor(timeout time.Duration) {
	go func() {
		for {
			time.Sleep(time.Minute)
			var nodes []model.Node
			storage.DB.Find(&nodes)
			now := time.Now()
			for _, node := range nodes {
				if node.Status == "active" && now.Sub(node.LastHeartbeat) > timeout {
					storage.DB.Model(&model.Node{}).
						Where("id = ?", node.ID).
						Update("status", "offline")
					// 可推送告警事件
					api.BroadcastEvent("node_status", gin.H{
						"name":   node.Name,
						"group":  node.Group,
						"status": "offline",
					})
				}
			}
		}
	}()
}

// =============== main函数 ================================
func main() {
	// 1. 配置、依赖初始化
	utils.LoadConfig("config/config.yaml")
	storage.InitDB(utils.Cfg.MySQL.DSN)
	storage.InitRedis(utils.Cfg.Redis.Addr)
	//log.InitLogger(...)

	// DB表自动迁移
	storage.DB.AutoMigrate(&model.Task{}, &model.Node{}, &model.User{})

	// 3. 健康检测 goroutine     service.StartNodeHealthMonitor(...)
	service.StartNodeHealthMonitor(2 * time.Minute)

	// 4. Gin 实例与中间件
	r := gin.Default()
	r.Use(utils.CORS())
	r.Use(utils.JWTAuthMiddleware())

	// 5. 统一注册所有HTTP和WS路由
	api.RegisterRoutes(r)
	api.RegisterHubRoutes(r) // WebSocket

	// 健康检测

	//r.GET("/ws/logs", api.wsLogsHandler)
	// CORS
	//r.Use(func(c *gin.Context) {
	//	c.Header("Access-Control-Allow-Origin", "*")
	//	c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
	//	c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type")
	//	if c.Request.Method == "OPTIONS" {
	//		c.AbortWithStatus(204)
	//		return
	//	}
	//	c.Next()
	//})
	//
	//// Token认证
	//r.Use(utils.TokenAuthMiddleware(utils.Cfg.Token))
	//
	//r.POST("/api/tasks", api.CreateTaskHandler)
	//r.GET("/api/tasks", api.ListTasksHandler)
	//
	//// ...初始化代码
	////api.RegisterRoutes(r)
	//r.POST("/api/register", api.RegisterHandler)
	//r.POST("/api/login", api.LoginHandler)
	//
	//r.Use(api.JWTAuthMiddleware())
	//r.GET("/api/user", api.GetUserHandler)
	//r.PUT("/api/user", api.UpdateUserHandler)
	//r.DELETE("/api/user", api.DeleteUserHandler)
	//
	//// 需在 Use(JWTAuthMiddleware()) 之后注册
	//r.GET("/api/tasks", api.ListTasksHandler)
	//r.POST("/api/tasks", api.CreateTaskHandler)
	//r.GET("/api/tasks/:id", api.GetTaskHandler)
	//r.PUT("/api/tasks/:id", api.UpdateTaskHandler)
	//r.DELETE("/api/tasks/:id", api.DeleteTaskHandler)
	//
	//r.POST("/api/nodes/register", api.RegisterNodeHandler)
	//r.POST("/api/nodes/heartbeat", api.NodeHeartbeatHandler)
	//
	//r.GET("/ws/logs", api.LogWebSocketHandler)
	//// 注册 WebSocket 日志/事件推送路由
	//api.RegisterHubRoutes(r)
	//
	//r.GET("/api/nodes", api.ListNodesHandler)
	r.GET("/ping", func(c *gin.Context) { c.JSON(200, gin.H{"pong": true}) })

	r.Run(":8090")
}
