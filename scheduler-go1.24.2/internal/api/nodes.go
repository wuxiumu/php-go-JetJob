package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"jetjob/internal/model"
	"jetjob/internal/storage"
	"net/http"
	"sync"
	"time"
)

// Node 模型（可存数据库，也支持后续扩展内存/redis等）
type Node struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	Name          string    `gorm:"uniqueIndex;size:64;not null" json:"name"` // 节点名称唯一
	Host          string    `gorm:"size:255" json:"host"`                     // 主机信息或IP
	LastHeartbeat time.Time `json:"last_heartbeat"`
	Status        string    `gorm:"size:16;default:'active'" json:"status"` // active, offline
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

var nodeLock sync.Mutex // 内存锁，防止并发冲突

// 注册节点请求体
type NodeRegisterReq struct {
	Name  string `json:"name" binding:"required"`
	Group string `json:"group"` // 可选分组名
	Host  string `json:"host" binding:"required"`
}

// 心跳请求体
type NodeHeartbeatReq struct {
	Name    string  `json:"name" binding:"required"`
	LoadAvg float64 `json:"load_avg"`
}

// 节点注册
func RegisterNodeHandler(c *gin.Context) {
	var req NodeRegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	if req.Group == "" {
		req.Group = "default"
	}
	nodeLock.Lock()
	defer nodeLock.Unlock()
	var count int64
	storage.DB.Model(&Node{}).Where("name = ?", req.Name).Count(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "节点已存在"})
		return
	}
	node := Node{
		Name:          req.Name,
		Host:          req.Host,
		LastHeartbeat: time.Now(),
		Status:        "active",
	}
	if err := storage.DB.Create(&node).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "注册失败"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "注册成功"})
}

// 节点心跳
func NodeHeartbeatHandler(c *gin.Context) {
	var req NodeHeartbeatReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	nodeLock.Lock()
	defer nodeLock.Unlock()
	var node Node
	if err := storage.DB.Where("name = ?", req.Name).First(&node).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "节点不存在"})
		return
	}
	node.LastHeartbeat = time.Now()
	node.Status = "active"
	node.LoadAvg = req.LoadAvg
	if err := storage.DB.Save(&node).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "心跳失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "心跳成功"})

	// 主动推送负载变化
	BroadcastLog("[节点负载] " + node.Name + " 当前负载: " + fmt.Sprintf("%.2f", req.LoadAvg))
}

func ListNodesHandler(c *gin.Context) {
	group := c.DefaultQuery("group", "")
	var nodes []model.Node
	query := storage.DB
	if group != "" {
		query = query.Where("group = ?", group)
	}
	query.Order("id desc").Find(&nodes)
	c.JSON(http.StatusOK, nodes)
}

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

func StartNodeHealthMonitor(timeout time.Duration) {
	go func() {
		for {
			time.Sleep(time.Minute)
			offlined, _ := nodeSvc.OfflineTimeoutNodes(timeout)
			for _, node := range offlined {
				// 离线推送事件
				api.BroadcastEvent("node_status", gin.H{
					"name":   node.Name,
					"group":  node.Group,
					"status": "offline",
				})
			}
		}
	}()
}
