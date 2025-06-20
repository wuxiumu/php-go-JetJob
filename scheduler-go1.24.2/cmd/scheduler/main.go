package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	_ "github.com/golang-jwt/jwt/v4"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/gorilla/websocket"
	"sync"
)

var (
	DB  *gorm.DB
	RDB *redis.Client
	Ctx = context.Background()
)

// ===================== 数据模型 ==========================
type Task struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"` // "shell" 支持扩展
	Command   string    `json:"command"`
	Schedule  string    `json:"schedule"`
	Status    string    `json:"status"` // active/paused/...
	Retry     int       `json:"retry"`  // 最大重试次数
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Node struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	Name          string    `json:"name"`
	Host          string    `json:"host"`
	Status        string    `json:"status"`
	LastHeartbeat time.Time `json:"last_heartbeat"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// ws升级器
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true }, // 允许任意跨域
}

// 日志订阅者
var wsClients = make(map[*websocket.Conn]bool)
var wsMutex sync.Mutex

func wsLogsHandler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	wsMutex.Lock()
	wsClients[conn] = true
	wsMutex.Unlock()

	defer func() {
		wsMutex.Lock()
		delete(wsClients, conn)
		wsMutex.Unlock()
		conn.Close()
	}()

	for {
		_, _, err := conn.ReadMessage() // 阻塞等待客户端关闭
		if err != nil {
			break
		}
	}
}

// 日志广播函数
func broadcastLog(msg string) {
	wsMutex.Lock()
	defer wsMutex.Unlock()
	for conn := range wsClients {
		_ = conn.WriteMessage(websocket.TextMessage, []byte(msg))
	}
}

// 占位代码，避免 "os" 和 "os/exec" 导入但未使用的问题
func init() {
	_ = os.Getenv
	_ = exec.Command
	_ = http.DefaultServeMux
}

// ===================== 初始化 ===========================
func InitDB() {
	dsn := "root:dz1991wqB!@tcp(127.0.0.1:3306)/jetjob?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}
	db.AutoMigrate(&Task{}, &Node{})
	DB = db
}

func InitRedis() {
	RDB = redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	if err := RDB.Ping(Ctx).Err(); err != nil {
		panic("failed to connect redis: " + err.Error())
	}
}

// ================ Token认证中间件 =======================
func TokenAuthMiddleware(token string) gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("Authorization header:", c.GetHeader("Authorization"))
		h := strings.TrimSpace(c.GetHeader("Authorization"))
		if !strings.HasPrefix(h, "Bearer ") {
			fmt.Println("token格式错误")
			c.AbortWithStatusJSON(401, gin.H{"message": "token格式错误"})
			return
		}
		userToken := strings.TrimSpace(h[7:])
		if userToken != token {
			fmt.Println("token无效")
			c.AbortWithStatusJSON(401, gin.H{"message": "token无效"})
			return
		}
		c.Next()
	}
}

// ================ 任务CRUD接口 ==========================
func createTask(c *gin.Context) {
	var t Task
	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if t.Retry == 0 {
		t.Retry = 3 // 默认3次重试
	}
	if err := DB.Create(&t).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, t)
}
func listTasks(c *gin.Context) {
	var ts []Task
	DB.Find(&ts)
	c.JSON(200, ts)
}
func getTask(c *gin.Context) {
	var t Task
	id := c.Param("id")
	if err := DB.First(&t, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "not found"})
		return
	}
	c.JSON(200, t)
}
func updateTask(c *gin.Context) {
	var t Task
	id := c.Param("id")
	if err := DB.First(&t, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "not found"})
		return
	}
	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	DB.Save(&t)
	c.JSON(200, t)
}
func deleteTask(c *gin.Context) {
	id := c.Param("id")
	DB.Delete(&Task{}, id)
	c.Status(204)
}

// ================ 节点注册/心跳 ==========================
func createNode(c *gin.Context) {
	var n Node
	if err := c.ShouldBindJSON(&n); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	n.Status = "active"
	n.LastHeartbeat = time.Now()
	if err := DB.Create(&n).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, n)
}
func listNodes(c *gin.Context) {
	var ns []Node
	DB.Find(&ns)
	c.JSON(200, ns)
}
func nodeHeartbeat(c *gin.Context) {
	var req struct {
		NodeID uint `json:"node_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	RDB.Set(Ctx, fmt.Sprintf("node:heartbeat:%d", req.NodeID), "1", time.Minute)
	DB.Model(&Node{}).Where("id = ?", req.NodeID).
		Update("last_heartbeat", time.Now())
	c.JSON(200, gin.H{"message": "心跳OK"})
}

