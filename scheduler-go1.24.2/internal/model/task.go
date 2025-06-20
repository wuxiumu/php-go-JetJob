// internal/model/task.go
package model

import (
	"time"
)

// 任务类型枚举
type TaskType string

const (
	TASK_TYPE_SHELL TaskType = "shell" // shell 脚本任务
	TASK_TYPE_HTTP  TaskType = "http"  // http 请求任务
	TASK_TYPE_FILE  TaskType = "file"  // 文件处理任务
	// 后续可扩展更多类型
)

// 任务结构体
type Task struct {
	ID        uint       `gorm:"primaryKey"       json:"id"`
	Name      string     `gorm:"size:128"         json:"name"`                    // 任务名称
	Type      TaskType   `gorm:"size:16"          json:"type"`                    // shell/http/file
	Command   string     `gorm:"size:512"         json:"command"`                 // shell:命令, http:URL, file:路径
	Params    string     `gorm:"type:text"        json:"params"`                  // 任务参数(json字符串)
	Schedule  string     `gorm:"size:64"          json:"schedule"`                // cron表达式，空为手动
	Trigger   string     `gorm:"size:32"          json:"trigger"`                 // manual/schedule/dependency
	DependsOn *uint      `gorm:"column:depends_on" json:"depends_on"`             // 依赖任务ID
	Status    string     `gorm:"size:32"          json:"status"`                  // 任务状态（如active, running, success, failed等）
	Retry     int        `gorm:"column:max_retry" json:"retry"`                   // 最大重试次数
	LastRunAt *time.Time `gorm:"column:last_run_at" json:"last_run_at,omitempty"` // 上次运行时间
	Output    string     `gorm:"type:text"        json:"output"`                  // 执行输出/日志（可选）
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}
