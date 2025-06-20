package api

import (
	_ "errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"jetjob/internal/model"
	"jetjob/internal/storage"
	_ "time"
)

// ================ Worker拉取/上报 ========================
func pullTasks(c *gin.Context) {
	// 你可以做节点轮询、随机分配、权重分配等，这里简单直接给全部active任务
	var ts []model.Task
	storage.DB.Where("status = ?", "active").Find(&ts)
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
		storage.DB.Model(&model.Task{}).Where("id = ?", report.TaskID).Update("status", "done")
	} else {
		// 如果失败，计数+1，超过重试最大次数就终止
		key := fmt.Sprintf("task:%d:failcount", report.TaskID)
		cnt, _ := storage.RDB.Incr(storage.Ctx, key).Result()
		var t model.Task
		storage.DB.First(&t, report.TaskID)
		if cnt >= int64(t.Retry) {
			storage.DB.Model(&model.Task{}).Where("id = ?", report.TaskID).Update("status", "failed")
			// 可在此触发告警
		} else {
			storage.DB.Model(&model.Task{}).Where("id = ?", report.TaskID).Update("status", "active")
		}
	}
	broadcastLog(fmt.Sprintf("任务%d 上报状态: %s", report.TaskID, report.Status))
	c.JSON(200, gin.H{"message": "上报成功"})
}
