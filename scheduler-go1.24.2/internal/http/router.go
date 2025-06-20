// internal/http/router.go
package http

import (
    "github.com/gin-gonic/gin"
    "net/http"
)

func RegisterRoutes(r *gin.Engine) {
    r.POST("/api/nodes/register", RegisterNode)
    r.POST("/api/nodes/heartbeat", NodeHeartbeat)
}

func RegisterNode(c *gin.Context) {
    // 解析参数，注册到内存/DB
    // ...略
    c.JSON(http.StatusOK, gin.H{"message": "注册成功"})
}

func NodeHeartbeat(c *gin.Context) {
    // 解析参数，更新节点心跳时间
    // ...略
    c.JSON(http.StatusOK, gin.H{"message": "心跳OK"})
}