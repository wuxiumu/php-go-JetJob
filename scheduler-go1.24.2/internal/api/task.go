package api

import (
	"github.com/gin-gonic/gin"
	"jetjob/internal/storage"
	"net/http"
	"strconv"
	"time"
)

// Task 模型（最基本的字段，可自行扩展）
type Task struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	Title     string     `gorm:"size:128;not null" json:"title"`
	Content   string     `gorm:"type:text" json:"content"`
	Status    string     `gorm:"size:32;default:'pending'" json:"status"` // pending, running, success, failed
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"-"`
}

// 创建任务请求体
type CreateTaskReq struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content"`
	Status  string `json:"status"`
}

// 更新任务请求体
type UpdateTaskReq struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Status  string `json:"status"`
}

// 统一错误处理
func handleError(c *gin.Context, err error, statusCode int, message string) {
	if err != nil {
		c.JSON(statusCode, gin.H{"error": message})
		return
	}
}

// 任务列表 GET /api/tasks
func ListTasksHandler(c *gin.Context) {
	var tasks []Task
	err := storage.DB.Order("id desc").Find(&tasks).Error
	handleError(c, err, http.StatusInternalServerError, "查询失败")
	c.JSON(http.StatusOK, tasks)
}

// 创建任务 POST /api/tasks
func CreateTaskHandler(c *gin.Context) {
	var req CreateTaskReq
	err := c.ShouldBindJSON(&req)
	handleError(c, err, http.StatusBadRequest, "参数错误")
	task := Task{
		Title:   req.Title,
		Content: req.Content,
		Status:  req.Status,
	}
	err = storage.DB.Create(&task).Error
	handleError(c, err, http.StatusInternalServerError, "创建失败")
	c.JSON(http.StatusCreated, task)
}

// 任务详情 GET /api/tasks/:id
func GetTaskHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	handleError(c, err, http.StatusBadRequest, "id参数错误")
	var task Task
	err = storage.DB.First(&task, id).Error
	if err != nil {
		handleError(c, err, http.StatusNotFound, "任务不存在")
		return
	}
	c.JSON(http.StatusOK, task)
}

// 编辑任务 PUT /api/tasks/:id
func UpdateTaskHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	handleError(c, err, http.StatusBadRequest, "id参数错误")
	var req UpdateTaskReq
	err = c.ShouldBindJSON(&req)
	handleError(c, err, http.StatusBadRequest, "参数错误")
	var task Task
	err = storage.DB.First(&task, id).Error
	if err != nil {
		handleError(c, err, http.StatusNotFound, "任务不存在")
		return
	}
	// 仅修改提供的字段
	if req.Title != "" {
		task.Title = req.Title
	}
	if req.Content != "" {
		task.Content = req.Content
	}
	if req.Status != "" {
		task.Status = req.Status
	}
	err = storage.DB.Save(&task).Error
	handleError(c, err, http.StatusInternalServerError, "更新失败")
	c.JSON(http.StatusOK, task)
}

// 删除任务 DELETE /api/tasks/:id
func DeleteTaskHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	handleError(c, err, http.StatusBadRequest, "id参数错误")
	err = storage.DB.Delete(&Task{}, id).Error
	handleError(c, err, http.StatusInternalServerError, "删除失败")
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}
