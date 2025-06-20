package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
)

// 日志WebSocket Hub
type LogHub struct {
	clients map[*websocket.Conn]bool
	lock    sync.Mutex
}

var logHub = &LogHub{
	clients: make(map[*websocket.Conn]bool),
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// 客户端连接
func LogWebSocketHandler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	logHub.lock.Lock()
	logHub.clients[conn] = true
	logHub.lock.Unlock()

	// 简单读协程，遇到错误断开连接
	go func() {
		defer func() {
			logHub.lock.Lock()
			delete(logHub.clients, conn)
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

// 广播日志到所有客户端
func BroadcastLog(msg string) {
	logHub.lock.Lock()
	defer logHub.lock.Unlock()
	for conn := range logHub.clients {
		conn.WriteMessage(websocket.TextMessage, []byte(msg))
	}
}