// ================ Worker拉取/上报 ========================
func pullTasks(c *gin.Context) {
	// 你可以做节点轮询、随机分配、权重分配等，这里简单直接给全部active任务
	var ts []Task
	DB.Where("status = ?", "active").Find(&ts)
	c.JSON(200, ts)
}
func reportTask(c *gin.Context) {
	var report struct {
		TaskID uint   `json:"task_id"`
		NodeID uint   `json:"node_id"`
		Status string `json:"status"`
		Log    string `json:"log"`
	}
	if err := c.ShouldBindJSON(&report); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	// 更新任务状态和日志（这里简单演示，实际应有tasks_log表等）
	if report.Status == "success" {
		DB.Model(&Task{}).Where("id = ?", report.TaskID).Update("status", "done")
	} else {
		// 如果失败，计数+1，超过重试最大次数就终止
		key := fmt.Sprintf("task:%d:failcount", report.TaskID)
		cnt, _ := RDB.Incr(Ctx, key).Result()
		var t Task
		DB.First(&t, report.TaskID)
		if cnt >= int64(t.Retry) {
			DB.Model(&Task{}).Where("id = ?", report.TaskID).Update("status", "failed")
			// 可在此触发告警
		} else {
			DB.Model(&Task{}).Where("id = ?", report.TaskID).Update("status", "active")
		}
	}
	broadcastLog(fmt.Sprintf("任务%d 上报状态: %s", report.TaskID, report.Status))
	c.JSON(200, gin.H{"message": "上报成功"})
}

// ============== 节点自动下线监控 ==========================
func monitorNodes() {
	for {
		var nodes []Node
		DB.Find(&nodes)
		now := time.Now()
		for _, n := range nodes {
			if now.Sub(n.LastHeartbeat) > time.Minute {
				if n.Status != "offline" {
					DB.Model(&Node{}).Where("id = ?", n.ID).
						Update("status", "offline")
				}
			} else {
				if n.Status != "active" {
					DB.Model(&Node{}).Where("id = ?", n.ID).
						Update("status", "active")
				}
			}
		}
		time.Sleep(30 * time.Second)
	}
}

// =============== 调度主循环（分发、并发、重试） ============
func schedulerLoop() {
	ticker := time.NewTicker(30 * time.Second)
	maxConcurrent := 3 // 可配置
	for {
		<-ticker.C
		var ts []Task
		DB.Where("status = ?", "active").Find(&ts)
		var ns []Node
		DB.Where("status = ?", "active").Find(&ns)
		for _, t := range ts {
			// 简单限流：看每个节点任务数
			for _, n := range ns {
				cKey := fmt.Sprintf("node:%d:running", n.ID)
				num, _ := RDB.Get(Ctx, cKey).Int()
				if num >= maxConcurrent {
					continue // 超并发不分配
				}
				// 分配任务给节点（你可以推送到队列/消息中间件/写Redis等，这里只是示例）
				_ = RDB.Incr(Ctx, cKey).Err()
				fmt.Printf("调度：任务%d 分配给节点%d\n", t.ID, n.ID)
				// 实际项目建议在此写任务分发/消息推送/队列/通知Worker拉取
				break
			}
		}
	}
}

// =============== 用户登录接口 ====================
func loginHandler(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "参数错误"})
		return
	}
	// 这里你可以查数据库、配置等，这里演示账号 admin 密码 123456
	if req.Username == "admin" && req.Password == "123456" {
		// 返回一个“伪token”，如"test-token"
		c.JSON(200, gin.H{"token": "test-token"})
	} else {
		c.JSON(401, gin.H{"error": "账号或密码错误"})
	}
}

// =============== main函数 ================================
func main() {
	InitDB()
	InitRedis()
	token := "test-token"
	go monitorNodes()
	go schedulerLoop()

	r := gin.Default()
	r.GET("/ws/logs", wsLogsHandler)
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})
	r.POST("/api/login", loginHandler)

	r.Use(TokenAuthMiddleware(token))
	r.POST("/api/tasks", createTask)
	r.GET("/api/tasks", listTasks)
	r.GET("/api/tasks/:id", getTask)
	r.PUT("/api/tasks/:id", updateTask)
	r.DELETE("/api/tasks/:id", deleteTask)
	r.POST("/api/nodes", createNode)
	r.GET("/api/nodes", listNodes)
	r.POST("/api/nodes/heartbeat", nodeHeartbeat)
	r.GET("/api/tasks/pull", pullTasks)
	r.POST("/api/tasks/report", reportTask)
	r.GET("/ping", func(c *gin.Context) { c.JSON(200, gin.H{"pong": true}) })

	r.Run(":8090")
}
