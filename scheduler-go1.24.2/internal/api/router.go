package api

import "github.com/gin-gonic/gin"

// RegisterRoutes 注册所有接口路由
func RegisterRoutes(r *gin.Engine) {

	// 注册与登录（无需token）

	r.POST("/api/nodes/register", RegisterNode)
	r.POST("/api/nodes/heartbeat", NodeHeartbeat)
	// r.GET("/api/nodes", ListNodes) // 可扩展
	// r.POST("/api/tasks", CreateTask) // 可扩展
}
