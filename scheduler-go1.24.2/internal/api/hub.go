package api

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
	"time"
)

// 消息类型常量
const (
	EventLog        = "log"
	EventNodeStatus = "node_status"
	EventAlarm      = "alarm"
	EventTaskAssign = "task_assign"
	EventBalance    = "balance_notice"
)

// WebSocket客户端结构
type Client struct {
	Conn     *websocket.Conn
	NodeName string // 管理端为""，节点端为节点名
}

// Hub管理全部连接
type Hub struct {
	clients map[*Client]bool
	lock    sync.Mutex
}

var logHub = &Hub{
	clients: make(map[*Client]bool),
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// WebSocket连接处理
func LogWebSocketHandler(c *gin.Context) {
	nodeName := c.Query("node")
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	client := &Client{
		Conn:     conn,
		NodeName: nodeName,
	}
	logHub.lock.Lock()
	logHub.clients[client] = true
	logHub.lock.Unlock()

	// 断开连接回收
	go func() {
		defer func() {
			logHub.lock.Lock()
			delete(logHub.clients, client)
			logHub.lock.Unlock()
			conn.Close()
		}()
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				break
			}
		}
	}()
}

// 广播多事件给所有客户端
func BroadcastEvent(event string, data interface{}) {
	msg, _ := json.Marshal(gin.H{
		"event": event,
		"data":  data,
	})
	logHub.lock.Lock()
	defer logHub.lock.Unlock()
	for client := range logHub.clients {
		_ = client.Conn.WriteMessage(websocket.TextMessage, msg)
	}
}

// 单节点推送（只推目标节点）
func SendEventToNode(nodeName, event string, data interface{}) {
	msg, _ := json.Marshal(gin.H{
		"event": event,
		"data":  data,
	})
	logHub.lock.Lock()
	defer logHub.lock.Unlock()
	for client := range logHub.clients {
		if client.NodeName == nodeName {
			_ = client.Conn.WriteMessage(websocket.TextMessage, msg)
		}
	}
}

// 分组推送（如所有worker/所有管理端）
func SendEventToGroup(nodeNames []string, event string, data interface{}) {
	nameSet := make(map[string]bool)
	for _, n := range nodeNames {
		nameSet[n] = true
	}
	msg, _ := json.Marshal(gin.H{
		"event": event,
		"data":  data,
	})
	logHub.lock.Lock()
	defer logHub.lock.Unlock()
	for client := range logHub.clients {
		if nameSet[client.NodeName] {
			_ = client.Conn.WriteMessage(websocket.TextMessage, msg)
		}
	}
}

// 负载均衡通知：推送当前推荐调度节点
func BalanceNotice(bestNode string, currentLoads map[string]float64) {
	notice := gin.H{
		"recommend": bestNode,
		"loads":     currentLoads,
		"time":      time.Now().Format(time.RFC3339),
	}
	BroadcastEvent(EventBalance, notice)
}

// 示例：调用负载均衡通知
func ExampleBalancePush() {
	loads := map[string]float64{
		"node-01": 0.12,
		"node-02": 1.56,
		"node-03": 0.41,
	}
	best := pickBestNode(loads) // 你的调度算法
	BalanceNotice(best, loads)
}

// 你的调度算法举例（简单选择负载最低节点）
func pickBestNode(loads map[string]float64) string {
	best := ""
	min := 1e9
	for n, l := range loads {
		if l < min {
			min = l
			best = n
		}
	}
	return best
}

// 路由注册
func RegisterHubRoutes(r *gin.Engine) {
	r.GET("/ws/logs", LogWebSocketHandler)
}
